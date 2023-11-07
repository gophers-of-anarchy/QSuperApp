package helpers

import (
	"strconv"
	"strings"
)

func AmountToEnglish(text string) (float64, error) {
	//Persian digits
	text = strings.NewReplacer(
		"۰", "0",
		"۱", "1",
		"۲", "2",
		"۳", "3",
		"۴", "4",
		"۵", "5",
		"۶", "6",
		"۷", "7",
		"۸", "8",
		"۹", "9",
	).Replace(text)
	//Arabic digits
	text = strings.NewReplacer(
		"٤", "4",
		"٥", "5",
		"٦", "6",
	).Replace(text)
	return strconv.ParseFloat(text, 64)
}
