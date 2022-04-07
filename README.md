# go-gemini-hn

## Build & Run

```shell
cd Development/personal/go-gemini-hn
go run ./cmd/...

docker build --tag geberl/go-gemini-hn:latest .

docker run --detach \
           --env "TZ=Europe/Berlin" \
           --env "HN_BASE_URL=hn.eberl.se" \
           --env "HN_LOG_LEVEL=debug" \
           --publish 1965:1965 \
           --name go-gemini-hn \
           --restart unless-stopped \
           geberl/go-gemini-hn:latest
```

## Dependencies

- https://pkg.go.dev/git.sr.ht/~adnano/go-gemini
- https://git.sr.ht/~adnano/go-gemini/tree/v0.2.2/text.go#L30
- https://sr.ht/~adnano/go-gemini/

- https://github.com/lukakerr/hkn
- https://github.com/HackerNews/API

- https://github.com/jaytaylor/html2text
