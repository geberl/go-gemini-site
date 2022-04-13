# go-gemini-site

![Go](https://img.shields.io/badge/go-1.18-orange.svg)
![Alpine](https://img.shields.io/badge/alpine-3.15-lightgrey.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)

## Build & Run

Running in place:

```shell
go run ./cmd/...
```

Building a Docker image locally:

```shell
docker build --tag geberl/go-gemini-site:latest .

docker run --detach \
           --env "TZ=Europe/Berlin" \
           --env "GEMINI_SITE_BASE_URL=localhost" \
           --env "GEMINI_SITE_LOG_LEVEL=debug" \
           --publish 1965:1965 \
           --name go-gemini-site \
           --restart unless-stopped \
           geberl/go-gemini-site:latest
```

Downloading and running the Docker image from GitHub Container Registry:

```shell
export CR_PAT=YOUR_TOKEN
echo $CR_PAT | docker login ghcr.io --username geberl --password-stdin

docker pull ghcr.io/geberl/go-gemini-site:latest

docker run --detach \
           --env "TZ=Europe/Berlin" \
           --env "GEMINI_SITE_BASE_URL=eberl.se" \
           --env "GEMINI_SITE_LOG_LEVEL=error" \
           --mount type=volume,source=go-gemini-site-certs,destination=/app/certs \
           --publish 1965:1965 \
           --name go-gemini-site \
           --restart unless-stopped \
           ghcr.io/geberl/go-gemini-site:latest
```

## Dependencies

- https://pkg.go.dev/git.sr.ht/~adnano/go-gemini
- https://git.sr.ht/~adnano/go-gemini/tree/v0.2.2/text.go#L30
- https://sr.ht/~adnano/go-gemini/

- https://github.com/lukakerr/hkn
- https://github.com/HackerNews/API

- https://github.com/jaytaylor/html2text
