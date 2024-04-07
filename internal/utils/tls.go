package utils

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"strings"
)

type TlsEnv struct {
	KeyPath  string
	CertPath string
	CertPool *x509.CertPool
}

func LoadTlsEnv() TlsEnv {
	var keyPath = os.Getenv("ROOT_CA_KEY")
	var certPath = os.Getenv("ROOT_CA_CERT")

	if strings.TrimSpace(keyPath) == "" {
		log.Fatal("[arbiter] Missing ROOT_CA_KEY environment variable")
		panic("Missing ROOT_CA_KEY environment variable")
	}

	if strings.TrimSpace(certPath) == "" {
		log.Fatal("[arbiter] Missing ROOT_CA_CERT environment variable")
		panic("Missing ROOT_CA_CERT environment variable")
	}

	// Check if the files exist
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		log.Fatal("[arbiter] Key file does not exist")
		panic("Key file does not exist")
	}

	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		log.Fatal("[arbiter] Cert file does not exist")
		panic("Cert file does not exist")
	}

	// Read and decode the certificate from PEM
	file, err := os.ReadFile(certPath)
	if err != nil {
		log.Fatal("[arbiter] Could not read cert file")
		panic("Could not read cert file")
	}

	block, _ := pem.Decode(file)
	if block == nil {
		log.Fatal("[arbiter] Could not decode cert file")
		panic("Could not decode cert file")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatal("[arbiter] Could not parse cert file")
		panic("Could not parse cert file")
	}

	// Create a certificate pool
	certPool := x509.NewCertPool()
	certPool.AddCert(cert)

	log.Println("[arbiter] TLS key and cert loaded successfully")
	return TlsEnv{
		KeyPath:  keyPath,
		CertPath: certPath,
		CertPool: certPool,
	}
}
