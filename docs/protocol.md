# Protocol: x360-som

XAM 17559 fetches `<UpdateUriRoot1>/fffe07d1/manifests/<SOMRevision>/system.online.manifest` when the hive
`SOMRevision` differs from its stored copy, then each listed file from
`<UpdateUriRoot1>/fffe07d1/files/<id>/<name>`. Plain HTTP GET, `X-FLOKEN` header present (never logged),
`If-None-Match` honoured, 200 with `application/octet-stream` or 304. Everything below is big-endian.

## Manifest container

| Offset | Size | Field |
| --- | --- | --- |
| 0x000 | 4 | `XMNP` |
| 0x004 | 4 | version, 1 (the rest of 0x004..0x023 is zero) |
| 0x024 | 20 | SHA1 of the body |
| 0x038 | 0x100 | RSA-2048 signature over that SHA1 |
| 0x138 | | body |

Signature: PKCS#1 v1.5 with SHA-1, the same scheme and OpenPak content key as x360-dashboard's YTGR envelope.
XAM verifies it with `XeKeysVerifyRSASignature` key slot 0; the padding XAM expects is not confirmed in the RE
note (see docs/evidence.md).

## Body

| Offset | Size | Field |
| --- | --- | --- |
| 0x00 | 4 | `xonm` |
| 0x04 | 4 | signed length (whole body) |
| 0x10 | 4 | type, 3 |
| 0x18 | 8 | revision (SOMRevision) |
| 0x28 | 4 | entry table offset, 0x30 |
| 0x30 | 4 | entry count N |
| 0x34 | N*0x60 | entries |
| | | string pool: `u16 length`, bytes, `0` |

Entry (offsets within the 0x60 bytes): 0x04 u64 version (1); 0x14 u32 flags (0); 0x18 u32 node type (3, module);
0x1C u32 body offset of the name string; 0x20 u32 display-name offset (same string); 0x24 20-byte secondary id
(the file id); 0x38 20-byte file id. The file id is the SHA1 of the served file. The URL id is the same 20 bytes
as `%016X%016X%08X` (u64, u64, u32), i.e. 40 upper-case hex digits.

## Revision

`X360_SOM_REVISION` (decimal u64) if set. Otherwise the first 8 bytes of the SHA1 of the body built with
revision 0, so the revision changes exactly when the module set changes. The hive value `SOMRevision` must be the
16-hex-digit form; the start-up log line shows it as `SOMRevision=`. A request for another revision is
answered with the current manifest and logged as a mismatch.

## Modules

`*.lex` and `*.xex` files in `X360_SOM_MODULES`. `Dash.Social.lex` is a Lua 5.1 big-endian chunk built by
`modules/build.sh` from `Dash.Social.lua`; the dash VM calls its globals `SetAppConfig`, `Main`, `BackgroundApp`.
