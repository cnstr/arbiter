package arbiter

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cnstr/arbiter/v2/internal/utils"
)

type Status struct {
	SchedulingStatus struct {
		LastRun string
		NextRun string
	}
	Peers []*Peer
}

// TODO: We want to turn this into a raft based peering system
// Over time I'll start implementing this and making the indexer in go
func StartServer() {
	// All of these exit as fatal if they fail
	keyPath, certPath := utils.LoadRootPaths()
	caCertPool := utils.LoadTrustChain(keyPath, certPath)
	tlsConfig := utils.GetTlsConfig(caCertPool)

	peerStore := NewPeerStore()
	wsHandler := GetHandlerFunc(peerStore)
	stop, job := ScheduleManifests(peerStore)
	if stop == nil || job == nil {
		log.Fatal("[scheduler] Could not start scheduler")
	}

	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/connect" {
				wsHandler(w, r)
				return
			}

			if r.Method == http.MethodPost {
				// Check for query parameter fast=true
				if r.URL.Query().Get("fast") == "true" {
					FastTrack = true
				}

				err := job.RunNow()
				if err != nil {
					log.Println("[scheduler] Could not run job:", err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Internal Server Error"))
					return
				}

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("OK"))
				return
			}

			// Return the peers as a JSON response
			peers := peerStore.GetAllPeers()

			lastRun, err := job.LastRun()
			if err != nil {
				log.Println("[http] Could not get last run:", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
				return
			}

			nextRun, err := job.NextRun()
			if err != nil {
				log.Println("[http] Could not get next run:", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Internal Server Error"))
				return
			}

			data, err := json.Marshal(&Status{
				SchedulingStatus: struct {
					LastRun string
					NextRun string
				}{
					LastRun: lastRun.Format("2006-01-02 15:04:05"),
					NextRun: nextRun.Format("2006-01-02 15:04:05"),
				},
				Peers: peers,
			})
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

	log.Println("[http] Server started on:", server.Addr)
	server.ListenAndServeTLS(certPath, keyPath)
	stop <- true
}
