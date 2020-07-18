package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

func readFromUpstream(w http.ResponseWriter, r http.Request) {
	resp, err := http.Get(r.URL.Path)
	if err != nil {
		log.Printf("Error in schedulig http GET: %v\n", err)
		return
	}

	st := time.Now()
	n, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Printf("Error in reading from response: %v\n", err)
	}
	log.Printf("Wrote %v bytes in %.2f seconds\n", n, time.Since(st).Seconds())
}
