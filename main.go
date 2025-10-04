package main

import (
	"encoding/json"
	"flag" // Użyjemy flag do przełączania trybu mock/real
	"fmt"
	"goHome/reku"
	"log"
	"net/http"
)

// APIServer przechowuje zależności, takie jak nasz kontroler.
// To jest tzw. wstrzykiwanie zależności (dependency injection).
type APIServer struct {
	controller reku.RecuperatorController
	auth       *BasicAuthManager
}

func main() {
	// Definiujemy flagę startową, aby móc wybrać implementację
	useMock := flag.Bool("mock", false, "Użyj mockowego kontrolera zamiast prawdziwego")
	flag.Parse() // Parsujemy flagi podane przy uruchomieniu

	// Tworzymy zmienną typu naszego interfejsu
	var ctrl reku.RecuperatorController

	// W zależności od flagi, tworzymy instancję mocka ALBO prawdziwego kontrolera
	if *useMock {
		ctrl = reku.NewMockRecuperator()
	} else {
		ctrl = reku.NewKomfventRecuperator()
	}

	// Inicjalizujemy system autoryzacji
	auth := NewBasicAuthManager()

	// Tworzymy instancję naszego serwera i "wstrzykujemy" mu kontroler
	server := &APIServer{
		controller: ctrl,
		auth:       auth,
	}

	// Rejestrujemy handlery z autoryzacją
	// Strona główna - dostępna dla wszystkich zalogowanych użytkowników
	http.Handle("/", server.auth.BasicAuthMiddleware("user")(http.FileServer(http.Dir("./static"))))

	// API endpoints - różne poziomy dostępu
	http.HandleFunc("/api/status", server.auth.BasicAuthMiddleware("user")(http.HandlerFunc(server.statusHandler)).ServeHTTP)
	http.HandleFunc("/api/moc", server.auth.BasicAuthMiddleware("user")(http.HandlerFunc(server.setPowerHandler)).ServeHTTP)

	// Admin endpoint - tylko dla administratorów
	http.HandleFunc("/api/admin/users", server.auth.BasicAuthMiddleware("admin")(http.HandlerFunc(server.adminHandler)).ServeHTTP)

	addr := ":8080"
	log.Printf("Server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

// statusHandler jest teraz metodą APIServer, więc ma dostęp do kontrolera
func (s *APIServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	status, err := s.controller.GetStatus()
	if err != nil {
		// Jeśli kontroler zwróci błąd, wysyłamy go do klienta
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// setPowerHandler również jest metodą APIServer
func (s *APIServer) setPowerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Niedozwolona metoda", http.StatusMethodNotAllowed)
		return
	}

	// Logowanie akcji użytkownika
	username := r.Header.Get("X-Auth-User")
	role := r.Header.Get("X-Auth-Role")

	// Używamy tymczasowej struktury, bo nie potrzebujemy jej nigdzie indziej
	var req struct {
		Moc int `json:"moc"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Logowanie zmiany mocy
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

// adminHandler - endpoint tylko dla administratorów
func (s *APIServer) adminHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("X-Auth-User")

	// Przykładowe informacje administracyjne
	adminInfo := map[string]interface{}{
		"message": "Panel administratora",
		"user":    username,
		"users":   []string{"admin", "family"},
		"system_info": map[string]interface{}{
			"uptime":  "System działa",
			"version": "1.0.0",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(adminInfo)
}
