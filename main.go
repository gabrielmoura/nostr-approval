package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"strconv"
	"sync"
	"time"
)

type Community struct {
	Id   string
	Name string
}

type Keys struct {
	pub string
	prv string
}

func (k *Keys) DecodedPub() string {
	_, v, _ := nip19.Decode(k.pub)
	return v.(string)
}
func (k *Keys) DecodedPrv() string {
	_, v, _ := nip19.Decode(k.prv)
	return v.(string)
}

type Connection struct {
	Context      context.Context
	Relay        *nostr.Relay
	CreatedList  []nostr.Event
	ApprovedList []nostr.Event
	ToApproval   []nostr.Event
	ToApproved   []nostr.Event
	ErrorList    []error
	Wg           sync.WaitGroup
	sync.Mutex
}

func NewConnection(ctx context.Context, relayURL string) (*Connection, error) {
	relay, err := nostr.RelayConnect(context.TODO(), relayURL)
	if err != nil {
		return nil, err
	}

	return &Connection{
		Relay:        relay,
		CreatedList:  []nostr.Event{},
		ApprovedList: []nostr.Event{},
		ToApproval:   []nostr.Event{},
		ToApproved:   []nostr.Event{},
		ErrorList:    []error{},
		Context:      ctx,
	}, nil
}

func (conn *Connection) SubscribeToFilter(filter nostr.Filter, eventList *[]nostr.Event, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(conn.Context, timeout)

	defer func() {
		cancel()
		conn.Wg.Done()
	}()

	sub, err := conn.Relay.Subscribe(ctx, nostr.Filters{filter})
	if err != nil {
		conn.Lock()
		conn.ErrorList = append(conn.ErrorList, err)
		conn.Unlock()
		return
	}

	for ev := range sub.Events {
		fmt.Println(ev)
		conn.Lock()
		*eventList = append(*eventList, *ev)
		conn.Unlock()
	}
}

func (conn *Connection) CreateApprovalEvents(keys *Keys, community *Community) {
	var toApproved []nostr.Event

	for _, event := range conn.ToApproval {
		j, _ := json.Marshal(event)
		ev := nostr.Event{
			PubKey: keys.DecodedPub(),
			Kind:   nostr.KindCommunityPostApproval,
			Tags: nostr.Tags{
				nostr.Tag{"a", fmt.Sprintf("%d:%s:%s", nostr.KindCommunityDefinition, community.Id, community.Name)},
				nostr.Tag{"e", event.ID},
				nostr.Tag{"p", event.PubKey},
				nostr.Tag{"k", strconv.Itoa(event.Kind)},
			},
			Content:   string(j),
			CreatedAt: nostr.Now(),
		}
		if err := ev.Sign(keys.DecodedPrv()); err == nil {
			toApproved = append(toApproved, ev)
		} else {
			conn.ErrorList = append(conn.ErrorList, err)
		}
	}

	conn.ToApproved = toApproved
}

func (conn *Connection) PublishEvents() {
	for _, ev := range conn.ToApproved {
		time.Sleep(300 * time.Millisecond)
		if err := conn.Relay.Publish(conn.Context, ev); err != nil {
			conn.ErrorList = append(conn.ErrorList, err)
		}
	}
}

func main() {
	fRelay := flag.String("relay", "wss://nos.lol", "Relay URL")
	fCommunityId := flag.String("cid", "", "Community ID")
	fCommunityName := flag.String("cname", "", "Community Name")
	fPubKey := flag.String("pub-key", "", "Public Key")
	fPrvKey := flag.String("prv-key", "", "Private Key")
	fSinceTime := flag.Int64("since-time", 24, "Since Time (Hours)")

	flag.Parse()

	if *fCommunityId == "" || *fCommunityName == "" || *fPubKey == "" || *fPrvKey == "" {
		fmt.Println("Missing required parameters")
		fmt.Println("Usage: nostr-approval -cid <ID> -cname <NAME> -pub-key <PUBKEY> -prv-key <PRVKEY>")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := NewConnection(ctx, *fRelay)
	if err != nil {
		fmt.Println("Failed to connect to relay:", err)
		return
	}
	conn.Wg.Add(2)

	community := &Community{
		Id:   *fCommunityId,
		Name: *fCommunityName,
	}
	keys := &Keys{
		pub: *fPubKey,
		prv: *fPrvKey,
	}

	var timestamp nostr.Timestamp
	timestamp = nostr.Timestamp(time.Now().Add(-time.Duration(*fSinceTime) * time.Hour).Unix())

	searchCreated := nostr.Filter{
		Kinds: []int{nostr.KindTextNote, nostr.KindReaction, nostr.KindZap, nostr.KindArticle, nostr.KindWikiArticle},
		Tags:  nostr.TagMap{"a": {fmt.Sprintf("%d:%s:%s", nostr.KindCommunityDefinition, community.Id, community.Name)}},
		Since: &timestamp,
	}

	searchApproved := nostr.Filter{
		Kinds:   []int{nostr.KindCommunityPostApproval},
		Authors: []string{community.Id},
		Tags:    nostr.TagMap{"a": {fmt.Sprintf("%d:%s:%s", nostr.KindCommunityDefinition, community.Id, community.Name)}},
		Since:   &timestamp,
	}
	go conn.SubscribeToFilter(searchCreated, &conn.CreatedList, 10*time.Second)
	go conn.SubscribeToFilter(searchApproved, &conn.ApprovedList, 10*time.Second)

	conn.Wg.Wait()
	for _, created := range conn.CreatedList {
		found := false
		for _, approved := range conn.ApprovedList {
			if created.ID == approved.Tags.GetFirst([]string{"e"}).Value() {
				found = true
				break
			}
		}
		if !found {
			conn.ToApproval = append(conn.ToApproval, created)
		}
	}

	conn.CreateApprovalEvents(keys, community)

	conn.PublishEvents()

	fmt.Printf("Criados: %d\n", len(conn.CreatedList))
	fmt.Printf("Aprovados: %d\n", len(conn.ApprovedList))
	fmt.Printf("Para aprovar: %d\n", len(conn.ToApproval))
	fmt.Printf("Novos Aprovados: %d\n", len(conn.ToApproved))
	fmt.Printf("Erros: %d\n", len(conn.ErrorList))

	for _, err := range conn.ErrorList {
		fmt.Println(err)
	}
}
