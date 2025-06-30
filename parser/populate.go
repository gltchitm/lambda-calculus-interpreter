package parser

import (
	"errors"
	"fmt"
	"strconv"
)

func parsePopulateBound(name string, value string) (int, error) {
	const (
		min = 0
		max = 5_000
	)

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.New("malformed " + name + " value")
	} else if intValue < min {
		return 0, errors.New(name + " cannot be less than " + fmt.Sprint(min))
	} else if intValue > max {
		return 0, errors.New(name + " cannot exceed " + fmt.Sprint(max))
	}

	return intValue, nil
}
