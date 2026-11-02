# Evidence: x360-som

From the static RE note `re-som-manifest` (XAM base 0x815F0000, dash base 0x92000000). [V] verified in code,
[I] inferred. No emulator run has confirmed any of it yet.

- [V] URL formats `%s/%08x/manifests/%s/%s` @0x815F4204 and `%s/%08x/files/%s/%s` @0x815F910C, id format
  `%016I64X%016I64X%08X` @0x815F90F4, built in FUN_8168bef0 and FUN_816ed888; root = hive `UpdateUriRoot1`,
  revision = hive `SOMRevision`.
- [V] Container and body layout: FUN_816ebe40 @0x816EBE40 (XMNP, SHA1 @0x24, signature @0x38, `mnox`),
  entry iterator FUN_816ece70 (stride 0x60), strings FUN_816ecda0 (`u16` length).
- [V] The manifest is RSA-verified with key slot 0 (retail key) on every install path: stock retail rejects it
  (0x8B050012); Wharf or a patched kernel can accept it. This corrects re-lua-modules.md section 1.
- [V] Transport is XAM's own HTTP client with `X-FLOKEN`, not the secure channel, so there is no gateway guard.
- [V] The dash VM does no hash or signature check on a `.lex` chunk; it needs `Main`.
- [I] PKCS#1 v1.5 SHA-1 padding (the note does not name it); the descriptor flag that makes XAM hash a `.lex` as a
  XEX header (FUN_8168cb20) is unknown, so a raw-Lua `.lex` may fail that check.
- [I] Entry fields 0x04, 0x20, 0x24 and the header bytes 0x04..0x27 are not load-bearing per the note.

Probes for a Wharf run: FUN_816ebe40, FUN_8168ba58 (`\sap\online`), FUN_8168bef0, FUN_8168cc38, FUN_816ed670.
