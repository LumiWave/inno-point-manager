package util

import (
	"fmt"
	"math/big"
)

// 4000000 => 0.00004 float64
func ToDecimalEncf(value string, decimal int64) float64 {
	fValue, _ := toDecimalEnc(value, decimal).Float64()
	return fValue
}

// 4000000 => "0.00004" string
func ToDecimalEncStr(value string, decimal int64) string {
	return toDecimalEnc(value, decimal).Text('f', -1)
}

// 0.00004 => 4000000 float64
func ToDecimalDecf(value float64, decimal int64) float64 {
	fValue, _ := toDecimalDec(value, decimal).Float64()
	return fValue
}

// 0.00004 => "4000000" string
func ToDecimalDecStr(value float64, decimal int64) string {
	return toDecimalDec(value, decimal).Text('f', -1)
}

func toDecimalEnc(value string, decimal int64) *big.Float {
	scale := new(big.Float).SetFloat64(1)
	scale.SetString("1e" + fmt.Sprintf("%d", decimal))
	newValue, _ := new(big.Float).SetString(value)
	return new(big.Float).Quo(newValue, scale)
}

func toDecimalDec(value float64, decimal int64) *big.Float {
	scale := new(big.Float).SetFloat64(1)
	scale.SetString("1e" + fmt.Sprintf("%d", decimal))
	newValue := new(big.Float).SetFloat64(value)
	return new(big.Float).Mul(newValue, scale)
}
