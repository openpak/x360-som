// Package som builds and serves the Xbox 360 system online manifest (system.online.manifest) and
// the module files it lists. Layout: docs/protocol.md, derived from XAM 17559 (re-som-manifest).
package som

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	hdrLen    = 0x138 // XMNP container prefix; the body starts here
	bodyHdr   = 0x30  // body header; the entry table follows it
	entryLen  = 0x60
	titlePath = "/fffe07d1/"
)

// File is one module served from the modules directory.
type File struct {
	Name string
	Data []byte
	ID   [20]byte // SHA1 of Data: the manifest file id and the URL id
}

// LoadDir reads every .lex and .xex file in dir (sources and build scripts beside them are ignored), sorted by name.
func LoadDir(dir string) ([]File, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fs []File
	for _, e := range ents {
		if ext := filepath.Ext(e.Name()); !e.Type().IsRegular() || (ext != ".lex" && ext != ".xex") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		fs = append(fs, File{Name: e.Name(), Data: b, ID: sha1.Sum(b)})
	}
	sort.Slice(fs, func(i, j int) bool { return fs[i].Name < fs[j].Name })
	return fs, nil
}

// URLID formats a file id as XAM does: %016I64X%016I64X%08X over the 20 bytes.
func URLID(id [20]byte) string {
	return fmt.Sprintf("%016X%016X%08X", binary.BigEndian.Uint64(id[0:]), binary.BigEndian.Uint64(id[8:]), binary.BigEndian.Uint32(id[16:]))
}

func body(rev uint64, files []File) []byte {
	T := bodyHdr
	pool := T + 4 + entryLen*len(files)
	b := make([]byte, pool)
	copy(b, "xonm")
	binary.BigEndian.PutUint32(b[0x10:], 3)
	binary.BigEndian.PutUint64(b[0x18:], rev)
	binary.BigEndian.PutUint32(b[0x28:], uint32(T))
	binary.BigEndian.PutUint32(b[T:], uint32(len(files)))
	for i, f := range files {
		off := uint32(len(b))
		n := len(f.Name) + 1 // XAM's string reader (FUN_816ECDA0) counts the NUL and checks the last byte is 0
		b = append(b, byte(n>>8), byte(n))
		b = append(b, f.Name...)
		b = append(b, 0)
		e := b[T+4+i*entryLen:]
		binary.BigEndian.PutUint32(e[0x00:], off) // lookup name (FUN_816ECF20 compares this)
		binary.BigEndian.PutUint64(e[0x04:], 1)
		binary.BigEndian.PutUint32(e[0x14:], 0) // flags: online entry
		binary.BigEndian.PutUint32(e[0x18:], 3) // node type: module
		binary.BigEndian.PutUint32(e[0x1c:], off)
		binary.BigEndian.PutUint32(e[0x20:], off)
		copy(e[0x24:], f.ID[:])
		copy(e[0x38:], f.ID[:])
	}
	binary.BigEndian.PutUint32(b[4:], uint32(len(b)))
	return b
}

func container(b []byte, key *rsa.PrivateKey) ([]byte, error) {
	out := make([]byte, hdrLen, hdrLen+len(b))
	copy(out, "XMNP")
	binary.BigEndian.PutUint32(out[4:], 1)
	sum := sha1.Sum(b)
	copy(out[0x24:], sum[:])
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA1, sum[:])
	if err != nil {
		return nil, err
	}
	copy(out[0x38:], sig)
	return append(out, b...), nil
}

// BuildManifest returns the signed manifest. With autoRev, the revision is the first 8 bytes
// (big-endian) of the SHA1 of the body built with revision 0, so it changes exactly when the
// file set changes; otherwise rev is used as given.
func BuildManifest(files []File, rev uint64, autoRev bool, key *rsa.PrivateKey) ([]byte, uint64, error) {
	if key.N.BitLen() != 2048 {
		return nil, 0, fmt.Errorf("signing key must be RSA-2048")
	}
	if autoRev {
		sum := sha1.Sum(body(0, files))
		rev = binary.BigEndian.Uint64(sum[:8])
	}
	m, err := container(body(rev, files), key)
	return m, rev, err
}

// Server serves the manifest and the module files.
type Server struct {
	Log      *slog.Logger
	Manifest []byte
	Rev      uint64
	Files    []File
}

// RevHex is the SOMRevision hive value: 16 hex digits.
func (s *Server) RevHex() string { return fmt.Sprintf("%016X", s.Rev) }

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/healthz" {
		w.Write([]byte("ok"))
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.NotFound(w, r)
		return
	}
	p, ok := strings.CutPrefix(r.URL.Path, titlePath)
	if !ok {
		http.NotFound(w, r)
		return
	}
	parts := strings.Split(p, "/")
	var data []byte
	switch {
	case len(parts) == 3 && parts[0] == "manifests" && parts[2] == "system.online.manifest":
		if v, err := strconv.ParseUint(parts[1], 16, 64); err != nil || v != s.Rev {
			s.Log.Warn("manifest revision mismatch", "asked", parts[1], "current", s.RevHex())
		}
		data = s.Manifest
	case len(parts) == 3 && parts[0] == "files":
		for _, f := range s.Files {
			if strings.EqualFold(parts[1], URLID(f.ID)) && parts[2] == f.Name {
				data = f.Data
			}
		}
	}
	if data == nil {
		s.Log.Info("not found", "path", r.URL.Path)
		http.NotFound(w, r)
		return
	}
	sum := sha1.Sum(data)
	s.Log.Info("serve", "path", r.URL.Path, "floken", r.Header.Get("X-FLOKEN") != "")
	w.Header().Set("ETag", `"`+hex.EncodeToString(sum[:])+`"`)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeContent(w, r, "", zeroTime, bytes.NewReader(data))
}

var zeroTime time.Time
