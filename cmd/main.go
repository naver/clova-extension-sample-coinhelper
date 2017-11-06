package main

import (
	"coinHelper/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/currency", handler.ServeHTTP)
	http.HandleFunc("/health_check", handler.HealthCheck)
	http.HandleFunc("/monitor/l7check", handler.HealthCheck)

	http.ListenAndServe(":10680", nil)
}
