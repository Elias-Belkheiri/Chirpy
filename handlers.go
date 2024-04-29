package main

import (
	"net/http"
	"fmt"
)

func greeting(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Mid2")
	fmt.Fprintf(w, "Good Morning!")
}
