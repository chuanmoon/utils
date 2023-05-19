package number

import "testing"

var (
	testCurrencys = []string{"USD", "EUR", "JPY"}

	testAmounts = []int64{-2, -22, -222, -222, -2222222, -20, -200, -2000, 2, 22, 222, 222, 2222222, 20, 200, 2000}
)

func TestFormatInt64Money(t *testing.T) {
	for _, currency := range testCurrencys {
		for _, amount := range testAmounts {
			result := FormatInt64Money(amount, currency)
			t.Log(result)
		}
	}
}
