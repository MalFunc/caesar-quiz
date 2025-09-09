package services

import (
	"math/rand"
	"regexp"
	"time"
	"unicode"
)

// Caesar encrypt - only shifts letters; other chars preserved
func CaesarEncrypt(input string, shift int) string {
	shift = (shift%26 + 26) % 26
	runes := []rune(input)
	for i, r := range runes {
		if unicode.IsLetter(r) {
			if r >= 'A' && r <= 'Z' {
				runes[i] = rune((int(r-'A')+shift)%26) + 'A'
			} else if r >= 'a' && r <= 'z' {
				runes[i] = rune((int(r-'a')+shift)%26) + 'a'
			}
		}
	}
	return string(runes)
}

func CaesarDecrypt(input string, shift int) string {
	return CaesarEncrypt(input, (26-shift)%26)
}

// name regex: unicode letters and space allowed
var nameRegex = regexp.MustCompile(`^[\p{L} ]+$`)

func ValidateName(name string) bool {
	if len(name) == 0 || len(name) > 255 {
		return false
	}
	return nameRegex.MatchString(name)
}

// --- tambahan untuk game ---
func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomShift return angka shift acak 1-25
func RandomShift() int {
	return rand.Intn(25) + 1
}

// GenerateGameCode return kode unik game 6 huruf
func GenerateGameCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, 6)
	for i := range code {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}
