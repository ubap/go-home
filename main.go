package main

import (
	"encoding/json"
	"flag"
	"fmt"
	auth2 "goHome/auth"
	"goHome/reku"
	"log"
	"net/http"
)

type APIServer struct {
	controller reku.RecuperatorController
	auth       *auth2.UserManager
}

const dbFile = "/data/users.json"

func main() {
	useMock := flag.Bool("mock", false, "Użyj mockowego kontrolera zamiast prawdziwego")
	flag.Parse() // Parsujemy flagi podane przy uruchomieniu

	var ctrl reku.RecuperatorController

	if *useMock {
		ctrl = reku.NewMockRecuperator()
	} else {
		ctrl = reku.NewKomfventRecuperator()
	}

	authManager, err := auth2.NewBasicAuthManager(dbFile)
	if err != nil {
		log.Fatalf("Failed to initialize auth manager: %v", err)
	}

	server := &APIServer{
		controller: ctrl,
		auth:       authManager,
	}

	http.Handle("/", server.auth.BasicAuthMiddleware()(http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/api/status", server.auth.BasicAuthMiddleware()(http.HandlerFunc(server.statusHandler)).ServeHTTP)
	http.HandleFunc("/api/moc", server.auth.BasicAuthMiddleware()(http.HandlerFunc(server.setPowerHandler)).ServeHTTP)

	addr := ":8080"
	log.Printf("Server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func (s *APIServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	status, err := s.controller.GetStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (s *APIServer) setPowerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Niedozwolona metoda", http.StatusMethodNotAllowed)
		return
	}

	username := r.Header.Get("X-Auth-User")
	role := r.Header.Get("X-Auth-Role")

	var req struct {
		Moc int `json:"moc"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("Użytkownik %s (%s) zmienia moc na %d%%\n", username, role, req.Moc)

	if err := s.controller.SetExtractAndSupplyFanSpeed(req.Moc, req.Moc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"user":   username,
		"action": fmt.Sprintf("Moc ustawiona na %d%%", req.Moc),
	})
}
