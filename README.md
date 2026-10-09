# oxzoo-worker-go

Deployed with [ox](https://deploywithox.com): deploy a repo to your own server with one command, no Docker. [Docs](https://deploywithox.com/docs) · [Stack guides](https://deploywithox.com/docs/guides)

An [ox](https://deploywithox.com) deploy example: an internal Go background worker built with the standard library only, deployed to your own Ubuntu server. There is no domain and no HTTP server; the logs are the product. ox builds `./worker` and runs it under systemd as a worker that restarts if it exits.

## Stack

| Component | Version | Purpose |
|---|---|---|
| Worker | Go 1.24 (from `go.mod`), standard library only | prints `hello world oxzoo-worker-go_<GREETING_TAG>` every 10 seconds |
| Toolchain | Go from `go.mod`, installed by ox with mise | builds `./worker` |

## ox.toml

```toml
# A Go background worker: no web process, no domain.

[app]
enabled = false

[workers]
worker = "./worker"

[build]
commands = ["go build -o worker ./cmd/worker"]
```

`[app] enabled = false` says the project has no web process, so ox adds none and asks for no domain.

## Environment flow

`cmd/worker/main.go` reads `GREETING_TAG` once at startup. If it is missing, the worker prints an error and exits 1, so a misconfigured deploy fails loudly instead of looping with an empty tag. Changing it with `ox vars set` redeploys, and the next lines carry the new value.

## Deploy with ox

```sh
curl -fsSL https://deploywithox.com/install.sh | sh
ox login
ox new https://github.com/saurav-codes/oxzoo-worker-go
printf 'GREETING_TAG=demo\n' | ox review oxzoo-worker-go --from-file - --wait
```

The plan, offline:

```console
$ ox check .
ox check . (manifest: ox.toml)

  build.commands[0]          go build -o worker ./cmd/worker                      declared
  workers.worker             ./worker                                             declared
  tools.go                   1.24                                                 detected:go.mod

  Provided by ox: PORT, HOST, OX_ENV, OX_PROJECT, OX_RELEASE, OX_DATA_DIR
  Set on the dashboard before the first deploy: GREETING_TAG
  hint: [build] commands replaces the detected build step "go build -o .ox/bin/app ./cmd/worker" (from go.mod); add it to the list if the release still needs it

Ready to deploy.
```

## Expected output

```sh
ox logs oxzoo-worker-go --follow
```

shows `hello world oxzoo-worker-go_<GREETING_TAG>` every 10 seconds.

## Local development

```sh
go build -o worker ./cmd/worker
GREETING_TAG=dev ./worker
```
