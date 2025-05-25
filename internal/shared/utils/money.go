package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func FloatToBRL(value float64) string {
	// Format with 2 decimal places
	formattedValue := fmt.Sprintf("%.2f", value)

	// Split the string into integer and decimal parts
	parts := strings.Split(formattedValue, ".")
	integer := parts[0]
	decimal := parts[1]

	// Add thousands separators to the integer part
	for i := len(integer) - 3; i > 0; i -= 3 {
		integer = integer[:i] + "." + integer[i:]
	}

	// Put it all back together
	return fmt.Sprintf("R$ %s,%s", integer, decimal)
}

func BRLToFloat(brl string) (float64, error) {
	brl = strings.Replace(brl, "R$", "", 1)
	brl = strings.Replace(brl, ".", "", -1)
	brl = strings.Replace(brl, ",", ".", 1)
	value, err := strconv.ParseFloat(strings.TrimSpace(brl), 64)
	if err != nil {
		return 0, fmt.Errorf("invalid BRL format: %v", err)
	}
	return value, nil
}
