#!/bin/sh
# Build Dash.Social.lex for the Xbox 360 17559 dash Lua 5.1 VM.
#   1. luac5.1 -s  (strip debug) -> host little-endian chunk
#   2. le2be.py    -> big-endian, size_t=4 console chunk (inverse of be2le.py)
# No signing here: the .lex itself has no content check at the dash VM. Its
# SHA1 goes into the (MS-signed) system.online.manifest entry; compute with
#   sha1sum Dash.Social.lex   (XAM compares XeCryptSha of the file, 0x8168CB20).
set -e
cd "$(dirname "$0")"
luac5.1 -s -o Dash.Social.luc Dash.Social.lua
python3 -I le2be.py Dash.Social.luc Dash.Social.lex
rm -f Dash.Social.luc
printf 'sha1  '; sha1sum Dash.Social.lex
printf 'size  '; wc -c < Dash.Social.lex
printf 'hdr   '; od -An -tx1 -N12 Dash.Social.lex
