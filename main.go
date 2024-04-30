package main

import (
	// "fmt"
	"fmt"
	"log"
	// "net"
	// "html"
	"net/http"
	"sync"
	"strconv"
)

type apiConfig struct {
	mu				sync.Mutex
	fileserverHits 	int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.mu.Lock()
		fmt.Println("Incrementing...")
		cfg.fileserverHits++
		cfg.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareGetMetrics() int {
	// cfg.mu.Lock()
	// defer cfg.mu.Unlock()
	return cfg.fileserverHits
}

func (cfg *apiConfig) middlewareResetMetrics() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.fileserverHits = 0
}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	counter := apiConfig{sync.Mutex{}, 0}

	mux := http.NewServeMux()
	corsMux := middlewareCors(http.FileServer(http.Dir("./")))	

	mux.Handle("GET /app/*", counter.middlewareMetricsInc(corsMux))
	mux.Handle("GET /healthz", counter.middlewareMetricsInc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})))
	mux.Handle("GET /count", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := "count: " + strconv.Itoa(counter.middlewareGetMetrics()) 
		fmt.Println(count)
		w.Write([]byte(count))
	}))
	mux.Handle("POST /reset", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter.middlewareResetMetrics()
		fmt.Println("Counter has been reset")
	}))
	
	log.Println("Starting server....")
	http.ListenAndServe(":1224", mux)
}