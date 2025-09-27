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
			Power:       45,
			Recovery:    88,
			TempExtract: 21.5,
			TempInlet:   19.8,
			TempOutside: 12.0,
		},
	}
}

// GetStatus dla mocka po prostu zwraca stan przechowywany w pamięci.
func (m *MockRecuperator) GetStatus() (Status, error) {
	log.Printf("[MOCK] Odczytano status: %+v", m.currentStatus)
	// Zwracamy skopiowany stan i nil jako błąd (mock zawsze działa).
	return m.currentStatus, nil
}

// SetPower dla mocka zmienia wartość mocy w jego wewnętrznym stanie.
func (m *MockRecuperator) SetPower(power int) error {
	log.Printf("[MOCK] Ustawiam moc na: %d%%", power)
	m.currentStatus.Power = power
	return nil // Sukces
}

// SetMode dla mocka symuluje zmianę trybu.
func (m *MockRecuperator) SetMode(mode string) error {
	log.Printf("[MOCK] Ustawiam tryb na: %s", mode)
	// Tutaj moglibyśmy dodać logikę zmiany temperatur w zależności od trybu.
	return nil // Sukces
}
