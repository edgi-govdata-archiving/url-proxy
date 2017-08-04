package main

import (
	"fmt"
	"io"
	"net/http"
)

// DiffHandler routes requests to the right func
func ProxyHandler(url string) func(w http.ResponseWriter, r *http.Request) {
	log.Println(url)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "OPTIONS":
			EmptyOkHandler(w, r)
		case "GET":
			resp, err := http.Get(url)
			if err != nil {
				w.Write([]byte(fmt.Sprintf("error: %s", err.Error())))
				return
			}
			if resp.StatusCode != http.StatusOK {
				w.Write([]byte(fmt.Sprintf("invalid response code: %d", resp.StatusCode)))
				return
			}
			defer resp.Body.Close()
			io.Copy(w, resp.Body)
		default:
			NotFoundHandler(w, r)
		}
	}
}

// HealthCheckHandler is a basic "hey I'm fine" for load balancers
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{ "status" : 200 }`))
}

// EmptyOkHandler is an empty 200 response, often used
// for OPTIONS requests that responds with headers set in addCorsHeaders
func EmptyOkHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// CertbotHandler pipes the certbot response for manual certificate generation
func CertbotHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, cfg.CertbotResponse)
}

// NotFoundHandler is a basic JSON 404
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{ "status" :  "not found" }`))
}
