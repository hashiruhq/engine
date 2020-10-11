package utils

import (
	"github.com/ericlagergren/decimal"
	"github.com/rs/zerolog/log"
	"math"
)

// DecimalToZeroCtx godoc
// The context used for calculating calculus for market orders
var DecimalToZeroCtx = decimal.Context{
	Precision:     256,
	RoundingMode:  decimal.ToZero,
	OperatingMode: decimal.GDA,
	Traps:         ^(decimal.Inexact | decimal.Rounded | decimal.Subnormal),
	MaxScale:      6144,
	MinScale:      -6143,
}

// DecimalAwayFromZeroCtx godoc
// The context used for calculating calculus for market orders
var DecimalAwayFromZeroCtx = decimal.Context{
	Precision:     256,
	RoundingMode:  decimal.AwayFromZero,
	OperatingMode: decimal.GDA,
	Traps:         ^(decimal.Inexact | decimal.Rounded | decimal.Subnormal),
	MaxScale:      6144,
	MinScale:      -6143,
}

// Divide two uint64 numbers with a 10^prec precision and return the result in the same format
func Divide(x, y uint64, xprec, yprec, prec int) uint64 {
	xDec := decimal.WithContext(DecimalToZeroCtx).SetUint64(x)
	xDec.Quo(xDec, decimal.WithContext(DecimalToZeroCtx).SetUint64(y))
	xDec.Mul(xDec, decimal.New(10, -1*(xprec-yprec+prec-1))).Quantize(0)
	z, _ := xDec.Uint64()
	return z
}

// Multiply two uint64 numbers with a 10^prec precision and return the result in the same format
func Multiply(x, y uint64, xprec, yprec, prec int) uint64 {
	multiplier := decimal.New(10, -1*(xprec+yprec-prec-1))
	multiplier.Context = DecimalAwayFromZeroCtx
	xDec := decimal.WithContext(DecimalAwayFromZeroCtx).SetUint64(x)
	xDec.Mul(xDec, decimal.WithContext(DecimalAwayFromZeroCtx).SetUint64(y))
	xDec.Quo(xDec, multiplier).Quantize(0)
	if xDec.Cmp(decimal.WithContext(DecimalAwayFromZeroCtx).SetUint64(math.MaxUint64)) > 0 {
		log.Warn().
			Str("section", "math").
			Str("action", "multiply").
			Uint64("x", x).Uint64("y", y).
			Int("xprec", xprec).Int("yprec", yprec).Int("prec", prec).
			Str("result", xDec.String()).
			Msg("Unable to convert to uint64, number probably exceeds alowed bounds")
		return 0
	}
	z, _ := xDec.Uint64()
	return z
}

// Max returns the maximum value between two numbers
func Max(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}

// Min returns the minimum value between two numbers
func Min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
