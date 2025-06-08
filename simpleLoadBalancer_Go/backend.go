package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from backend 3031!")
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Backend running at :3031")
	log.Fatal(http.ListenAndServe(":3031", nil))
}
