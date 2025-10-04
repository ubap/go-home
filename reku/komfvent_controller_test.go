package reku

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"golang.org/x/text/encoding/charmap"
)

func Test_Komfvent_login(t *testing.T) {
	server := setupTestServer()
	defer server.Close()
	recu := *NewKomfventRecuperator()
	recu.address = server.URL

	loginError := recu.login()

	if loginError != nil {
		t.Errorf("login error, expected success")
	}
}

func Test_Komfvent_GetStatus(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	recu := *NewKomfventRecuperator()
	recu.address = server.URL

	status, _ := recu.GetStatus()
	if status.OutdoorAirTemperature != "5.6 °C" {
		t.Errorf("Incorrect temp")
	}
}

func TestKomfventRecuperator_SetExtractAndSupplyFanSpeed(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	recu := *NewKomfventRecuperator()
	recu.address = server.URL

	err := recu.SetExtractAndSupplyFanSpeed(20, 30)
	if err != nil {
		t.Errorf("SetExtractAndSupplyFanSpeed error, expected success")
	}

	status, _ := recu.GetStatus()
	if status.ActualExtractFanSpeed != "20 %" && status.ActualSupplyFanSpeed != "30 %" {
		t.Errorf("Incorrect fan speed")
	}
}

func setupTestServer() *httptest.Server {
	mux := http.NewServeMux()
	var extractFanSpeed int
	var supplyFanSpeed int

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Login successful")
	})

	mux.HandleFunc("/ajax.xml", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "text/plain;charset=UTF-8" {
			errMsg := fmt.Sprintf("Unsupported Media Type")
			http.Error(w, errMsg, http.StatusUnsupportedMediaType)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "mock server could not read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		values, err := url.ParseQuery(string(bodyBytes))
		if err != nil {
			http.Error(w, "mock server could not parse body", http.StatusBadRequest)
			return
		}

		extractStr := values.Get("248")
		supplyStr := values.Get("256")

		extractFanSpeed, _ = strconv.Atoi(extractStr)
		supplyFanSpeed, _ = strconv.Atoi(supplyStr)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Payload received by mock server")
	})

	mux.HandleFunc("/i.asp", func(w http.ResponseWriter, r *http.Request) {
		utf8Body := "<?xml version=\"1.0\" encoding=\"windows-1250\"?> <A><OMO>NORMALNY      </OMO><AI0>19.8 °C  </AI0><AI1>22.2 °C  </AI1><AI2>5.6 °C   </AI2><SP>35 </SP><SAF>" + strconv.Itoa(supplyFanSpeed) + " %              </SAF><EAF>" + strconv.Itoa(extractFanSpeed) + " %              </EAF><SAFS>35 %              </SAFS><EAFS>34 %              </EAFS><FCG>71 %  </FCG><EC1>87 %  </EC1><EC2> 1079 W  </EC2><EC3>48 W     </EC3><EC4>0 W      </EC4><EC5A>0.29        </EC5A><EC5D>0.30        </EC5D><EC6D>1.05 kWh    </EC6D><EC6M>30.54 kWh   </EC6M><EC6T>1156.64 kWh </EC6T><EC7D>0.00 kWh    </EC7D><EC7M>0.00 kWh    </EC7M><EC7T>21.54 kWh   </EC7T><EC8D>19.86 kWh   </EC8D><EC8M>84.30 kWh   </EC8M><EC8T>11717.77 kWh</EC8T><ST>30.0 °C  </ST><ET>--.- °C  </ET><AQS>--.- %    </AQS><AQ>--.- %    </AQ><AHS>--.- %    </AHS><AH>--.- %    </AH><VF>203571212 </VF></A>"
		encoder := charmap.Windows1250.NewEncoder()
		win1250Bytes, _ := encoder.String(utf8Body)

		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, win1250Bytes)
	})

	return httptest.NewServer(mux)
}
