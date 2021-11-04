package main

import (
	"errors"
	"fmt"
	"sort"
)

//Vector is a 2D vector
type Vector struct {
	x int
	y int
}

func (vect1 Vector) isOutOfBounds() bool {
	return vect1.x < 0 || vect1.x > 7 || vect1.y < 0 || vect1.y > 7
}

func (vect1 Vector) add(vect2 Vector) Vector {
	return Vector{vect1.x + vect2.x, vect1.y + vect2.y}
}

func (vect1 Vector) mult(c int) Vector {
	return Vector{vect1.x * c, vect1.y * c}
}

func (vect1 Vector) toString() string {
	return fmt.Sprintf("(%d, %d)", vect1.x, vect1.y)
}

func geometricAbs(value float64) float64 {
	if value >= 1 {
		return value
	}
	if value <= 0 {
		errors.New("cannot get geometric absolute value of negative number")
	}
	return 1 / value
}

func sortDiminishing(values []float64) {
	sort.Slice(values, func(i, j int) bool { return geometricAbs(values[i]) > geometricAbs(values[j]) })
}
