package utils

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Response struct {
	Success bool
	Message string
}

func WriteResponse(
	success bool,
	message string,
	mt int,
	c *websocket.Conn,
) error {
	res := Response{
		Success: success,
		Message: message,
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		log.Println("[ws] Could not marshal response:", err)
	}

	err = c.WriteMessage(mt, resJson)
	if err != nil {
		log.Println("[ws] Could not write response:", err)
	}

	return err
}
