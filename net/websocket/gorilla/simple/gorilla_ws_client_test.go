package simple

import (
	"log"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestGorillaWSClient(t *testing.T) {

	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:8080",
		Path:   "/echo",
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer c.Close()

	//支持binary、string、
	sendMsg := []byte("hello")
	err = c.WriteMessage(websocket.TextMessage, sendMsg)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("client send: %s ", sendMsg)

	_, message, err := c.ReadMessage()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("client receive: %s \n", message)

}
