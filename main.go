package main

import (
	"context"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	//defer cancel()

	conn, _, err := websocket.Dial(ctx, URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := wsjson.Write(ctx, conn, []string{}); err != nil {
		log.Fatal(err)
	}

	var ack AckMessage
	if err = wsjson.Read(ctx, conn, &ack); err != nil {
		log.Fatal(err)
	}

	log.Printf("connected to broker endpoint=%s version=%s", ack.EndpointUUID, ack.Version)

	stopChan := make(chan struct{}, 1)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	counter := atomic.Int64{}

	go func() {
		<-sigChan
		stopChan <- struct{}{}
		if err := conn.Close(websocket.StatusNormalClosure, ""); err != nil {
			log.Println(err)
		}
		cancel()
		log.Printf("exited with counter=%d", counter.Load())
		os.Exit(0)
	}()

	buf := []byte(event)

Loop:
	for {
		select {
		case <-stopChan:
			break Loop
		default:
			if err = conn.Write(ctx, websocket.MessageText, buf); err != nil {
				break Loop
			}
			counter.Add(1)
		}
	}

	cancel()
	log.Printf("exited with error=%s and counter=%d", err.Error(), counter.Load())
}
