# Block and rewrite header or body

Block Rewrite is a middleware plugin for [Traefik](https://github.com/traefik/traefik) which block all path and rewrite response headers or body.

## Configuration

## Static

```toml
[pilot]
    token="xxx"

[experimental.plugins.blockrewrite]
    modulename = "github.com/czc006/block-rewrite-header-body"
    version = "v0.0.2"
```

## Dynamic

To configure the `Block Rewrite` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in 
your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). The following example creates
and uses the `Block Rewrite` middleware plugin to block all HTTP requests and rewrite response headers or body. 

```yaml
# Dynamic configuration
http:
  routers:
    my-router:
      rule: host(`localhost`)
      service: my-service
      entryPoints:
        - web
      middlewares:
        - my-plugin

  services:
   service-foo:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
  
  middlewares:
    my-plugin:
      plugin:
        blockrewrite:
          headers:
            Foo: Bar   # set header Foo=Bar
            code: 200  # set http status code
          body: "test body"  # write response body
```
