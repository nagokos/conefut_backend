package validation

import (
	"errors"
	"unicode/utf8"
)

func CheckStringLen(maxLen int) func(s string) error {
	return func(s string) error {
		if utf8.RuneCountInString(s) > maxLen {
			return errors.New("value is greater than the required length")
		}
		return nil
	}
}
