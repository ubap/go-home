package reku

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type KomfventRecuperator struct {
	address  string
	username string
	password string
}

// KomfventData odzwierciedla strukturę odpowiedzi XML.
// Nazwa struktury może być dowolna, ale dobrze, by była opisowa.
type KomfventData struct {
	// To specjalne pole mówi parserowi, że oczekujemy,
	// że główny element XML będzie miał nazwę "A".
	XMLName xml.Name `xml:"A"`

	// Tag `xml:"OMO"` mapuje element <OMO> na pole Mode w strukturze.
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
	Ec2                   TrimmedString `xml:"EC2"`
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

func NewKomfventRecuperator() *KomfventRecuperator {
	// 1. Odczytaj zmienną środowiskową o nazwie RECUPERATOR_ADDRESS
	//    Nazwa zmiennej jest dowolna, ale dobrze, by była opisowa.
	address := os.Getenv("RECUPERATOR_ADDRESS")

	// 2. Sprawdź, czy zmienna została ustawiona.
	//    Jeśli os.Getenv() nie znajdzie zmiennej, zwróci pusty string "".
	if address == "" {
		// 3. Jeśli zmienna jest pusta, użyj wartości domyślnej.
		address = "http://192.168.1.24"
		fmt.Println("Zmienna środowiskowa RECUPERATOR_ADDRESS nie jest ustawiona. Używam domyślnego adresu:", address)
	} else {
		fmt.Println("Znaleziono zmienną środowiskową RECUPERATOR_ADDRESS. Używam adresu:", address)
	}

	// Zwróć strukturę z adresem (ze zmiennej środowiskowej lub domyślnym)
	return &KomfventRecuperator{address: address, username: "user", password: "user"}
}

func (k *KomfventRecuperator) GetStatus() (Status, error) {
	loginError := k.login()
	if loginError != nil {
		return Status{}, loginError
	}

	resp, err := http.Get(k.address + "/i.asp")
	if err != nil {
		fmt.Println("Error making request:", err)
		return Status{}, err
	}
	defer resp.Body.Close()

	var data KomfventData

	// Tworzymy nowy dekoder, który czyta bezpośrednio z ciała odpowiedzi.
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = makeCharsetReader
	// Dekodujemy XML do naszej struktury.
	// Musimy przekazać wskaźnik do naszej zmiennej `data`.
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding XML:", err)
		return Status{}, err
	}

	return komfventDataToStatus(data), nil
}

func (k *KomfventRecuperator) SetPower(power int) error {
	//TODO implement me
	panic("implement me")
}

func (k *KomfventRecuperator) SetMode(mode string) error {
	//TODO implement me
	panic("implement me")
}

func (k *KomfventRecuperator) login() error {
	// Create the form data
	formData := url.Values{}
	formData.Set("1", k.username)
	formData.Set("2", k.password)

	// Make the POST request
	// http.PostForm automatically sets Content-Type to application/x-www-form-urlencoded
	resp, err := http.PostForm(k.address, formData)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
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

func komfventDataToStatus(komfventData KomfventData) Status {
	status := Status{}
	status.ExtractAirTemperature = string(komfventData.ExtractAirTemperature)
	status.OutdoorAirTemperature = string(komfventData.OutdoorAirTemperature)
	status.SupplyAirTemperature = string(komfventData.SupplyAirTemperature)
	return status
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
