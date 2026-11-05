package som

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/binary"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestManifest(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	data := []byte("chunk")
	files := []File{{Name: "Dash.Social.lex", Data: data, ID: sha1.Sum(data)}}
	m, rev, err := BuildManifest(files, 0, true, key)
	if err != nil {
		t.Fatal(err)
	}
	if string(m[:4]) != "XMNP" || string(m[0x138:0x13c]) != "xonm" {
		t.Fatal("magic")
	}
	b := m[0x138:]
	sum := sha1.Sum(b)
	if string(m[0x24:0x38]) != string(sum[:]) || binary.BigEndian.Uint32(b[4:]) != uint32(len(b)) {
		t.Fatal("body sha1 or signed length")
	}
	if err := rsa.VerifyPKCS1v15(&key.PublicKey, crypto.SHA1, sum[:], m[0x38:0x138]); err != nil {
		t.Fatal(err)
	}
	if binary.BigEndian.Uint32(b[0x10:]) != 3 || binary.BigEndian.Uint64(b[0x18:]) != rev || binary.BigEndian.Uint32(b[0x28:]) != 0x30 || binary.BigEndian.Uint32(b[0x30:]) != 1 {
		t.Fatal("body header")
	}
	e := b[0x34:]
	off := binary.BigEndian.Uint32(e[0x1c:])
	if binary.BigEndian.Uint32(e[0x18:]) != 3 || string(e[0x38:0x4c]) != string(files[0].ID[:]) || binary.BigEndian.Uint16(b[off:]) != 16 || binary.BigEndian.Uint32(e[0x00:]) != off || string(b[off+2:off+17]) != "Dash.Social.lex" || b[off+17] != 0 {
		t.Fatal("entry")
	}
	if len(URLID(files[0].ID)) != 40 {
		t.Fatal("url id")
	}
}

func TestServe(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	data := []byte("chunk")
	files := []File{{Name: "Dash.Social.lex", Data: data, ID: sha1.Sum(data)}}
	m, rev, _ := BuildManifest(files, 0, true, key)
	ts := httptest.NewServer(&Server{Log: slog.New(slog.NewTextHandler(io.Discard, nil)), Manifest: m, Rev: rev, Files: files})
	defer ts.Close()
	get := func(p, inm string) *http.Response {
		req, _ := http.NewRequest("GET", ts.URL+p, nil)
		if inm != "" {
			req.Header.Set("If-None-Match", inm)
		}
		r, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		return r
	}
	p := "/fffe07d1/files/" + URLID(files[0].ID) + "/Dash.Social.lex"
	r := get(p, "")
	b, _ := io.ReadAll(r.Body)
	if r.StatusCode != 200 || string(b) != "chunk" || r.Header.Get("ETag") == "" {
		t.Fatal(r.StatusCode, string(b))
	}
	if r := get(p, r.Header.Get("ETag")); r.StatusCode != 304 {
		t.Fatal(r.StatusCode)
	}
	if r := get("/fffe07d1/files/00/x", ""); r.StatusCode != 404 {
		t.Fatal(r.StatusCode)
	}
	r = get("/fffe07d1/manifests/0000000000000001/system.online.manifest", "")
	b, _ = io.ReadAll(r.Body)
	if r.StatusCode != 200 || string(b) != string(m) {
		t.Fatal(r.StatusCode)
	}
}
