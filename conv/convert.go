package conv

// ToUnits converts the given price to uint64 units used by the trading engine
func ToUnits(amounts string, precision uint8) uint64 {
	bytes := []byte(amounts)
	size := len(bytes)
	start := false
	pointPos := 0
	var dec uint64
	i := 0
	for i = 0; i < size && (!start || (start && i-pointPos <= int(precision))); i++ {
		if !start && bytes[i] == '.' {
			start = true
			pointPos = i
		} else {
			dec = 10*dec + uint64(bytes[i]-48) // ascii char for 0
		}
	}
	for i-pointPos <= int(precision) {
		dec *= 10
		i++
	}
	return dec
}

// FromUnits converts the given price to uint64 units used by the trading engine
func FromUnits(number uint64, precision uint8) string {
	bytes := []byte{48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48, 48}
	i := 0
	for (number != 0 || i < int(precision)) && i <= 28 {
		add := uint8(number % 10)
		number /= 10
		bytes[28-i] = 48 + add
		if i == int(precision)-1 {
			i++
			bytes[28-i] = 46 // . char
		}
		i++
	}
	i--
	if bytes[28-i] == 46 {
		return string(bytes[28-i-1:])
	}

	return string(bytes[28-i:])
}
