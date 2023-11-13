package utils

import "testing"

func TestIsValidTRXAddress(t *testing.T) {
	result := IsValidTRXAddress("TTy7o4hXwuiztVe24EesKAB8haMcE5Keyo")
	if result {
		t.Log("the address is trx address")
		return
	}
	t.Error("the address is not trx address")
}
