# sample maven repository

a sample maven repository that serves a jar that serves a hello world over
:8080.

## usage

1. generate the kubernetes deployment (under `./release/release.yaml`)

```
KO_DOCKER_REPO=ghcr.io/cirocosta/sample-maven-repository make release
```

1. deploy

```
kapp deploy -a maven-repository -f ./release
```
