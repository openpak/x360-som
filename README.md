# Xbox 360 System Online Modules

OpenPak's Xbox 360 system online module server: the manifest and module files XAM installs into `\sap\online`
after sign-in.

| | |
| --- | --- |
| Network | Xbox 360 |
| Role | `system.online.manifest` plus the dash modules it lists (`Dash.Social.lex`) |
| Depends on | Nothing (standard library only) |
| Image | `ghcr.io/openpak/x360-som` |
| Used by | XAM, over plain HTTP from `UpdateUriRoot1` (no security gateway) |
| Licence | AGPL-3.0-only |

## What it serves

| Path | What |
| --- | --- |
| `GET /fffe07d1/manifests/<rev>/system.online.manifest` | The signed manifest (any `<rev>` gets the current one) |
| `GET /fffe07d1/files/<id>/<name>` | A module file; `<id>` is `%016X%016X%08X` of its SHA1 |
| `GET /healthz` | Health probe |

Both honour `If-None-Match` (ETag = SHA1 hex). The files in [`modules/`](modules/) are OpenPak-authored; no
Microsoft binary is stored here. Not verified in Wharf yet, and stock retail consoles reject the manifest
(retail RSA key) until the console patch replaces key slot 0.

## Run

```sh
X360_SOM_SIGNING_KEY=content-key.pem GOWORK=off go run ./cmd/x360-som
GOWORK=off go test ./...
```

## Docs

- [Protocol](docs/protocol.md)
- [Evidence](docs/evidence.md)
- [Deploy](docs/deploy.md)

## Licence

AGPL-3.0-only. See [LICENSE](LICENSE).
