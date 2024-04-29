package main

import (
	// "fmt"
	"fmt"
	"log"
	// "net"
	// "html"
	"net/http"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Mid1")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		http.StripPrefix("/app/", next).ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(http.FileServer(http.Dir("./")))	
	mux.Handle("/app/*", corsMux)
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	// http.HandleFunc("/greeting", greeting)
	
	log.Println("Starting server....")
	http.ListenAndServe(":1224", mux)
}