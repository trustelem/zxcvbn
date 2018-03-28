package fuzz

import (
	"github.com/trustelem/zxcvbn"
)

func Fuzz(data []byte) int {
	password := string(data)

	_ = zxcvbn.PasswordStrength(password, nil)
	return 1
}
