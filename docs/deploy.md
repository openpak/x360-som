# Deploy

Primary host, quadlet `deploy/openpak-x360-som.container` (IP 10.89.2.209 pinned in openpak.network), Traefik
router `x360-som-web` from `deploy/traefik.yml` merged into the shared dynamic file:
`Host(content.xboxlive.com) && PathPrefix(/fffe07d1/)`, priority 50, above the dashboard's Host-only router.
No gateway key: XAM fetches over plain HTTP, not the secure channel.

## Environment

| Variable | Default | |
| --- | --- | --- |
| `X360_SOM_ADDR` | `:10045` | listen address |
| `X360_SOM_MODULES` | `./modules` (`/srv/modules` in the image) | `.lex`/`.xex` files to serve |
| `X360_SOM_SIGNING_KEY` | required | RSA-2048 PEM |
| `X360_SOM_REVISION` | body-hash derived | decimal u64 override |

## Signing key

The existing OpenPak content key `~/openpak/deploy/x360-dashboard/signing/content-key.pem` (RSA-2048, e=65537)
is reused, so the `--xsigner_public_key` pin Wharf already needs for x360-dashboard also covers this service.
Public modulus (hex, big-endian):

```
94B7559F2A70DBB07B145BB98B725842C86367D05F4301AD435ADEC362554D6E32F2CE00F2903FA51CED335DF51BA30B8DAF3057BB221FD17B53E1727361E3A503804B4314880330431CD4D2610D576DA30CF5DF962C1B16007DECF7B7FE1D58897280006873CD96A2485AC2DB42A4DD67D7176DC36EE70674A66382A7DA8E558489A22E5F1BFFAA62B1BBE9968EC7F6DDEA48E7F54D661879619BA0B9283424A403015B2210D956C9E0F4DC1FA26391454A0A453DBDEEDF37935B664A2527E5353184B1A9D3E861102295D763B8D341F90A18F60B0528A4723C48F1CC246071B916EB80598D1ACE3C097E46D0A9FD88C3108F791855F40E881C4FCE78D8537D
```

## LiveHive

x360-liveinfo must advertise `UpdateUriRoot1=http://content.xboxlive.com` and `SOMRevision=<16 hex digits>`
(`X360_LIVEINFO_HIVE_EXTRA`). Read the revision from the x360-som start-up log (`SOMRevision=`); it changes when
the module set changes, so update the hive then. Modules are rebuilt with `modules/build.sh` (luac5.1, python3).
