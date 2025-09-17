# caddy-quic-only

A Caddy module for QUIC-only listeners.

## What is this?

This module allows you to configure Caddy to only accept QUIC connections (HTTP/3) or only TCP connections (HTTP/1.1, HTTP/2), rather than both.

By default, Caddy will bind to both TCP and UDP ports when configured for HTTPS, with TCP handling HTTP/1.1 and HTTP/2, and UDP handling HTTP/3 (QUIC). This module allows you to restrict Caddy to only listen on one protocol or the other.

## Requirements

- Caddy v2

## Installation

### Build from source

```bash
xcaddy build --with github.com/aUsernameWoW/caddy-quic-only=.
```

## Usage

### Caddyfile

```
{
    servers {
        listener_wrappers {
            quic_only {
                mode quic_only
            }
            tls
        }
    }
}

:8443 {
    respond "Hello, QUIC-only world!" 200
}
```

### JSON Configuration

```json
{
  "apps": {
    "http": {
      "servers": {
        "example": {
          "listen": [":8443"],
          "listener_wrappers": [
            {
              "wrapper": "quic_only",
              "mode": "quic_only"
            },
            {
              "wrapper": "tls"
            }
          ],
          "routes": [
            {
              "handle": [
                {
                  "handler": "static_response",
                  "body": "Hello, QUIC-only world!",
                  "status_code": 200
                }
              ]
            }
          ]
        }
      }
    }
  }
}
```

## Modes

- `quic_only` - Only allow QUIC (HTTP/3) connections
- `tcp_only` - Only allow TCP (HTTP/1.1, HTTP/2) connections
- `default` - Allow both QUIC and TCP connections (default behavior)

## How it works

This module implements the `caddy.ListenerWrapper` interface to modify how Caddy handles incoming connections. The wrapper is applied to listeners, and based on the configured mode, it can restrict which protocols are allowed.

Note that fully implementing QUIC-only or TCP-only behavior requires deeper integration with Caddy's server configuration, as the protocol selection happens at the server level rather than just the listener level. This module provides a starting point for that integration.

## License

MIT