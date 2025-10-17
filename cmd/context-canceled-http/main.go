// Package main demonstrates handling context cancellation in an HTTP server.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	// spawn http server
	go func() {
		err := http.ListenAndServe(":8080", http.HandlerFunc(handler))
		log.Printf("http server exited with error: %v", err)
	}()

	// create a request that will be cancelled manually
	ctx, cancelWithCause := context.WithCancelCause(context.Background())

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
	if err != nil {
		log.Fatalf("client: failed to create request: %v", err)
	}

	var wg sync.WaitGroup

	// spawn the request that will never complete by itself but we will cancel it later
	wg.Go(func() {

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("client: request failed: %v", err)
			return
		}
		defer resp.Body.Close()

		log.Printf("client: unexpected: request completed with status: %s", resp.Status)
	})

	// now cancel the request after a short delay

	log.Printf("client: will cancel the request in 1 second")

	time.Sleep(time.Second)

	log.Printf("client: cancelling the request now")

	cancelWithCause(fmt.Errorf("client: manual cancellation"))

	wg.Wait()
}

func handler(w http.ResponseWriter, r *http.Request) {

	log.Println("server: will wait until request is cancelled")

	ctx := r.Context()
	<-ctx.Done()

	reason := ctx.Err()

	log.Printf("server: request cancelled, reason: %v", reason)

	fmt.Fprintf(w, "server: request cancelled, reason: %v\n", reason)
}
