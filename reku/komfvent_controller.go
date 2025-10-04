package reku

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type KomfventRecuperator struct {
	address  string
	username string
	password string
}

// KomfventData is the struct to deserialize status XML response into.
type KomfventData struct {
	// Root element is named A
	XMLName xml.Name `xml:"A"`

	Mode                  TrimmedString `xml:"OMO"`
	SupplyAirTemperature  TrimmedString `xml:"AI0"`
	ExtractAirTemperature TrimmedString `xml:"AI1"`
	OutdoorAirTemperature TrimmedString `xml:"AI2"`
	Unknown1              TrimmedString `xml:"SP"`
	ActualSupplyFanSpeed  TrimmedString `xml:"SAF"`
	ActualExtractFanSpeed TrimmedString `xml:"EAF"`
	TargetSupplyFanSpeed  TrimmedString `xml:"SAFS"`
	TargetExtractFanSpeed TrimmedString `xml:"EAFS"`
	FilterLifePercentage  TrimmedString `xml:"FCG"`
	Ec1                   TrimmedString `xml:"EC1"`
	HeatRecoveryPower     TrimmedString `xml:"EC2"`
	Ec3                   TrimmedString `xml:"EC3"`
	Ec4                   TrimmedString `xml:"EC4"`
	Ec5a                  TrimmedString `xml:"EC5A"`
	Ec5d                  TrimmedString `xml:"EC5D"`
	Ec6d                  TrimmedString `xml:"EC6D"`
	Ec6m                  TrimmedString `xml:"EC6M"`
	Ec6t                  TrimmedString `xml:"EC6T"`
	Ec7d                  TrimmedString `xml:"EC7D"`
	Ec7m                  TrimmedString `xml:"EC7M"`
	Ec7t                  TrimmedString `xml:"EC7T"`
	Ec8d                  TrimmedString `xml:"EC8D"`
	Ec8m                  TrimmedString `xml:"EC8M"`
	Ec8t                  TrimmedString `xml:"EC8T"`
	St                    TrimmedString `xml:"ST"`
	Et                    TrimmedString `xml:"ET"`
	Aqs                   TrimmedString `xml:"AQS"`
	Aq                    TrimmedString `xml:"AQ"`
	Ahs                   TrimmedString `xml:"AHS"`
	Ah                    TrimmedString `xml:"AH"`
	Vf                    TrimmedString `xml:"VF"`
}

func komfventDataToStatus(komfventData KomfventData) Status {
	status := Status{}
	status.Mode = string(komfventData.Mode)
	status.ExtractAirTemperature = string(komfventData.ExtractAirTemperature)
	status.OutdoorAirTemperature = string(komfventData.OutdoorAirTemperature)
	status.SupplyAirTemperature = string(komfventData.SupplyAirTemperature)
	status.ActualSupplyFanSpeed = string(komfventData.ActualSupplyFanSpeed)
	status.ActualExtractFanSpeed = string(komfventData.ActualExtractFanSpeed)
	status.HeatRecoveryPower = string(komfventData.HeatRecoveryPower)
	return status
}

func NewKomfventRecuperator() *KomfventRecuperator {
	return &KomfventRecuperator{address: "http://192.168.1.24", username: "user", password: "user"}
}

func (k *KomfventRecuperator) GetStatus() (Status, error) {
	start := time.Now()
	defer func() {
		log.Printf("GetStatus function execution took %s", time.Since(start))
	}()

	data, err := k.getStatusImpl()
	if errors.Is(err, ErrUnauthorized) {
		log.Printf("Unauthorized, trying to log in.")
		loginError := k.login()
		if loginError != nil {
			return Status{}, loginError
		}
		return k.getStatusImpl()
	}
	if err != nil {
		log.Printf("Error getting status response: %s", err)
		return Status{}, err
	}

	return data, nil
}

func (k *KomfventRecuperator) getStatusImpl() (Status, error) {
	resp, err := http.Get(k.address + "/i.asp")
	if err != nil {
		return Status{}, err
	}
	defer resp.Body.Close()

	data, err := processResponse(resp)
	if err != nil {
		return Status{}, err
	}

	return komfventDataToStatus(data), nil
}

func processResponse(resp *http.Response) (KomfventData, error) {
	var data KomfventData

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, fmt.Errorf("failed to read response body: %w", err)
	}

	win1250Decoder := charmap.Windows1250.NewDecoder()

	utf8Bytes, _, err := transform.Bytes(win1250Decoder, bodyBytes)
	if err != nil {
		return data, fmt.Errorf("failed to transform body to UTF-8: %w", err)
	}

	// Check for Unauthorized screen. Unfortunately the device returns 200 with some HTML in the body.
	if bytes.Contains(utf8Bytes, []byte("Niepoprawne")) {
		return data, ErrUnauthorized
	}

	bodyReader := bytes.NewReader(bodyBytes)
	decoder := xml.NewDecoder(bodyReader)
	decoder.CharsetReader = makeCharsetReader

	err = decoder.Decode(&data)
	if err != nil {
		// The XML parsing failed. It's helpful to include the body here too.
		return data, fmt.Errorf("failed to decode XML: %w. Body: %s", err, string(utf8Bytes))
	}

	return data, nil
}

func (k *KomfventRecuperator) SetExtractAndSupplyFanSpeed(extractFanSpeed int, supplyFanSpeed int) error {
	start := time.Now()
	defer func() {
		log.Printf("SetExtractAndSupplyFanSpeed function execution took %s", time.Since(start))
	}()

	payload := fmt.Sprintf("248=%d&256=%d&", extractFanSpeed, supplyFanSpeed)

	// Create a new request with the POST method
	req, err := http.NewRequest("POST", k.address+"/ajax.xml", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "text/plain;charset=UTF-8")

	// Create an HTTP client
	client := &http.Client{}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed
	log.Printf("Response Status: %s\n", resp.Status)

	// no error
	return nil
}

func (k *KomfventRecuperator) login() error {
	start := time.Now()
	defer func() {
		log.Printf("login function execution took %s", time.Since(start))
	}()

	formData := url.Values{}
	formData.Set("1", k.username)
	formData.Set("2", k.password)

	// http.PostForm automatically sets Content-Type to application/x-www-form-urlencoded
	resp, err := http.PostForm(k.address, formData)
	if err != nil {
		log.Printf("Error making request: %s", err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response: %s", err)
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Login Failed; body: " + string(body))
	}

	return nil
}

func makeCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	if charset == "windows-1250" {
		// Zwracamy specjalny reader, który w locie konwertuje
		// windows-1250 na UTF-8.
		return transform.NewReader(input, charmap.Windows1250.NewDecoder()), nil
	}
	// Zwracamy błąd, jeśli napotkamy inne, nieobsługiwane kodowanie.
	return nil, fmt.Errorf("nieznane kodowanie: %s", charset)
}

// TrimmedString to własny typ, który zachowuje się jak string,
// ale automatycznie usuwa białe znaki z początku i końca
// podczas dekodowania z XML.
type TrimmedString string

// UnmarshalXML implementuje interfejs xml.Unmarshaler dla naszego typu.
// Ta metoda zostanie automatycznie wywołana przez dekoder XML
// dla każdego pola, które jest typu TrimmedString.
func (ts *TrimmedString) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	// DecodeElement odczytuje zawartość tekstową elementu do zwykłego stringa.
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	// Kluczowy moment: usuwamy białe znaki i przypisujemy wynik
	// do naszej wartości TrimmedString.
	*ts = TrimmedString(strings.TrimSpace(s))

	return nil
}
