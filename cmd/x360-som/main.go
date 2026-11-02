// Command x360-som serves the Xbox 360 system online manifest and dash modules (see internal/som).
package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/openpak/x360-som/internal/som"
)

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func loadKey(path string) (*rsa.PrivateKey, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	blk, _ := pem.Decode(b)
	if blk == nil {
		return nil, errors.New("signing key is not PEM")
	}
	k, err := x509.ParsePKCS8PrivateKey(blk.Bytes)
	if err != nil {
		return x509.ParsePKCS1PrivateKey(blk.Bytes)
	}
	rk, ok := k.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("signing key is not RSA")
	}
	return rk, nil
}

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	fatal := func(msg string, err error) { log.Error(msg, "err", err); os.Exit(1) }
	key, err := loadKey(os.Getenv("X360_SOM_SIGNING_KEY"))
	if err != nil {
		fatal("signing key", err)
	}
	files, err := som.LoadDir(env("X360_SOM_MODULES", "./modules"))
	if err != nil {
		fatal("modules", err)
	}
	var rev uint64
	rs := os.Getenv("X360_SOM_REVISION")
	if rs != "" {
		if rev, err = strconv.ParseUint(rs, 10, 64); err != nil {
			fatal("X360_SOM_REVISION", err)
		}
	}
	m, rev, err := som.BuildManifest(files, rev, rs == "", key)
	if err != nil {
		fatal("manifest", err)
	}
	s := &som.Server{Log: log, Manifest: m, Rev: rev, Files: files}
	log.Info("manifest built", "SOMRevision", s.RevHex(), "files", len(files), "bytes", len(m))

	addr := env("X360_SOM_ADDR", ":10045")
	srv := &http.Server{Addr: addr, Handler: s, ReadHeaderTimeout: 10 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		sc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(sc)
	}()
	log.Info("x360-som listening", "addr", addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fatal("serve", err)
	}
}
