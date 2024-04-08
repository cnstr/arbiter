package arbiter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cnstr/arbiter/v2/internal/utils"
)

func StartServer() {
	// All of these exit as fatal if they fail
	keyPath, certPath := utils.LoadRootPaths()
	caCertPool := utils.LoadTrustChain(keyPath, certPath)
	tlsConfig := utils.GetTlsConfig(caCertPool)

	peerStore := NewPeerStore()
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the TLS certificate from the request
			cert := r.TLS.PeerCertificates[0]
			// print cert hash
			log.Println("[http] Received request from:", cert.SerialNumber)

			// Return the peers as a JSON response
			peers := peerStore.GetAllPeers()
			data, err := json.Marshal(peers)
			if err != nil {
				log.Println("[http] Could not marshal peers:", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		}),
	}

	peerStore.AddPeer(&Peer{
		SerialNumber: "test1",
	})

	peerStore.AddPeer(&Peer{
		SerialNumber: "test2",
	})

	stop := ScheduleManifests()
	log.Println("[http] Server started on:", server.Addr)
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
	stop <- true
}
