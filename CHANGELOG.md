# Changelog

## v0.1.2 (2026-11-03)

- Every request logs one stderr line, `http <method> <path> -> <status> <duration>ms` (x360-common v0.5.1 `httplog`): no query, headers, cookies, tokens or bodies; `/healthz` skipped.

## v0.1.1

- Quadlet IP moved to 10.89.2.108 (10.89.2.209 belongs to x360-xmessaging).

## v0.1.0

- Serves the signed `system.online.manifest` and the OpenPak-authored `Dash.Social.lex` (Xbox 360 17559).
