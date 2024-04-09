package arbiter

import (
	"crypto/x509"
	"log"
)

type State string

const (
	NotReady     State = "NotReady"
	Ready        State = "Ready"
	Running      State = "Running"
	Disconnected State = "Disconnected"
)

type Peer struct {
	Id                string
	State             State
	Certificate       *x509.Certificate
	AssignedManifests []Manifest
}

type PeerStore struct {
	Manifests []Manifest
	Peers     map[string]*Peer
}

// We have a mutable list of peers globally available for the http handler
// and the scheduler which modifies the list to rebalance the manifests

func NewPeerStore() *PeerStore {
	return &PeerStore{
		Peers:     make(map[string]*Peer),
		Manifests: fetchManifests(),
	}
}

func (ps *PeerStore) AddPeer(peer *Peer) {
	ps.Peers[peer.Id] = peer
	ps.RebalancePeers()
}

func (ps *PeerStore) UpdatePeerState(Id string, state State) {
	ps.Peers[Id].State = state
	ps.RebalancePeers()
}

func (ps *PeerStore) RebalancePeers() {
	// If there are ready peers with no manifests, then rebalancing is needed
	// If there are no ready peers, then we can't rebalance
	shouldRebalance := false
	readyPeers := []*Peer{}
	for _, peer := range ps.Peers {
		if peer.State == Ready {
			readyPeers = append(readyPeers, peer)

			if len(peer.AssignedManifests) == 0 {
				shouldRebalance = true
			}
		}
	}

	peerLength := len(readyPeers)
	if peerLength == 0 {
		// Sentry if we are long running?
		log.Println("[scheduler] No peers ready to rebalance")
		return
	}

	if !shouldRebalance {
		return
	}

	// Filter the manifests into 5 groups based on ranking
	// Split the group by the number of peers and combine
	group1 := filterManifestsByRanking(ps.Manifests, 1)
	group2 := filterManifestsByRanking(ps.Manifests, 2)
	group3 := filterManifestsByRanking(ps.Manifests, 3)
	group4 := filterManifestsByRanking(ps.Manifests, 4)
	group5 := filterManifestsByRanking(ps.Manifests, 5)

	// Split the groups by the number of peers
	subgroups1 := splitManifestsAcrossCount(group1, peerLength)
	subgroups2 := splitManifestsAcrossCount(group2, peerLength)
	subgroups3 := splitManifestsAcrossCount(group3, peerLength)
	subgroups4 := splitManifestsAcrossCount(group4, peerLength)
	subgroups5 := splitManifestsAcrossCount(group5, peerLength)

	// Combine the subgroups
	m := []Manifest{}
	index := 0

	for _, peer := range readyPeers {
		m = append(m, subgroups1[index]...)
		m = append(m, subgroups2[index]...)
		m = append(m, subgroups3[index]...)
		m = append(m, subgroups4[index]...)
		m = append(m, subgroups5[index]...)
		peer.AssignedManifests = m

		index += 1
		m = []Manifest{}
	}

	log.Println("[scheduler] Rebalanced jobs across", peerLength, "peers")
}

func (ps *PeerStore) GetAllPeers() []*Peer {
	result := make([]*Peer, 0, len(ps.Peers))
	for _, peer := range ps.Peers {
		result = append(result, peer)
	}

	return result
}

func (ps *PeerStore) GetAllPeerJobs() map[string][]Manifest {
	result := make(map[string][]Manifest)
	for _, peer := range ps.Peers {
		result[peer.Id] = peer.AssignedManifests
	}

	return result
}

func filterManifestsByRanking(manifests []Manifest, ranking int) []Manifest {
	var filtered []Manifest
	for _, manifest := range manifests {
		if manifest.Quality == ranking {
			filtered = append(filtered, manifest)
		}
	}
	return filtered
}

func splitManifestsAcrossCount(manifests []Manifest, groups int) [][]Manifest {
	length := len(manifests)
	groupSize := length / groups

	var split [][]Manifest
	for i := 0; i < groups; i++ {
		start := i * groupSize
		end := start + groupSize

		// Handle the last group
		if i == groups-1 {
			end = length
		}

		split = append(split, manifests[start:end])
	}

	return split
}
