# Monzo prometheus exporter

Very alpha

## Roadmap

- [x] Export metrics from Monzo
- [ ] OAuth token capture & refresh

## Instructions

```
docker run \
       --rm -it \
       tlwr/monzo_exporter \
       --monzo-access-tokens=token1,token2,token3
```

Where `token1`, `token2`, `token3` are OAuth tokens (you can get these from
Monzo playground).

As per the roadmap above OAuth token capture and refresh is not implemented
yet.
