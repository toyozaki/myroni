package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/toyozaki/negroni_sample/myroni"
)

type HelloResponse struct {
	Msg string `json:"msg"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	res := &HelloResponse{
		Msg: "Hello, world!",
	}
	resJSON, _ := json.Marshal(res)
	w.Header().Add("Content-Type", "text/json")
	w.Write(resJSON)
	log.Println("[helloHandler] Sent a message to client")
}

func sampleMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("[sample middleware] Before executing controller")

	next(w, r)

	log.Println("[sample middleware] A the controller was executed")
}

func main() {
	mux := http.NewServeMux()

	myroniHello := myroni.New(
		myroni.HandlerFunc(sampleMiddleware),
		myroni.Wrap(http.HandlerFunc(helloHandler)),
	)

	mux.Handle("/hello", myroniHello)

	http.ListenAndServe(":8080", mux)
}
