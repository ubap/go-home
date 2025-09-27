package reku

import "testing"

func Test_Komfvent_login(t *testing.T) {
	recu := *NewKomfventRecuperator()

	loginError := recu.login()

	if loginError != nil {
		t.Errorf("login error, expected success")
	}
}

func Test_Komfvent_GetStatus(t *testing.T) {
	recu := *NewKomfventRecuperator()

	recu.GetStatus()
}
