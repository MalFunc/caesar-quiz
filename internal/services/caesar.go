package services

import (
	cryptorand "crypto/rand"
	"math/big"
	"math/rand"
	"regexp"
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

// ValidShift memastikan shift berada di rentang Caesar yang masuk akal (1-25).
func ValidShift(shift int) bool {
	return shift >= 1 && shift <= 25
}

// RandomShift return angka shift acak 1-25
func RandomShift() int {
	return rand.Intn(25) + 1
}

// GenerateGameCode return kode unik game 6 huruf (crypto/rand)
func GenerateGameCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	code := make([]byte, 6)
	for i := range code {
		n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			// fallback tak terhindarkan: tetap gunakan sumber acak lain
			code[i] = letters[rand.Intn(len(letters))]
			continue
		}
		code[i] = letters[n.Int64()]
	}
	return string(code)
}
