## Este é apenas um esboço de como seria a implementtação pada validar a ideia.

## Definições
	KindTextNote = 1
	KindReaction = 7
	KindZap = 9735
	KindCommunityPostApproval = 4550
	KindCommunityDefinition = 34550
    KindArticle             = 30023
    KindWikiArticle         = 30818

## Passos
1. **Conectar ao Relay:**
* O código primeiro se conecta a um servidor de retransmissão Nostr usando `nostr.RelayConnect`. Um servidor de retransmissão ajuda a descobrir e se comunicar com outros participantes na rede Nostr.
* O servidor de retransmissão específico usado aqui é "[https://relay.openintents.org](https://relay.openintents.org)".

2. **Defina Comunidade e Chaves:**
* O código define um objeto `Comunidade` que representa a comunidade de interesse. Neste caso, o ID da comunidade é definido como "135d2b016eb41672477291ea7bcafe5f00e007dc6612610ff58a08655bc1b095" e o nome é definido como "Brasil".
* Ele também define um objeto `Chaves`, mas deixa os campos de chave pública e privada vazios (`Pub` e `Prv` são definidos como strings vazias).

3. **Definir filtros:**
* O código define três filtros Nostr usando a estrutura `nostr.Filter`:
* `searchCreated`: Este filtro pesquisa eventos dos tipos `nostr.KindTextNote`, `nostr.KindReaction` e `nostr.KindZap` que são marcados com um formato específico, indicando que estão relacionados à definição da comunidade com o ID e o nome fornecidos.
* `searchApproved`: Este filtro pesquisa eventos do tipo `nostr.KindCommunityPostApproval` criados pela comunidade e marcados com o mesmo formato que `searchCreated`.
* `searchCommunityDefinition`: Este filtro pesquisa eventos do tipo `nostr.KindCommunityDefinition` criados pela comunidade e marcados com o nome da comunidade.

4. **Inscrever-se em eventos:**
* O código se inscreve em eventos que correspondem aos filtros definidos usando `relay.Subscribe`. Ele assina os filtros `searchCommunityDefinition` e `searchCreated` para recuperar informações sobre a comunidade e os eventos criados dentro dela.
* Os resultados são armazenados em canais separados (`reqC` e `reqA`).

5. **Processar eventos:**
* O código itera pelos eventos recebidos em ambos os canais (`reqC` e `reqA`) e os armazena em listas separadas: `createdList` para eventos criados e `approvedList` para eventos aprovados.

6. **Identificar eventos para aprovação:**