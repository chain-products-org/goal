package mathx

import (
	"math/big"

	"github.com/gophero/goal/constraintx"
)

func Add(addend1 float64, addend2 float64) float64 {
	result, _ := NewFromFloat(addend1).Add(NewFromFloat(addend2)).Float64()
	return result
}

func BigAdd(a, b *big.Int) *big.Int {
	return new(big.Int).Add(a, b)
}

func BigSub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

func BigMul(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(a, b)
}

func BigDiv(a, b *big.Int) *big.Int {
	return new(big.Int).Div(a, b)
}

func BigMod(a, b *big.Int) *big.Int {
	return new(big.Int).Mod(a, b)
}

func BigPow(a *big.Int, b int64) *big.Int {
	return new(big.Int).Exp(a, big.NewInt(b), nil)
}

func BigFloatAdd(a, b *big.Float) *big.Float {
	return new(big.Float).Add(a, b)
}

func BigFloatSub(a, b *big.Float) *big.Float {
	return new(big.Float).Sub(a, b)
}

func BigFloatMul(a, b *big.Float) *big.Float {
	return new(big.Float).Mul(a, b)
}

func BigFloatDiv(a, b *big.Float) *big.Float {
	return new(big.Float).Quo(a, b)
}

func BigFloatStrAdd(a, b string) *big.Float {
	af, _ := new(big.Float).SetString(a)
	bf, _ := new(big.Float).SetString(b)
	return new(big.Float).Add(af, bf)
}

func BigFloatStrSub(a, b string) *big.Float {
	af, _ := new(big.Float).SetString(a)
	bf, _ := new(big.Float).SetString(b)
	return new(big.Float).Sub(af, bf)
}

func BigFloatStrMul(a, b string) *big.Float {
	af, _ := new(big.Float).SetString(a)
	bf, _ := new(big.Float).SetString(b)
	return new(big.Float).Mul(af, bf)
}

func BigFloatStrDiv(a, b string) *big.Float {
	af, _ := new(big.Float).SetString(a)
	bf, _ := new(big.Float).SetString(b)
	return new(big.Float).Quo(af, bf)
}

func Sub(minuend float64, subtrahend float64) float64 {
	result, _ := NewFromFloat(minuend).Sub(NewFromFloat(subtrahend)).Float64()
	return result
}

func Mul(multiplicand float64, multiplier float64) float64 {
	result, _ := NewFromFloat(multiplicand).Mul(NewFromFloat(multiplier)).Float64()
	return result
}

func Div(dividend float64, divisor float64, precision int32) float64 {
	result, _ := NewFromFloat(dividend).DivRound(NewFromFloat(divisor), precision).Float64()
	return result
}

func DivFloatUp(dividend float64, divisor float64, precision int32) float64 {
	result, _ := NewFromFloat(dividend).Div(NewFromFloat(divisor)).RoundUp(precision).Float64()
	return result
}

func DivTrunc(dividend float64, divisor float64, precision int32) float64 {
	result, _ := NewFromFloat(dividend).Div(NewFromFloat(divisor)).Truncate(precision).Float64()
	return result
}

func DivInt64(dividend int64, divisor int64, precision int32) float64 {
	result, _ := NewFromInt(dividend).DivRound(NewFromInt(divisor), precision).Float64()
	return result
}

func DivTruncInt64(dividend int64, divisor int64, divisorTwo int64, precision int32) float64 {
	result, _ := NewFromInt(dividend).Div(NewFromInt(divisor)).Div(NewFromInt(divisorTwo)).Truncate(precision).Float64()
	return result
}

func MulInt64AndFloat64(multiplicand int64, multiplier float64) float64 {
	result, _ := NewFromInt(multiplicand).Mul(NewFromFloat(multiplier)).Float64()
	return result
}

func MulBigFloat(multiplicand *big.Float, multiplier *big.Float) *big.Float {
	m1, _ := multiplicand.Float64()
	m2, _ := multiplier.Float64()
	result := Mul(m1, m2)
	return big.NewFloat(result)
}

// func CompareToBigFloat(f1 *big.Float, f2 *big.Float) bool {
//	return f1.
// }

func Maxn[T constraintx.Number](ns ...T) T {
	max := ns[0]
	var a, b T
	for i, j := 0, len(ns)-1; i < len(ns)/2 && j > i; {
		a, b = ns[i], ns[j]
		if a > b {
			if a > max {
				max = a
			}
		} else {
			if b > max {
				max = b
			}
		}
		i++
		j--
	}
	return max
}

func Minn[T constraintx.Number](ns ...T) T {
	min := ns[0]
	var a, b T
	for i, j := 0, len(ns)-1; i < len(ns)/2 && j > i; {
		a, b = ns[i], ns[j]
		if a < b {
			if a < min {
				min = a
			}
		} else {
			if b < min {
				min = b
			}
		}
		i++
		j--
	}
	return min
}

func Sumn[T constraintx.Number](ns ...T) T {
	var r T
	for _, v := range ns {
		r += v
	}
	return r
}
