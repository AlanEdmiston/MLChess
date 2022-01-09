package main

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//Vector is a 2D vector
type Vector struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (vect1 Vector) isOutOfBounds() bool {
	return vect1.X < 0 || vect1.X > 7 || vect1.Y < 0 || vect1.Y > 7
}

func (vect1 Vector) add(vect2 Vector) Vector {
	return Vector{vect1.X + vect2.X, vect1.Y + vect2.Y}
}

func (vect1 Vector) mult(c int) Vector {
	return Vector{vect1.X * c, vect1.Y * c}
}

func (vect1 Vector) toString() string {
	return fmt.Sprintf("(%d, %d)", vect1.X, vect1.Y)
}

func fromString(vect string) Vector {
	nums := strings.Split(vect[1:len(vect)-1], ",")
	x, _ := strconv.Atoi(nums[0])
	y, _ := strconv.Atoi(nums[1])
	return Vector{X: x, Y: y}
}

func (vect1 Vector) boardPosition() string {
	return fmt.Sprintf("%s%d", string(rune('a'+vect1.X)), vect1.Y+1)
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
