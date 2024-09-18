package main

import (
	"github.com/gorilla/websocket"
	"log"
)

const URL = "ws://localhost:16666/v1/messages/json"

const event = `
{
  "@data-type": "vector",
  "data": [
    {
      "@data-type": "count",
      "data": 1
    },
    {
      "@data-type": "count",
      "data": 1
    },
    {
      "@data-type": "vector",
      "data": [
        {
          "@data-type": "string",
          "data": "Test::evt"
        },
        {
          "@data-type": "vector",
          "data": [
          ]
        }
      ]
    }
  ],
  "topic": "/simeonmiteff/test",
  "type": "data-message"
}
`

type AckMessage struct {
	ConstType    string `json:"type"`
	EndpointUUID string `json:"endpoint"`
	Version      string `json:"version"`
}

func main() {
	conn, _, err := websocket.DefaultDialer.Dial(URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := conn.WriteJSON([]string{}); err != nil {
		log.Fatal(err)
	}

	var ack AckMessage
	if err = conn.ReadJSON(&ack); err != nil {
		log.Fatal(err)
	}

	log.Printf("connected to broker endpoint=%s version=%s", ack.EndpointUUID, ack.Version)

	buf := []byte(event)

	for {
		if err = conn.WriteMessage(websocket.TextMessage, buf); err != nil {
			break
		}
	}

	log.Printf("failed with error=%s", err.Error())
}
