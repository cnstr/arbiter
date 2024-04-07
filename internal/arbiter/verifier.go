package arbiter

import (
	"crypto/x509"
	"encoding/pem"
	"log"
)

// We assume the TLS certificate is a child of the Canister CA
func VerifyTls(cert string, pool *x509.CertPool) bool {
	certBytes := []byte(cert)
	block, _ := pem.Decode(certBytes)

	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Println("[arbiter] Could not parse TLS certificate")
		return false
	}

	opts := x509.VerifyOptions{
		Roots: pool,
	}

	_, err = x509Cert.Verify(opts)
	if err != nil {
		log.Println("[arbiter] Could not verify TLS certificate")
		return false
	}

	return true
}
