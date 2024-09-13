# bitcoin-amb-proxy

proxy bitcoin rpc requests to ["Amazon Managed Blockchain" (AMB)](https://docs.aws.amazon.com/managed-blockchain/latest/ambbtc-dg/what-is-service.html) without the v4 signing dance.


### run

run the proxy
```
$ go run . -aws-key="<<aws-key>>" -aws-secret=<<aws-secret>>
2024/09/12 11:11:11 Starting server on 0.0.0.0:8787
```


use the vanilla bitcoin-cli (or your app) to send rpc requests to localy running proxy \
and they proxy to AMB
```
$ bitcoin-cli -rpcconnect=127.0.0.1 -rpcport=8787 getbestblockhash
00000000000000000002895e8b890f316fee10df337d8f52d54bc276f2be6d18


$ bitcoin-cli -rpcconnect=127.0.0.1 -rpcport=8787 getblockchaininfo
{
  "chain": "main",
  "blocks": 861098,
  "headers": 861098,
  "bestblockhash": "00000000000000000002895e8b890f316fee10df337d8f52d54bc276f2be6d18",
  "difficulty": 92671576265123.32,
  "time": 1726200123,
  "mediantime": 1726196123,
  "verificationprogress": 0.9999984212345678,
  "initialblockdownload": false,
  "chainwork": "00000000000000000000000000000000000000008e350859c34b9bcba8aeb731",
  "size_on_disk": 682174999191,
  "pruned": false,
  "warnings": "..."
}
```


### todo
- metrics
- more credential providers
- tests
- ...