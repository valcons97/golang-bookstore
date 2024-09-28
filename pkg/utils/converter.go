package converter

import "math/big"

func ConvertToBigInt(value *float64) *big.Int {

	if value == nil {
		return big.NewInt(0)
	}

	cents := *value * 100

	bigInt := big.NewInt(int64(cents))
	return bigInt
}

func ConvertToFloat(value *big.Int) *float64 {
	if value == nil {
		return nil
	}

	floatVal := new(big.Float).SetInt(value)
	floatVal.Quo(floatVal, big.NewFloat(100))

	result, _ := floatVal.Float64()

	return &result
}
