package sse

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"testing"
)

func formatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	return []byte(eventPayload + "\n")
}

var messageChannels = make(map[chan []byte]bool)

func sayHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	message := r.FormValue("message")

	jsonStructure, _ := json.Marshal(map[string]string{
		"name":    name,
		"message": message})

	go func() {
		for messageChannel := range messageChannels {
			messageChannel <- []byte(jsonStructure)
		}
	}()

	w.Write([]byte(`<script type="text/javascript">
        let eventListener = new EventSource("http://localhost:9000/listen")
        eventListener.onmessage = (event) => {
            let {type, data} = event
            alert("received event: ${type} with data: ${data}")
        }
    </script>`))
}

func listenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_messageChannel := make(chan []byte)
	messageChannels[_messageChannel] = true

	for {
		select {
		case _msg := <-_messageChannel:
			w.Write(formatSSE("message", string(_msg)))
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			delete(messageChannels, _messageChannel)
			return
		}
	}
}

func TestSSEServer(t *testing.T) {
	http.HandleFunc("/say", sayHandler)
	http.HandleFunc("/listen", listenHandler)

	log.Println("Running at :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
