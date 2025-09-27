package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Definiujemy handler (funkcję obsługującą zapytania)
	http.HandleFunc("/i.asp", func(w http.ResponseWriter, r *http.Request) {
		// Jest to wygodniejsze dla czytelności XMLa
		xmlContent := "<?xml version=\"1.0\" encoding=\"windows-1250\"?> <A><OMO>NORMALNY      </OMO><AI0>19.8 °C  </AI0><AI1>22.2 °C  </AI1><AI2>7.6 °C   </AI2><SP>35 </SP><SAF>35 %              </SAF><EAF>34 %              </EAF><SAFS>35 %              </SAFS><EAFS>34 %              </EAFS><FCG>71 %  </FCG><EC1>87 %  </EC1><EC2> 1079 W  </EC2><EC3>48 W     </EC3><EC4>0 W      </EC4><EC5A>0.29        </EC5A><EC5D>0.30        </EC5D><EC6D>1.05 kWh    </EC6D><EC6M>30.54 kWh   </EC6M><EC6T>1156.64 kWh </EC6T><EC7D>0.00 kWh    </EC7D><EC7M>0.00 kWh    </EC7M><EC7T>21.54 kWh   </EC7T><EC8D>19.86 kWh   </EC8D><EC8M>84.30 kWh   </EC8M><EC8T>11717.77 kWh</EC8T><ST>30.0 °C  </ST><ET>--.- °C  </ET><AQS>--.- %    </AQS><AQ>--.- %    </AQ><AHS>--.- %    </AHS><AH>--.- %    </AH><VF>203571212 </VF></A>"

		// Ustawiamy nagłówek Content-Type, aby klient wiedział, że otrzymuje XML
		w.Header().Set("Content-Type", "text/xml")

		// Ustawiamy kod odpowiedzi na 200 OK
		w.WriteHeader(http.StatusOK)

		// Wpisujemy zawartość XML do ciała odpowiedzi
		fmt.Fprint(w, xmlContent)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		htmlContent := "<html></html>"

		w.Header().Set("Content-Type", "text/html")

		// Ustawiamy kod odpowiedzi na 200 OK
		w.WriteHeader(http.StatusOK)

		// Wpisujemy zawartość XML do ciała odpowiedzi
		fmt.Fprint(w, htmlContent)
	})

	// Informacja o starcie serwera
	fmt.Println("Serwer HTTP działa na porcie 8000...")

	// Uruchamiamy serwer. Funkcja nasłuchuje na porcie 8000.
	// log.Fatal opakowuje wywołanie, aby złapać i wyświetlić ewentualne błędy przy starcie.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
