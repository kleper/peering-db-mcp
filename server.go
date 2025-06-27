package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env if present
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/locations/", locationsHandler)
	http.HandleFunc("/openapi.json", openapiHandler)

	log.Printf("HTTP server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func locationsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ASN from path
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "ASN required", http.StatusBadRequest)
		return
	}
	asn, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, "invalid ASN", http.StatusBadRequest)
		return
	}

	netID, err := getNetworkID(asn)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting network ID: %v", err), http.StatusInternalServerError)
		return
	}

	locations, err := getPeeringLocations(netID)
	if err != nil {
		http.Error(w, fmt.Sprintf("error getting locations: %v", err), http.StatusInternalServerError)
		return
	}

	var response strings.Builder
	response.WriteString(fmt.Sprintf("🌐 IX Locations for AS%d:\n\n", asn))
	if len(locations) == 0 {
		response.WriteString("No IX locations found for this ASN.")
	} else {
		for i, loc := range locations {
			status := "🟢 Operational"
			if !loc.Status {
				status = "🔴 Not Operational"
			}
			response.WriteString(fmt.Sprintf("%d. %s\n   📍 %s\n   Status: %s\n\n", i+1, loc.IXName, loc.City, status))
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response.String()))
}

func openapiHandler(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile("openapi.json")
	if err != nil {
		http.Error(w, "spec not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
