package converter

const freezingPoint = 32

func CelsiusToFahrenheit(cel float64) float64 {
	return (cel * (9.0 / 5.0)) + freezingPoint
}

func FahrenheitToCelsius(far float64) float64 {
	return (far - freezingPoint) * (5.0 / 9.0)
}
