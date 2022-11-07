package util

const (
	USD = "USD"
	THB = "THB"
)

func IsValidCurrency(currency string) bool {
	switch currency {
	case USD, THB:
		return true
	}
	return false
}
