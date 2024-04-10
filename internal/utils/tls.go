package utils

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"strings"
)

func LoadRootPaths() (string, string) {
	var keyPath = pickOr("ROOT_KEY", "/etc/ssl/private/root.key")
	var certPath = pickOr("ROOT_CERT", "/etc/ssl/certs/root.crt")

	return keyPath, certPath
}

func LoadTrustChain(keyPath string, certPath string) x509.CertPool {
	caCert, err := os.ReadFile(certPath)
	if err != nil {
		log.Fatal("[arbiter] Could not load CA certificate:", err)
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatal("[arbiter] Could not append CA certificate to pool")
	}

	log.Println("[arbiter] Loaded CA certificate from:", certPath)
	return *caCertPool
}

func pickOr(env string, name string) string {
	out := os.Getenv(env)
	if strings.TrimSpace(out) == "" {
		return name
	}

	return out
}

func GetTlsConfig(caCertPool x509.CertPool) *tls.Config {
	return &tls.Config{
		ClientCAs:  &caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
}
