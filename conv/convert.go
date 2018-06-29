package conv

import (
	"math"
	"strconv"
	"strings"
)

// ToUnits converts the given price to uint64 units used by the trading engine
func ToUnits(amount string, precision uint8) uint64 {
	amounts := strings.Split(amount, ".")
	var dec uint64
	var remainder uint64

	if len(amounts) == 0 {
		return 0
	}
	if len(amounts) == 2 {
		decimals := amounts[1]
		if uint8(len(decimals)) > precision {
			decimals = decimals[:precision]
		}
		remainder, _ = strconv.ParseUint(decimals, 10, 64)
	}
	mult := uint64(math.Pow(10, float64(precision)))
	dec, _ = strconv.ParseUint(amounts[0], 10, 64)
	return dec*mult + remainder
}

// FromUnits converts the given price to uint64 units used by the trading engine
func FromUnits(number uint64, precision uint8) string {
	mult := uint64(math.Pow(10, float64(precision)))
	dec := number / mult
	remainder := number % mult
	rem := strconv.FormatUint(remainder, 10)
	i := precision - uint8(len(rem))
	for i > 0 {
		rem = "0" + rem
		i--
	}
	return strconv.FormatUint(dec, 10) + "." + rem
}

func pad(s, max string) string {
	if i := len(s); i < len(max) {
		s = max[i:] + s
	}
	return s
}
