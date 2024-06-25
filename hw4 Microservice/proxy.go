package main

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()

	router.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		service1URL := "http://localhost:8080"
		service2URL := "http://localhost:8081"

		targetURL := service1URL
		if rand.Intn(2) == 0 {
			targetURL = service2URL
		}

		target, err := url.Parse(targetURL)
		if err != nil {
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(target)

		proxy.ServeHTTP(w, r)
	})

	log.Println("Starting proxy server on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
