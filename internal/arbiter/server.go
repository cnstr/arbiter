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
	wsHandler := GetHandlerFunc(peerStore)

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/connect" {
				wsHandler(w, r)
				return
			}

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

	stop := ScheduleManifests()
	log.Println("[http] Server started on:", server.Addr)
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
	stop <- true
}
