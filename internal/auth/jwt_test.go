package auth

import (
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	tokenSecret := "one two three"
	userID := uuid.New()
	jwt, err := MakeJWT(userID, tokenSecret, time.Second*10)
	if err != nil {
		t.Error(err)
	}

	jwtUserId, err := ValidateJWT(jwt, tokenSecret)
	if err != nil {
		t.Error(err)
	}

	if jwtUserId != userID {
		t.Error("jwt user id does not match")
	}
}

func TestDifferentSecretJWT(t *testing.T) {
	tokenSecret := "three four five"
	userID := uuid.New()
	jwt, err := MakeJWT(userID, tokenSecret, time.Second*10)
	if err != nil {
		t.Error(err)
	}

	_, err = ValidateJWT(jwt, "a differnt secret")
	wantErr := regexp.MustCompile("signature is invalid")
	if err == nil {
		t.Error("jwt user id should not not match")
	} else if !wantErr.Match([]byte(err.Error())) {
		t.Error("jwt user id should not not match")
	}
}

func TestExpiredJWT(t *testing.T) {
	tokenSecret := "five six seven"
	userID := uuid.New()
	duration := time.Millisecond * 10
	jwt, err := MakeJWT(userID, tokenSecret, duration)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(duration * 2)

	_, err = ValidateJWT(jwt, tokenSecret)
	wantErr := regexp.MustCompile("token is expired")
	if err == nil {
		t.Error("jwt should have expired")
	} else if !wantErr.Match([]byte(err.Error())) {
		t.Error("jwt should have expired")
	}
}
