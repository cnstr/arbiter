package constructs

import (
	"encoding/json"
)

func SayCommand(payload json.RawMessage, mt int, p *Peer) error {
	p.sched.broadcast <- payload
	return nil
}
