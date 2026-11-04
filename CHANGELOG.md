# Changelog

## v0.1.3 (2026-11-04)

- Fix the manifest body magic: XAM FUN_816EBE40 requires the bytes `xonm` (online) at body offset 0, not `mnox`; with `mnox` XAM rejected the downloaded manifest with 8B050004 (found in a Wharf run with `\sap` mounted). The signature, SOMRevision and module set are unchanged apart from the signed body bytes.

## v0.1.2 (2026-11-03)

- Every request logs one stderr line, `http <method> <path> -> <status> <duration>ms` (x360-common v0.5.1 `httplog`): no query, headers, cookies, tokens or bodies; `/healthz` skipped.

## v0.1.1

- Quadlet IP moved to 10.89.2.108 (10.89.2.209 belongs to x360-xmessaging).

## v0.1.0

- Serves the signed `system.online.manifest` and the OpenPak-authored `Dash.Social.lex` (Xbox 360 17559).
