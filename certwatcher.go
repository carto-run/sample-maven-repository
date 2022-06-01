package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

// thx cartographer-conventions!
// https://github.com/vmware-tanzu/cartographer-conventions/blob/9634b26dfb5fe7bbdbfa4f0607180c719b589493/webhook/server.go
//
type certWatcher struct {
	CrtFile string
	KeyFile string

	m       sync.Mutex
	keyPair *tls.Certificate
}

func (w *certWatcher) Watch(ctx context.Context) error {
	// refresh the certs periodically even if we miss a fs event
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := w.Load(ctx); err != nil {
					log.Println("error loading TLS key pair", err)
				}
			}
		}
	}()

	<-ctx.Done()
	return nil
}

func (w *certWatcher) Load(ctx context.Context) error {
	w.m.Lock()
	defer w.m.Unlock()

	crt, err := ioutil.ReadFile(w.CrtFile)
	if err != nil {
		return err
	}
	key, err := ioutil.ReadFile(w.KeyFile)
	if err != nil {
		return err
	}
	keyPair, err := tls.X509KeyPair(crt, key)
	if err != nil {
		return err
	}
	leaf := keyPair.Leaf
	if leaf == nil {
		leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
		if err != nil {
			return err
		}
	}

	w.keyPair = &keyPair
	log.Println("loaded TLS key pair", "not-after", leaf.NotAfter)
	return nil
}

func (w *certWatcher) GetCertificate() *tls.Certificate {
	w.m.Lock()
	defer w.m.Unlock()

	return w.keyPair
}
