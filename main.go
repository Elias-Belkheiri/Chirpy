package main

import (
	// "fmt"
	"fmt"
	"log"
	// "net"
	// "html"
	"net/http"
	"sync"
	// "strconv"
	"html/template"
)

type apiConfig struct {
	mu				sync.Mutex
	FileserverHits 	int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.mu.Lock()
		fmt.Println("Incrementing...")
		cfg.FileserverHits++
		cfg.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareGetMetrics() int {
	// cfg.mu.Lock()
	// defer cfg.mu.Unlock()
	return cfg.FileserverHits
}

func (cfg *apiConfig) middlewareResetMetrics() {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.FileserverHits = 0
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
	templ:= template.Must(template.New("temp").Parse("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited {{.FileserverHits}} times!</p></body></html>"))

	mux := http.NewServeMux()
	corsMux := middlewareCors(http.FileServer(http.Dir("./")))	

	mux.Handle("GET /app/*", counter.middlewareMetricsInc(corsMux))
	mux.Handle("GET /api/healthz", counter.middlewareMetricsInc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})))
	mux.Handle("GET /admin/count", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Execute(w, counter)
		// w.Write([]byte(count))
	}))
	mux.Handle("POST /api/reset", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter.middlewareResetMetrics()
		fmt.Println("Counter has been reset")
	}))

	mux.Handle("POST /api/users", http.HandlerFunc(addUser))

	log.Println("Starting server....")
	http.ListenAndServe(":1224", mux)
}