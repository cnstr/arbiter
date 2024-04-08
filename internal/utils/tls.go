package utils

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"strings"
)

func LoadRootPaths() (string, string) {
	var keyPath = ensureValidEnv("CANISTER_ROOT_KEY")
	var certPath = ensureValidEnv("CANISTER_ROOT_CERT")

	return keyPath, certPath
}

func LoadTrustChain(keyPath string, certPath string) x509.CertPool {
	caCert, err := os.ReadFile(certPath)
	if err != nil {
		log.Fatal("[arbiter] Could not load CA key pair:", err)
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatal("[arbiter] Could not append CA certificate to pool")
	}

	return *caCertPool
}

func ensureValidEnv(name string) string {
	out := os.Getenv(name)
	if strings.TrimSpace(out) == "" {
		log.Fatal("[arbiter] Missing ", name, " environment variable")
	}

	return out
}

func GetTlsConfig(caCertPool x509.CertPool) *tls.Config {
	return &tls.Config{
		ClientCAs:  &caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
}
