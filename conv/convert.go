package conv

import "math"

// ToUnits converts the given price to uint64 units used by the trading engine
func ToUnits(number float64, precision uint8) uint64 {
	return uint64(number * math.Pow(10, float64(precision)))
}

// FromUnits converts the given price to uint64 units used by the trading engine
func FromUnits(number uint64, precision uint8) float64 {
	return float64(number) / math.Pow(10, float64(precision))
}
