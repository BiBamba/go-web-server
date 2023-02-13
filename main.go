// Including Prometheus instrumentation
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var REQUESTS = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "http_requests_count",
		Help: "Number of request handled by the server.",
	},
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	REQUESTS.Inc()
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err:%v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful")
	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Address = %s", address)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	REQUESTS.Inc()
	if r.URL.Path != "/welcome" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Welcome To my first Go Server")
}

func main() {
	prometheus.MustRegister(REQUESTS)

	fileServer := http.FileServer(http.Dir("./static"))
	promHandler := promhttp.Handler()

	http.Handle("/", fileServer)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/welcome", welcomeHandler)

	// Prometheus endpoint
	http.Handle("/prometheus", promHandler)

	fmt.Printf("Starting server on port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
