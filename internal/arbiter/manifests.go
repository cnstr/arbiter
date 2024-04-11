package arbiter

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-co-op/gocron/v2"
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

type Payload struct {
	AssignedManifests []Manifest
	FastTrack         bool
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

// Is there a better way to pass state into gocron?
var FastTrack bool = false

func ScheduleManifests(PeerStore *PeerStore) (chan bool, gocron.Job) {
	stop := make(chan bool)
	sched, err := gocron.NewScheduler()

	if err != nil {
		log.Println("[scheduler] Could not create scheduler:", err)
		return nil, nil
	}

	job, err := sched.NewJob(
		gocron.CronJob("0 * * * *", false),
		gocron.NewTask(func() {
			manifests := fetchManifests()
			if manifests == nil {
				log.Println("[scheduler] Could not fetch manifests")
			}

			if FastTrack {
				log.Println("[scheduler] Fast tracking manifests")
			}

			count := 0
			for _, peer := range PeerStore.PeerArray {
				if peer.State != Ready {
					continue
				}

				if peer.Connection == nil {
					log.Println("[scheduler] Peer has no connection")
					continue
				}

				peer.Connection.WriteJSON(&Payload{
					AssignedManifests: peer.AssignedManifests,
					FastTrack:         FastTrack,
				})

				log.Println("[scheduler] Sent manifests to:", peer.Id)
				count++
			}

			FastTrack = false
			if count == 0 {
				// TODO: Sentry Error here, this is probably bad
				log.Println("[scheduler] No peers to send manifests to")
			}
		}),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)

	if err != nil {
		log.Println("[scheduler] Could not create job:", err)
		return nil, nil
	}

	log.Println("[scheduler] Scheduled job:", job.ID())

	go func() {
		sched.Start()
		select {
		case <-stop:
			err := sched.Shutdown()
			if err != nil {
				log.Println("[scheduler] Could not shutdown scheduler:", err)
			}
		}
	}()

	return stop, job
}
