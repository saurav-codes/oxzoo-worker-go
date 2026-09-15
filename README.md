# oxzoo-worker-go

An official ox deploy example: an internal Go background worker built with the standard library only, deployed to a single Ubuntu VPS by the [ox](https://github.com/saurav-codes/vps-ctl) control plane from one `ox.toml` manifest at the repo root. There is no domain, no nginx routing, and no HTTP server; the journal is the product. ox compiles `./worker` with the Go toolchain installed from apt, runs it as a systemd process with `Restart=always`, and the printed stdout lines land in that unit's journal, where `journalctl` reads them back.

## Stack

| Component | Version | Purpose |
|---|---|---|
| Worker | Go 1.24 (`go.mod` directive), standard library only (`os`, `fmt`, `time`) | Prints `hello world oxzoo-worker-go_<GREETING_TAG>` to stdout every 10 seconds |
| Build | apt `go` toolchain, pulled by `required_packages = ["go"]` | Runs `go build -o worker ./cmd/worker` as the deploy install hook |
| Process | systemd via ox, `restart_policy = "always"` | Keeps `./worker` alive and restarts it if it exits |
| Deploy | ox | `ox.toml` defines the process, runtime env file, and install hook |

## Environment flow

One variable, one path:

**`GREETING_TAG`** — `cmd/worker/main.go` reads it once at startup with `os.Getenv("GREETING_TAG")`. If it is missing, the worker prints a clear error to stderr and exits 1, so a misconfigured deploy fails loudly instead of looping with an empty tag. With the value set, every line printed to stdout carries it, and a restart with a new value changes every line from then on. `GREETING_TAG` arrives as runtime env from the ox Environment editor (backed by the `environment_file` in `ox.toml`); `.env.example` documents the variable with a placeholder, and real values live in the ox dashboard, never in git.

The `port = 9119` in `ox.toml` is declared because ox requires a project port even for non-listening workers; this worker never opens a socket, so nothing listens there.

## Deploy with ox

1. Add the repo in the ox dashboard: paste the clone URL `https://github.com/saurav-codes/oxzoo-worker-go`.
2. In the Environment editor, set `GREETING_TAG=w3-05`.
3. Press **Deploy**. ox runs `go build -o worker ./cmd/worker` in the release worktree, injects `GREETING_TAG` as runtime env, and starts `./worker` under systemd with `Restart=always`. No domain is needed at any step.

## Expected output

Follow the worker's journal on the host:

```bash
journalctl -u ox-oxzoo-worker-go-worker.service
```

With `GREETING_TAG=w3-05` set in the Environment editor, the journal shows the exact line `hello world oxzoo-worker-go_w3-05`, repeating every 10 seconds:

```
hello world oxzoo-worker-go_w3-05
hello world oxzoo-worker-go_w3-05
```

Unit name anatomy: ox names every process unit `ox-<project>-<process>.service`. The project name is `oxzoo-worker-go` and the `[[processes]]` entry is `worker`, so the unit is `ox-oxzoo-worker-go-worker.service`. If the worker exits, systemd restarts it per `restart_policy = "always"` and the restarts show up inline in the same journal.

## Local development

```bash
go build -o worker ./cmd/worker
GREETING_TAG=dev ./worker
```

Without `GREETING_TAG` the worker exits 1 with an error on stderr instead of printing. Pass env inline per the commands above; never commit a real `.env`.
