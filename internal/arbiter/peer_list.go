package arbiter

import (
	"log"
)

type Peer struct {
	SerialNumber      string
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
	ps.Peers[peer.SerialNumber] = peer
	ps.RebalancePeers()
}

func (ps *PeerStore) RemovePeer(peer *Peer) {
	delete(ps.Peers, peer.SerialNumber)
	ps.RebalancePeers()
}

func (ps *PeerStore) RebalancePeers() {
	// Filter the manifests into 5 groups based on ranking
	// Split the group by the number of peers and combine
	group1 := filterManifestsByRanking(ps.Manifests, 1)
	group2 := filterManifestsByRanking(ps.Manifests, 2)
	group3 := filterManifestsByRanking(ps.Manifests, 3)
	group4 := filterManifestsByRanking(ps.Manifests, 4)
	group5 := filterManifestsByRanking(ps.Manifests, 5)

	// Split the groups by the number of peers
	subgroups1 := splitManifestsAcrossCount(group1, len(ps.Peers))
	subgroups2 := splitManifestsAcrossCount(group2, len(ps.Peers))
	subgroups3 := splitManifestsAcrossCount(group3, len(ps.Peers))
	subgroups4 := splitManifestsAcrossCount(group4, len(ps.Peers))
	subgroups5 := splitManifestsAcrossCount(group5, len(ps.Peers))

	// Combine the subgroups
	m := []Manifest{}
	index := 0

	for _, peer := range ps.Peers {
		m = append(m, subgroups1[index]...)
		m = append(m, subgroups2[index]...)
		m = append(m, subgroups3[index]...)
		m = append(m, subgroups4[index]...)
		m = append(m, subgroups5[index]...)
		peer.AssignedManifests = m

		index += 1
		m = []Manifest{}
	}

	log.Println("[scheduler] Rebalanced jobs across", len(ps.Peers), "peers")
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
		result[peer.SerialNumber] = peer.AssignedManifests
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
