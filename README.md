# Bot para Aprovação de Postagens no Nostr

Este é um bot para facilitar a aprovação de postagens no Nostr, verificando e aprovando automaticamente as postagens mais recentes dentro de um intervalo de tempo configurável.

## Uso

Para aprovar postagens no Nostr, execute o seguinte comando. Ele buscará as últimas postagens e aprovações realizadas nas últimas 24 horas e aprovará aquelas que ainda não foram aprovadas.

### Comando

```bash
nostr-approval -cid <ID> -cname <NAME> -pub-key <PUBKEY> -prv-key <PRVKEY>
```

### Parâmetros

- `-cid`: **Obrigatório**. ID da comunidade.
- `-cname`: **Obrigatório**. Nome da comunidade.
- `-pub-key`: **Obrigatório**. Chave pública.
- `-prv-key`: **Obrigatório**. Chave privada.
- `-relay`: **Opcional**. Endereço do servidor de retransmissão.
- `-since-time`: **Opcional**. Define o intervalo de tempo (em horas) para busca de postagens.
- `-pow`: **Opcional**. Define o nível de dificuldade do PoW.

## Recomendações

É recomendado usar um cron para rodar o comando periodicamente e garantir que todas as postagens sejam aprovadas no tempo adequado.

### Exemplo de Cron (a cada 1 hora)

```bash
0 * * * * nostr-approval -cid <ID> -cname <NAME> -pub-key <PUBKEY> -prv-key <PRVKEY> -since-time 1
```

Isso permite aprovar as mensagens de forma mais rápida, dentro de um intervalo de uma hora.

## Contribuições

Sinta-se à vontade para contribuir com melhorias ou novas funcionalidades para o projeto. Abra um PR ou reporte problemas na aba de issues.
Ou enviar Sats para `verdantkite75@walletofsatoshi.com` para apoiar o desenvolvimento.

---

Feito com ❤️ para a comunidade Nostr.
