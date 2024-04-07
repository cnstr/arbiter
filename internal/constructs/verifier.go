package constructs

import (
	"crypto/x509"
	"encoding/pem"
	"log"

	"github.com/cnstr/arbiter/v2/internal/utils"
)

var TlsEnv = utils.LoadTlsEnv()

// We assume the TLS certificate is a child of the Canister CA
func VerifyTls(cert string) bool {
	certBytes := []byte(cert)
	block, _ := pem.Decode(certBytes)

	if block == nil {
		log.Println("[arbiter] Recieved invalid TLS certificate")
		return false
	}

	x509Cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Println("[arbiter] Could not parse TLS certificate")
		return false
	}

	opts := x509.VerifyOptions{
		Roots: TlsEnv.CertPool,
	}

	_, err = x509Cert.Verify(opts)
	if err != nil {
		log.Println("[arbiter] Could not verify TLS certificate")
		return false
	}

	return true
}
