package arbiter

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout:  time.Second * 5,
	EnableCompression: true,
}

func GetHandlerFunc(PeerStore *PeerStore) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cert := r.TLS.PeerCertificates[0]
		if cert == nil {
			// TODO: Sentry Error here, this isn't possible
			log.Println("[http] Got a request without a valid certificate?")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden"))
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("[http] Could not upgrade connection:", err)
			PeerStore.UpdatePeerState(cert.SerialNumber.String(), Disconnected)

			return
		}

		PeerStore.AddPeer(&Peer{
			Id:          cert.SerialNumber.String(),
			State:       NotReady,
			Connection:  conn,
			Certificate: cert,
		})

		conn.SetPingHandler(func(data string) error {
			state := State(data)

			if state != NotReady && state != Ready && state != Running {
				log.Println("[ws] Invalid ping data from:", cert.SerialNumber)
				log.Println("[ws] Got:", state)
				return nil
			}

			PeerStore.UpdatePeerState(cert.SerialNumber.String(), state)
			return nil
		})

		go func() {
			log.Println("[ws] Connected to:", cert.SerialNumber)
			defer conn.Close()

			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					log.Println("[ws] Could not read message:", err)
					PeerStore.UpdatePeerState(cert.SerialNumber.String(), Disconnected)
					break
				}

			}
		}()
	}
}
