package reku

// Status reprezentuje aktualny stan rekuperatora.
// Użyjemy tej struktury do przesyłania danych wewnątrz programu.
type Status struct {
	Mode                  string `json:"tryb_pracy"`
	Recovery              int    `json:"odzysk"`
	ExtractAirTemperature string `json:"temp_wyciagu"`
	SupplyAirTemperature  string `json:"temp_nawiewu"`
	OutdoorAirTemperature string `json:"temp_zewnatrz"`
	ActualSupplyFanSpeed  string `json:"wentylator_nawiewu"`
	ActualExtractFanSpeed string `json:"przeplyw_wyciagu"`
	HeatRecoveryPower     string `json:"odzysk_moc"`
}

// RecuperatorController to nasz interfejs.
// Definiuje on zestaw metod, które każda implementacja (prawdziwa, mockowa) musi posiadać.
type RecuperatorController interface {
	// GetStatus pobiera aktualny stan rekuperatora.
	GetStatus() (Status, error)

	SetExtractAndSupplyFanSpeed(extractFanSpeed int, supplyFanSpeed int) error

	// SetMode ustawia tryb pracy (np. chłodzenie/grzanie).
	// Dodajemy go dla przykładu, aby interfejs był bardziej kompletny.
	SetMode(mode string) error
}
