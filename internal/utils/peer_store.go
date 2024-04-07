package utils

import (
	"encoding/json"
	"log"
	"os"
)

// I know I should be using sqlite or something but a JSON store is easier
// I don't need to be able to query it, just read and write for assignments

type PeerData struct {
	IpAddr       string
	CnstrVersion string
}

type PeerStore struct {
	Peers map[string]PeerData
}

func NewOrLoadPeerStore() *PeerStore {
	ps := newPeerStore()

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
	err = json.Unmarshal(data, &ps.Peers)
	if err != nil {
		log.Println("[store] Could not unmarshal peer store from JSON:", err)
		return ps
	}

	log.Println("[store] Loaded peer store from disk")
	return ps
}

func newPeerStore() *PeerStore {
	return &PeerStore{
		Peers: make(map[string]PeerData),
	}
}

func (ps *PeerStore) AddPeer(peerKey string, peerData PeerData) {
	ps.Peers[peerKey] = peerData
	ps.flushToDisk()
}

func (ps *PeerStore) GetPeer(peerKey string) (PeerData, bool) {
	peer, ok := ps.Peers[peerKey]
	return peer, ok
}

func (ps *PeerStore) RemovePeer(peerKey string) {
	delete(ps.Peers, peerKey)
	ps.flushToDisk()
}

func (ps *PeerStore) GetAllPeers() map[string]PeerData {
	return ps.Peers
}

func (ps *PeerStore) flushToDisk() {
	// Marshal the map to JSON
	data, err := json.Marshal(ps.Peers)
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
