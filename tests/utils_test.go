package tests

import (
	"testing"

	"caesar-quiz/internal/services"
)

func TestCaesarEncryptDecrypt(t *testing.T) {
	cases := []struct {
		plain string
		shift int
		want  string
	}{
		{"HALO", 3, "KDOR"},
		{"abc", 1, "bcd"},
		{"xyz", 2, "zab"},
	}
	for _, c := range cases {
		enc := services.CaesarEncrypt(c.plain, c.shift)
		if enc != c.want {
			t.Fatalf("enc(%s,%d) = %s; want=%s", c.plain, c.shift, enc, c.want)
		}
		dec := services.CaesarDecrypt(enc, c.shift)
		if dec != c.plain {
			t.Fatalf("dec(%s,%d) = %s; want=%s", enc, c.shift, dec, c.plain)
		}
	}
}
