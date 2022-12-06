package nmcli

import "testing"

func Test_GetStatus(t *testing.T) {
	msg, err := Status()
	if err != nil {
		t.Errorf("Failed to get status with message: \n%v\n", msg)
	}
}
