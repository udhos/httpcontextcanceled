// Package main demonstrates handling context cancellation in a chi server.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	// spawn chi server

	r := chi.NewRouter()
	r.Get("/", handler)

	log.Printf("test us with: curl http://localhost:8080")
	log.Printf("then cancel with ctrl-c")

	go func() {
		err := http.ListenAndServe(":8080", r)
		log.Printf("http server exited with error: %v", err)
	}()

	select {} // block forever
}

func handler(w http.ResponseWriter, r *http.Request) {

	log.Println("server: will wait until request is cancelled")

	ctx := r.Context()
	<-ctx.Done()

	reason := ctx.Err()

	log.Printf("server: request cancelled, reason: %v", reason)

	fmt.Fprintf(w, "server: request cancelled, reason: %v\n", reason)
}
