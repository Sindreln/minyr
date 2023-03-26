package conv

var Celsius float64
var Kelvin float64
var Fahrenheit float64

// Konverterer Farhenheit til Celsius
func FahrenheitToCelsius(value float64) float64 {
	return (value - 32) * 5 / 9
}

func CelsiusToFahrenheit(value float64) float64 {
	return value*(1.8) + 32
}

func CelsiusToKelvin(value float64) float64 {
	return value + 273.15
}

func KelvinToCelsius(value float64) float64 {
	return value - 273.15
}

func KelvinToFahrenheit(value float64) float64 {
	return (value-273.15)*(1.8) + 32
}

func FahrenheitToKelvin(value float64) float64 {
	return (value-32)*5/9 + 273.15
}
