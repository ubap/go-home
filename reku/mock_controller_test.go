package reku

import "testing"

func Test_GetStatus(t *testing.T) {
	var recu RecuperatorController = NewMockRecuperator()

	status, _ := recu.GetStatus()
	if status.Power != 45 {
		t.Fail()
	}
}
