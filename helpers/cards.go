package helpers

import (
	"strconv"
	"strings"
)

func CardIsValid(number string) bool {
	var sum int
	for i := 0; i < len(number); i++ {
		digit, _ := strconv.Atoi(string(number[i]))
		var result int
		if (i+1)%2 == 0 {
			result = digit * 1
		} else {
			doubled := digit * 2
			if doubled > 9 {
				result = doubled - 9
			} else {
				result = doubled
			}
		}
		sum += result
	}
	return sum%10 == 0
}

func CardNumberToEnglish(text string) string {
	//Remove dashes
	text = strings.Replace(text, "-", "", -1)
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
	return strings.NewReplacer(
		"٤", "4",
		"٥", "5",
		"٦", "6",
	).Replace(text)
}
