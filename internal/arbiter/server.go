package arbiter

import (
	"flag"
	"log"
	"net/http"

	"github.com/cnstr/arbiter/v2/internal/constructs"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func StartServer() {
	scheduler := constructs.NewScheduler()
	go scheduler.Listen()

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		constructs.PeerHandler(scheduler, w, r)
	})

	log.Println("[ws] Server started on:", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
