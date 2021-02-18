package main

import (
	"log"
	"net/http"
	"testing"
)

func TestWebAssemblyHttpServer(t *testing.T) {
	log.Fatal(http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`))))
}
