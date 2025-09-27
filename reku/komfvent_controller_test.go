package reku

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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

func setupTestServer() *httptest.Server {
	// http.NewServeMux is a router that lets us define responses for different paths.
	mux := http.NewServeMux()

	// Handler for the login path.
	// For this example, we assume a successful login just returns a 200 OK status.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// You could add more logic here, like checking the username/password
		// from the request if your login function sends them.
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Login successful")
	})

	// Handler for the GetStatus path.
	// We return a sample JSON response that your GetStatus function would expect.
	mux.HandleFunc("/i.asp", func(w http.ResponseWriter, r *http.Request) {
		utf8Body := "<?xml version=\"1.0\" encoding=\"windows-1250\"?> <A><OMO>NORMALNY      </OMO><AI0>19.8 °C  </AI0><AI1>22.2 °C  </AI1><AI2>5.6 °C   </AI2><SP>35 </SP><SAF>35 %              </SAF><EAF>34 %              </EAF><SAFS>35 %              </SAFS><EAFS>34 %              </EAFS><FCG>71 %  </FCG><EC1>87 %  </EC1><EC2> 1079 W  </EC2><EC3>48 W     </EC3><EC4>0 W      </EC4><EC5A>0.29        </EC5A><EC5D>0.30        </EC5D><EC6D>1.05 kWh    </EC6D><EC6M>30.54 kWh   </EC6M><EC6T>1156.64 kWh </EC6T><EC7D>0.00 kWh    </EC7D><EC7M>0.00 kWh    </EC7M><EC7T>21.54 kWh   </EC7T><EC8D>19.86 kWh   </EC8D><EC8M>84.30 kWh   </EC8M><EC8T>11717.77 kWh</EC8T><ST>30.0 °C  </ST><ET>--.- °C  </ET><AQS>--.- %    </AQS><AQ>--.- %    </AQ><AHS>--.- %    </AHS><AH>--.- %    </AH><VF>203571212 </VF></A>"
		encoder := charmap.Windows1250.NewEncoder()
		win1250Bytes, _ := encoder.String(utf8Body)

		w.Header().Set("Content-Type", "text/xml")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, win1250Bytes)
	})

	// httptest.NewServer starts a server on a random available port.
	return httptest.NewServer(mux)
}
