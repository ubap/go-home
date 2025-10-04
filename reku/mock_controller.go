package reku

import (
	"log"
)

type MockRecuperator struct {
	currentStatus Status
}

func NewMockRecuperator() *MockRecuperator {
	log.Println("Używam mockowego kontrolera rekuperatora.")
	return &MockRecuperator{
		currentStatus: Status{
			Recovery:              88,
			ExtractAirTemperature: "21.5 °C",
			SupplyAirTemperature:  "19.8 °C",
			OutdoorAirTemperature: "12.0 °C",
		},
	}
}

func (m *MockRecuperator) GetStatus() (Status, error) {
	log.Printf("[MOCK] Odczytano status: %+v", m.currentStatus)
	return m.currentStatus, nil
}

func (m *MockRecuperator) SetExtractAndSupplyFanSpeed(extractFanSpeed int, supplyFanSpeed int) error {
	log.Printf("[MOCK] Ustawiam moc na: %d%%", extractFanSpeed)
	return nil // Sukces
}
