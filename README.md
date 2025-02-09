
# Mangia Nastri

Authors: Andrea Leone, Daniele Brugnara

**How to get stuff going**

Start the proxy server in one terminal

```shell
# Run the proxy
go run .
```

or watch using wgo
```shell
brew install wgo

wgo run -file conf.yaml .
```

and in the other send a request through the proxy

```shell
# Send a request through the proxy
curl \
  -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer 123" \
  -d '{"content": "stuff", "anther_param": true, "nested_object": {"i": {"j": "k", "l": {"m": "n"}}, "c": {"d": "e", "f": ["g", "h"]}, "a": "b"}}' \
  localhost:8080/some/endpoint

# send this too. The sort of headers and or body should generate the same hash
curl \
  -X POST \
  -H "Authorization: Bearer 123" \
  -H "Content-Type: application/json" \
  -d '{"anther_param": true, "content": "stuff", "nested_object": {"i": {"j": "k", "l": {"m": "n"}}, "c": {"f": ["g", "h"], "d": "e"}, "a": "b"}}' \
  localhost:8080/some/endpoint
```

# commander

Commander allow to send commands to the proxy servers.

```shell
# tells the proxy to record requests
curl localhost:1333/do-record

# tells the proxy to stop recording requests
curl localhost:1333/do-not-record
```
