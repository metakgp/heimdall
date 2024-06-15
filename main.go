package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/likexian/whois"
	"github.com/rs/cors"
)

func handleCampusCheck(res http.ResponseWriter, req *http.Request) {
	clientIP := req.Header.Get("X-Forwarded-For")
	whoisResponse, err := whois.Whois(clientIP)
	if err != nil {
		fmt.Println(err)
	}

	// Define a regular expression pattern to match the netname
	pattern := `netname:\s+(.*)`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find the netname using the regular expression
	match := re.FindStringSubmatch(whoisResponse)

	response := make(map[string]bool)

	if len(match) >= 2 {
		netname := match[1]
		if netname == "IITKGP-IN" {
			response["is_inside_kgp"] = true
		} else {
			response["is_inside_kgp"] = false
		}
	} else {
		fmt.Println("Netname not found in the whois response.")
		response["is_inside_kgp"] = false
	}

	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	jsonResp, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	res.Write(jsonResp)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/campus_check", handleCampusCheck)

	handler := cors.AllowAll().Handler(mux)
	err := http.ListenAndServe(":3333", handler)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
