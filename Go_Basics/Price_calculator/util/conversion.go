package util

import (
	"errors"
	"strconv"
)

func StringsToFloat(strings []string) ([]float64, error) {
	var floats []float64
	for _, stringValue := range strings {
		floatValue, err := strconv.ParseFloat(stringValue, 64)
		if err != nil {
			return nil, errors.New("Fail to convert string.")
		}
		floats = append(floats, floatValue)
	}
	return floats, nil
}
