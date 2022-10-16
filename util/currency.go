package util

const (
	USD = "USD"
	GHS = "GHS"
	NGN = "NGN"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, GHS, NGN:
		return true
	}
	return false
}
