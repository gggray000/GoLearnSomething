package util

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func WriteFloatValueToFile(value float64, fileName string) {
	valueText := fmt.Sprint(value)
	os.WriteFile(fileName, []byte(valueText), 0644)
}

func GetFloatValueFromFile(fileName string) (float64, error) {
	data, err := os.ReadFile(fileName)

	if err != nil {
		return 0, errors.New("Failed to find file ")
	}

	valueText := string(data)
	value, err := strconv.ParseFloat(valueText, 64)

	if err != nil {
		return 0, errors.New("Failed to parse stored value ")
	}
	return value, nil
}
