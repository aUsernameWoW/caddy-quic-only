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

#### QUIC-only mode:
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

#### TCP-only mode:
```
{
    servers {
        listener_wrappers {
            quic_only {
                mode tcp_only
            }
            tls
        }
    }
}

:8443 {
    respond "Hello, TCP-only world!" 200
}
```

### JSON Configuration

#### QUIC-only mode:
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

#### TCP-only mode:
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
              "mode": "tcp_only"
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
                  "body": "Hello, TCP-only world!",
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

In addition to wrapping listeners, the module also modifies the server's protocol configuration to ensure that only the appropriate protocols are enabled:
- In `quic_only` mode, only HTTP/3 (h3) is enabled
- In `tcp_only` mode, only HTTP/1.1 (h1) and HTTP/2 (h2) are enabled
- In `default` mode, all protocols are enabled as configured by the server

This approach provides a more complete solution than just wrapping listeners, as it prevents Caddy from even attempting to create listeners for protocols that are not wanted.

## License

MIT