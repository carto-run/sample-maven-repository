package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

var (
	dir  = flag.String("d", ".", "directory to serve")
	addr = flag.String("addr", ":8080", "address to bind to")

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

	http.Handle("/", handler())

	log.Println("listening on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		panic(err)
	}
}
