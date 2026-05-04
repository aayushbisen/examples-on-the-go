package safedivider

import "errors"

var ErrNegative = errors.New("negative numbers not allowed")
var ErrZero = errors.New("cannot divide by zero")

func Divide(a, b float64) (result float64, err error) {
	if b == 0 {
		return 0, ErrZero
	}
	if a < 0 || b < 0 {
		return 0, ErrNegative
	}

	return (a / b), nil
}
