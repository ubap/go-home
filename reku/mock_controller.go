package reku

import (
	"log"
)

// MockRecuperator to nasza implementacja testowa (mock).
// Przechowuje ona swój stan w pamięci.
type MockRecuperator struct {
	// Przechowujemy aktualny, symulowany status.
	currentStatus Status
}

// tworzenie zmiennej *MockRecuperator o wartosci nil. To jest tylko w compile time.
var _ RecuperatorController = (*MockRecuperator)(nil)

// NewMockRecuperator tworzy nową instancję naszego mocka z domyślnymi wartościami.
func NewMockRecuperator() *MockRecuperator {
	log.Println("Używam mockowego kontrolera rekuperatora.")
	return &MockRecuperator{
		currentStatus: Status{
			Power:                 45,
			Recovery:              88,
			ExtractAirTemperature: "21.5 °C",
			SupplyAirTemperature:  "19.8 °C",
			OutdoorAirTemperature: "12.0 °C",
		},
	}
}

// GetStatus dla mocka po prostu zwraca stan przechowywany w pamięci.
func (m *MockRecuperator) GetStatus() (Status, error) {
	log.Printf("[MOCK] Odczytano status: %+v", m.currentStatus)
	// Zwracamy skopiowany stan i nil jako błąd (mock zawsze działa).
	return m.currentStatus, nil
}

func (m *MockRecuperator) SetExtractAndSupplyFanSpeed(extractFanSpeed int, supplyFanSpeed int) error {
	log.Printf("[MOCK] Ustawiam moc na: %d%%", extractFanSpeed)
	m.currentStatus.Power = extractFanSpeed
	return nil // Sukces
}

// SetMode dla mocka symuluje zmianę trybu.
func (m *MockRecuperator) SetMode(mode string) error {
	log.Printf("[MOCK] Ustawiam tryb na: %s", mode)
	// Tutaj moglibyśmy dodać logikę zmiany temperatur w zależności od trybu.
	return nil // Sukces
}
