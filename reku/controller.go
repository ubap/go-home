package reku

// Status structure is used for serializing data to send to the fronend.
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

type RecuperatorController interface {
	GetStatus() (Status, error)

	SetExtractAndSupplyFanSpeed(extractFanSpeed int, supplyFanSpeed int) error
}
