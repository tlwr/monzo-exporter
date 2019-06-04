# Monzo prometheus exporter

Very alpha

## Roadmap

- [x] Export metrics from Monzo
- [x] OAuth token capture
- [x] OAuth token refresh
- [ ] OAuth token persistent storage

## Instructions

### CLI usage:

```
$ monzo_exporter --help

usage: monzo_exporter [<flags>]

Flags:
  --help                         Show context-sensitive help (also try --help-long and --help-man).
  --monzo-oauth-client-id=""     Monzo OAuth client id
  --monzo-oauth-client-secret=""
                                 Monzo OAuth client secret
  --monzo-oauth-port=8080        The port to bind to for serving OAuth
  --monzo-oauth-external-url=""  The URL on which the exporter will be reachable
  --monzo-access-tokens=""       Monzo access tokens comma separated
  --scrape-interval=30           Time in seconds between scrapes
  --metrics-port=9036            The port to bind to for serving metrics
```

### Access tokens from Monzo playground

Using one or many access keys from Monzo API Playground you can run:

```
monzo_exporter --monzo-access-tokens=token1,token2,token3
```

These tokens are only valid for 6 hours.

### Using OAuth flow

This exporter has the ability to export metrics and also do OAuth flows for
Monzo.

You can run:

```
monzo_exporter                                                       \
  --monzo-oauth-client-id     my-client-id-from-monzo-playground     \
  --monzo-oauth-client-secret my-client-secret-from-monzo-playground \
  --monzo-oauth-external-url  https://external-url-for-server
```

Please do TLS termination (HTTPS) in front of the Monzo exporter. This will
involve something like traefik or nginx.

You can configure the port on which the OAuth component listens on with the
flag: `--monzo-oauth-port`, which defaults to port 8080.

The OAuth flow uses a cookie for ensuring that there is no tampering with
authentication. This means that you have to complete the OAuth journey using
the same browser.

Presently this exporter has no persistent state, so restarting the process will
require all users to reauthenticate.
