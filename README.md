
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

Basic request

```shell
curl localhost:8080/posts/1
# do a second request to have it from cache (check logs)
curl localhost:8080/posts/1
```

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
# tells the proxy to record requests on every proxy
curl localhost:1333/*/do-record

# tells the proxy to stop recording requests on proxies that match the pattern
curl localhost:1333/redi*/do-not-record
```

Note that the commander is not protected by any kind of authentication.
The proxy-name is an easy implementation of glob pattern that matches the proxy name: `redi*` matches every proxy name that starts with `redi`.
While `*` matches every proxy name and `redis` matches only the proxy named `redis`.
Use `*dis` to match every proxy name that ends with `dis`.

# Datasources

## in-memory

Does not require any configuration. It just works until the proxy is running.
A reboot will clear the cache.

## redis

Requires a running redis server. The configuration is in the `conf.yaml` file.

```yaml
datasource:
  type: redis
  uri: "redis://localhost:6379/0"
```

Please consider using Docker to run a redis server.

```shell
docker run -d -p 6379:6379 redis
```

## SQLite

Check the `conf.yaml` file for an example configuration.

```yaml
datasource:
  type: sqlite
  uri: "file::memory:?cache=shared"
```

or

```yaml
datasource:
  type: sqlite
  uri: "file.db"
```
