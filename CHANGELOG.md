# Changelog

## v0.1.4 (2026-11-05)

- Fix manifest entries so XAM can find modules by name (verified in Wharf: `XdfCacheItem` now fetches `/fffe07d1/files/<id>/<name>` and stores it in `\sap\online`): the lookup name offset must be at entry+0x00 (FUN_816ECF20), and the string length word counts the terminating NUL (FUN_816ECDA0 checks the last counted byte is 0). Before this every by-name lookup failed with 8B050009 and no file was ever requested. The signed body changes, so the auto SOMRevision is now `91CBB7CA471F1F2E` (was `90466AC43FB4F50C`): update the LiveHive `SOMRevision`.

## v0.1.3 (2026-11-04)

- Fix the manifest body magic: XAM FUN_816EBE40 requires the bytes `xonm` (online) at body offset 0, not `mnox`; with `mnox` XAM rejected the downloaded manifest with 8B050004 (found in a Wharf run with `\sap` mounted). The signature, SOMRevision and module set are unchanged apart from the signed body bytes.

## v0.1.2 (2026-11-03)

- Every request logs one stderr line, `http <method> <path> -> <status> <duration>ms` (x360-common v0.5.1 `httplog`): no query, headers, cookies, tokens or bodies; `/healthz` skipped.

## v0.1.1

- Quadlet IP moved to 10.89.2.108 (10.89.2.209 belongs to x360-xmessaging).

## v0.1.0

- Serves the signed `system.online.manifest` and the OpenPak-authored `Dash.Social.lex` (Xbox 360 17559).
