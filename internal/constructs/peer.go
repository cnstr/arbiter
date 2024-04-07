package constructs

import (
	"encoding/json"
	"log"

	"github.com/cnstr/arbiter/v2/internal/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024,
	WriteBufferSize: 1024 * 1024,
}

type Peer struct {
	sched        *Scheduler
	conn         *websocket.Conn
	send         chan []byte
	ipAddr       string
	cnstrVersion string
	isValidated  bool
}

type Message struct {
	PeerKey string
	Payload json.RawMessage
}

func (p *Peer) readCycle() {
	// Clean up after the socket disconnects
	defer func() {
		p.sched.unregister <- p
		p.conn.Close()
	}()

	for {
		mt, message, err := p.conn.ReadMessage()
		if err != nil {
			log.Println("[ws] Read error:", err)
			break
		}

		var payload Message
		err = json.Unmarshal(message, &payload)
		if err != nil {
			log.Println("[ws] Could not unmarshal message:", err)
			err := utils.WriteResponse(false, "Could not unmarshal message", mt, p.conn)
			if err != nil {
				// Fatal for connection, bail
				break
			}

			continue
		}

		if !p.isValidated {
			status := VerifyTls(payload.PeerKey)
			if !status {
				err := utils.WriteResponse(false, "Could not verify TLS certificate", mt, p.conn)
				if err != nil {
					// Fatal for connection, bail
					break
				}

				continue
			}

			p.isValidated = true
			log.Println("[ws] Peer validated:", p.ipAddr)
		}

		data := payload.Payload
		if len(data) == 0 {
			err := utils.WriteResponse(false, "Data is empty", mt, p.conn)
			if err != nil {
				// Fatal for connection, bail
				break
			}

			continue
		}

		err = SayCommand(data, mt, p)
		if err != nil {
			log.Println("[ws] Fatal error processing command:", err)
			break
		}
	}
}

func (p *Peer) writeCycle() {
	defer func() {
		p.conn.Close()
	}()

	for {
		select {
		case message, ok := <-p.send:
			if !ok {
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := p.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)
			n := len(p.send)

			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-p.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
