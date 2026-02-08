# ollama-proxy

A lightweight reverse proxy for [Ollama](https://ollama.com) that logs all incoming requests and responses. Useful for debugging and inspecting traffic between clients and your Ollama instance.

## Features

- Logs request method, path, headers, and body
- Logs response status codes
- Streams responses (flush-on-write)
- Truncates large request bodies in logs (>10KB)

## Usage

```sh
go build -o ollama-proxy && ./ollama-proxy
```

By default the proxy listens on `:8080` and forwards to `http://localhost:11434`.

## Configuration

Command line flags take priority over environment variables.

| Flag | Environment Variable | Default                    | Description                  |
|------|----------------------|----------------------------|------------------------------|
| `-port` | `PROXY_PORT`      | `:8080`                    | Port the proxy listens on    |
| `-url`  | `OLLAMA_URL`      | `http://localhost:11434`   | Upstream Ollama instance URL |

```sh
# Using flags
./ollama-proxy -port :9090 -url http://192.168.1.50:11434

# Using environment variables
PROXY_PORT=:9090 OLLAMA_URL=http://192.168.1.50:11434 ./ollama-proxy
```
