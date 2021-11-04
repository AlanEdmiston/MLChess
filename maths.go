package main

func PI(values []float64) float64 {
	val := 1.0
	for _, value := range values {
		val *= value
	}
	return val
}
