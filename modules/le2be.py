#!/usr/bin/env python3
# Inverse of store-re/be2le.py. Read a host luac5.1 -s chunk (little-endian,
# sizeof int=4, size_t=8, Instruction=4, lua_Number=8 double) and re-serialise
# as an Xbox 360 dash Lua 5.1 chunk: BIG-endian, sizeof int=4, size_t=4,
# Instruction=4, lua_Number=8 (BE double), integral=0.
import struct, sys
src = open(sys.argv[1], 'rb').read()
o = [0]
def ru8():
    v = src[o[0]]; o[0]+=1; return v
def ri32():
    v = struct.unpack('<I', src[o[0]:o[0]+4])[0]; o[0]+=4; return v
def rsz():
    v = struct.unpack('<Q', src[o[0]:o[0]+8])[0]; o[0]+=8; return v
def rraw(n):
    v = src[o[0]:o[0]+n]; o[0]+=n; return v
W = []
def w8(v): W.append(bytes([v]))
def w32(v): W.append(struct.pack('>I', v))      # int / Instruction / size_t(4)
def wstr():
    n = rsz()                                    # host size_t (8)
    W.append(struct.pack('>I', n))               # console size_t (4, BE)
    if n: W.append(rraw(n))
def func():
    wstr()                                       # source
    w32(ri32()); w32(ri32())                     # line defined / last
    for _ in range(4): w8(ru8())                 # nups,numparams,isvararg,maxstack
    n = ri32(); w32(n)                           # code
    for _ in range(n): w32(ri32())
    n = ri32(); w32(n)                           # constants
    for _ in range(n):
        t = ru8(); w8(t)
        if t == 1: w8(ru8())
        elif t == 3: W.append(struct.pack('>d', struct.unpack('<d', rraw(8))[0]))
        elif t == 4: wstr()
        elif t != 0: raise SystemExit('bad const %d' % t)
    n = ri32(); w32(n)                           # protos
    for _ in range(n): func()
    n = ri32(); w32(n)                           # lineinfo
    for _ in range(n): w32(ri32())
    n = ri32(); w32(n)                           # locvars
    for _ in range(n): wstr(); w32(ri32()); w32(ri32())
    n = ri32(); w32(n)                           # upvalues
    for _ in range(n): wstr()
# host header is 12 bytes; validate then skip
assert src[:6] == b'\x1bLuaQ\x00' and src[6]==1 and src[7]==4 and src[8]==8 and src[9]==4 and src[10]==8, 'unexpected host header'
o[0] = 12
W.append(b'\x1bLuaQ\x00\x00\x04\x04\x04\x08\x00')   # BE, size_t=4
func()
open(sys.argv[2], 'wb').write(b''.join(W))
print('wrote', sys.argv[2], sum(len(x) for x in W), 'bytes')
