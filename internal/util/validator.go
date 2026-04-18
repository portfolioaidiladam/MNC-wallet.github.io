package util

import "errors"

// ValidatePhone memeriksa format nomor HP Indonesia:
//   - prefix "08"
//   - panjang 10..13 digit
//   - semua karakter angka
//
// Tidak pakai regex supaya mudah dibaca reviewer dan menghindari dependency.
func ValidatePhone(s string) error {
	if len(s) < 10 || len(s) > 13 {
		return errors.New("phone number must be 10 to 13 digits")
	}
	if len(s) < 2 || s[0] != '0' || s[1] != '8' {
		return errors.New("phone number must start with 08")
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return errors.New("phone number must contain digits only")
		}
	}
	return nil
}

// ValidatePIN memeriksa PIN 6 digit numerik.
func ValidatePIN(s string) error {
	if len(s) != 6 {
		return errors.New("pin must be exactly 6 digits")
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return errors.New("pin must contain digits only")
		}
	}
	return nil
}