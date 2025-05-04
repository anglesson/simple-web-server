package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func FloatToBRL(value float64) string {
	brl := fmt.Sprintf("R$ %.2f", value)
	brl = strings.Replace(brl, ".", ",", 1)
	if idx := strings.LastIndex(brl, ","); idx > 3 {
		brl = brl[:idx-3] + "." + brl[idx-3:]
	}
	return brl
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
