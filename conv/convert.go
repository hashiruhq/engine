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
	str := strconv.FormatUint(number, 10)
	n := uint8(len(str))
	if n <= precision {
		str = strings.Repeat("0", int(precision-n+1)) + str
		n = precision + 1
	}
	return str[:n-precision] + "." + str[n-precision:]
}
