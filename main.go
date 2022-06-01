package main

import (
	"context"
	"crypto/tls"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

var (
	httpAddr  = flag.String("http-addr", ":8080", "address to bind to")
	httpsAddr = flag.String("https-addr", ":8443", "address to bind to")
	cert      = flag.String("cert", "/cert/crt.pem", "public cert filepath")
	key       = flag.String("key", "/cert/key.pem", "private key filepath")

	//go:embed content
	content embed.FS
)

func handler() http.Handler {
	fsys := fs.FS(content)

	html, err := fs.Sub(fsys, "content")
	if err != nil {
		panic(fmt.Errorf("fs sub: %w", err))
	}

	return http.FileServer(http.FS(html))
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := run(ctx); err != nil {
		panic(err)
	}
}

func run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return serveHTTP(*httpAddr, handler())
	})

	g.Go(func() error {
		return serveHTTPS(ctx, *httpsAddr, *cert, *key, handler())
	})

	return g.Wait()
}

func serveHTTP(addr string, handler http.Handler) error {
	log.Println("serving http", addr)

	return (&http.Server{
		Addr:    addr,
		Handler: handler,
	}).ListenAndServe()
}

func serveHTTPS(
	ctx context.Context,
	addr, cert, key string,
	handler http.Handler,
) error {
	log.Println("serving https", addr, cert, key)

	watcher := certWatcher{
		CrtFile: cert,
		KeyFile: key,
	}

	if err := watcher.Load(ctx); err != nil {
		return err
	}

	go watcher.Watch(ctx)

	getter := func(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
		cert := watcher.GetCertificate()
		return cert, nil
	}

	return (&http.Server{
		Addr:    addr,
		Handler: handler,
		TLSConfig: &tls.Config{
			GetCertificate:           getter,
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS13,
		},
	}).ListenAndServeTLS(cert, key)
}
