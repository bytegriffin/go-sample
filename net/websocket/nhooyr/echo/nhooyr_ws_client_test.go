package echo

import (
	"context"
	"log"
	"testing"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

func TestNhooyrWSClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, "ws://localhost:8080/echo", &websocket.DialOptions{
		Subprotocols:         []string{"echo"},
		CompressionMode:      websocket.CompressionNoContextTakeover,
		CompressionThreshold: 128,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(websocket.StatusInternalError, "Client端内部出错。")

	for i := 0; i < 5; i++ {
		err = wsjson.Write(ctx, conn, "hello world")
		if err != nil {
			t.Fatal(err)
		}

		var v interface{}
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("Client %v count read msg %s ", i, v)
	}

	conn.Close(websocket.StatusNormalClosure, "")
}
