package arbiter

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/cnstr/arbiter/v2/internal/utils"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{}
var tlsEnv = utils.LoadTlsEnv()

type Message struct {
	PeerKey string
	Command string
	Payload string
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("[ws] Read error, closing:", err)
			break
		}

		var payload Message
		err = json.Unmarshal(message, &payload)
		if err != nil {
			log.Println("[ws] Could not unmarshal message:", err)
			err := utils.WriteResponse(false, "Could not unmarshal message", mt, c)
			if err != nil {
				// Fatal for connection, bail
				break
			}

			continue
		}

		// Verify the TLS certificate that was sent
		status := VerifyTls(payload.PeerKey, tlsEnv.CertPool)

		if !status {
			err := utils.WriteResponse(false, "Could not verify TLS certificate", mt, c)
			if err != nil {
				// Fatal for connection, bail
				break
			}

			continue
		}

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func StartServer() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)

	log.Println("[ws] Server started on:", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
