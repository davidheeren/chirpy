package auth

import (
	"testing"
)

func TestPassword(t *testing.T) {
	password := "hollow moon"
	hash, err := HashPassword(password)
	if err != nil {
		t.Error(err)
	}

	valid, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Error(err)
	}

	if !valid {
		t.Errorf("password and hash do not match")
	}
}

func TestNotPassword(t *testing.T) {
	hash, err := HashPassword("one two three")
	if err != nil {
		t.Error(err)
	}

	valid, err := CheckPasswordHash("four five six", hash)
	if err != nil {
		t.Error(err)
	}

	if valid {
		t.Errorf("password and hash should not match")
	}
}
