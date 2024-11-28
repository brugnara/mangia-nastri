
# Mangia Nastri

Authors: Daniele "vrut" Brugnara, Andrea Leone

**How to get stuff going**

Start the proxy server in one terminal

```shell
# Run the proxy
go run .
```

and in the other send a request through the proxy

```shell
# Send a request through the proxy
curl \
  -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 123" \
  -d '{"content": "stuff"}' \
  localhost:8080/some/endpoint
```
