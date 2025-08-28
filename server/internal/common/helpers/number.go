package helpers

import (
	"strconv"
	"strings"
)

func FormatDecimal(f float64, precision int) string {
	s := strconv.FormatFloat(f, 'f', precision, 64)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}
