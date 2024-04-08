package arbiter

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type Manifest struct {
	Uri       string   `json:"uri"`
	Id        string   `json:"slug"`
	Suite     string   `json:"suite"`
	Bootstrap bool     `json:"bootstrap"`
	Quality   int      `json:"ranking"`
	Component string   `json:"component"`
	Binary    string   `json:"binary"`
	Aliases   []string `json:"aliases"`
}

func fetchManifests() []Manifest {
	url := "https://source.canister.me/index-repositories.json"

	// Fetch the manifests from the source
	httpClient := http.Client{}
	resp, err := httpClient.Get(url)
	if err != nil {
		log.Println("[scheduler] Could not fetch manifests:", err)
		return nil
	}

	// Print the body of the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[scheduler] Could not read response body:", err)
	}

	// Parse the JSON response
	var manifests []Manifest
	err = json.Unmarshal(body, &manifests)
	if err != nil {
		log.Println("[scheduler] Could not parse manifests:", err)
		return nil
	}

	return manifests
}

// Spread the manifests across the peers evenly
// Ensures that each peer splits compute load
func rebalancePeers() {
}

func ScheduleManifests() chan bool {
	stop := make(chan bool)
	go func() {
		for {
			manifests := fetchManifests()
			if manifests == nil {
				log.Println("[scheduler] Could not fetch manifests")
			}

			select {
			case <-time.After(10 * time.Minute):
			case <-stop:
				return
			}
		}
	}()

	return stop
}
