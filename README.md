# Return Static Response

Return static response is a middleware plugin for Traefik which takes an incoming request on a specific HTTP method and path match, and then returns a static response.

## Configuration

### Static

```yaml
pilot:
  token: "xxxx"

experimental:
  plugins:
    returnStaticResponse:
      moduleName: github.com/jamesmcroft/traefik-plugin-return response
      version: "v1.0.0"
```

### Dynamic

To configure the Return Static Response plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in your dynamic configuration as explained [here](https://doc.traefik.io/traefik/middlewares/overview/). The following example creates and uses the `returnStaticResponse` middleware plugin to return a static response for a specific path.

```yaml
http:
  services:
    serviceRoot:
      loadBalancer:
        servers:
          - url: "http://localhost:8080"

  middlewares:
    options-static-response:
      plugin:
        returnStaticResponse:
          response:
            method: "OPTIONS"
            url_match: "^http://localhost:8080/(.+)$"
            status_code: 200

  routers:
    routerRoot:
      rule: "PathPrefix(`/`)"
      service: "serviceRoot"
      middlewares:
        - "options-static-response"
```
