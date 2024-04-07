package commands

import (
	"encoding/json"

	"github.com/cnstr/arbiter/v2/internal/utils"
	"github.com/gorilla/websocket"
)

var Peers = utils.NewOrLoadPeerStore()

func Peer(payload json.RawMessage, mt int, c *websocket.Conn) error {
	// Try to unmarshal the payload into a PeerData struct
	var peerData utils.PeerData
	err := json.Unmarshal(payload, &peerData)
	if err != nil {
		return utils.WriteResponse(false, "Invalid payload", mt, c)
	}

	// Check if the peer is already in the store
	_, ok := Peers.GetPeer(peerData.IpAddr)
	if ok {
		return utils.WriteResponse(true, "Peer already exists", mt, c)
	}

	Peers.AddPeer(peerData.IpAddr, peerData)
	return utils.WriteResponse(true, "Peer added", mt, c)
}
