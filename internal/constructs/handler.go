package constructs

import (
	"log"
	"net/http"
)

func PeerHandler(sched *Scheduler, w http.ResponseWriter, r *http.Request) {
	ipAddr := r.RemoteAddr
	version := r.Header.Get("X-Canister-Version")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Peer{
		sched:        sched,
		conn:         conn,
		send:         make(chan []byte, 256),
		ipAddr:       ipAddr,
		cnstrVersion: version,
		isValidated:  false,
	}

	client.sched.register <- client
	go client.writeCycle()
	go client.readCycle()
}
