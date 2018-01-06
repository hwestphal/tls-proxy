TLS Proxy
=========

This is a tiny TLS reverse proxy written in Go, suitable for running in front of a Varnish.

It sets the `X-Forwared-For`, `X-Forwared-Protocol` and `X-Forwared-Port` request headers correctly and rewrites a `Location` response header is necessary.


Usage:
------

```
Usage of tls-proxy:
  -cert string
        path to PEM certificate (default "cert.pem")
  -flush-interval duration
        minimum duration between flushes to the client (default: off)
  -key string
        path to PEM key (default "key.pem")
  -listen string
        bind address to listen on (default ":8443")
  -logging
        log requests (default true)
  -where string
        address to forward connections to (default "127.0.0.1:8080")
```
