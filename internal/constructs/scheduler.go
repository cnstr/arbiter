package constructs

import (
	"encoding/json"
	"log"
	"os"
)

type Scheduler struct {
	peers      map[string]*Peer
	broadcast  chan []byte
	register   chan *Peer
	unregister chan *Peer
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		broadcast:  make(chan []byte),
		register:   make(chan *Peer),
		unregister: make(chan *Peer),
		peers:      NewOrLoadPeerStore(),
	}
}

func (s *Scheduler) Listen() {
	for {
		select {
		// Load into active peer list
		case peer := <-s.register:
			s.peers[peer.ipAddr] = peer
			flushPeerStore(s)

		// Remove from active peer list
		case peer := <-s.unregister:
			if _, ok := s.peers[peer.ipAddr]; ok {
				delete(s.peers, peer.ipAddr)
				flushPeerStore(s)
				close(peer.send)
			}

		// Share all messages between peers
		case message := <-s.broadcast:
			for ip := range s.peers {
				peer := s.peers[ip]
				if peer == nil || !peer.isValidated {
					continue
				}

				select {
				case peer.send <- message:
				default:
					if peer.send != nil {
						close(peer.send)
						delete(s.peers, ip)
						flushPeerStore(s)
					}
				}
			}
		}
	}
}

func NewOrLoadPeerStore() map[string]*Peer {
	ps := make(map[string]*Peer)

	path := os.Getenv("PEER_STORE_PATH")
	if len(path) == 0 {
		path = "peers.json"
	}

	log.Println("[store] Attempting to read peer file:", path)

	// Load the peer store from disk
	data, err := os.ReadFile("peers.json")
	if err != nil {
		log.Println("[store] Could not read peer store from disk:", err)
		log.Println("[store] Creating new peer store")
		return ps
	}

	// Unmarshal the JSON into the map
	err = json.Unmarshal(data, &ps)
	if err != nil {
		log.Println("[store] Could not unmarshal peer store from JSON:", err)
		return ps
	}

	log.Println("[store] Loaded peer store from disk")
	return ps
}

func flushPeerStore(ps *Scheduler) {
	// Marshal the map to JSON
	data, err := json.Marshal(ps.peers)
	if err != nil {
		// TODO: Fatal or sentry?
		log.Println("[store] Could not marshal peer store to JSON:", err)
		return
	}

	// Write the JSON to disk
	err = os.WriteFile("peers.json", data, 0644)
	if err != nil {
		log.Println("[store] Could not write peer store to disk:", err)
		return
	}
}
