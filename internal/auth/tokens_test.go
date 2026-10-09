package auth

import (
	"net/http"
	"testing"
)

func TestBearerToken(t *testing.T) {
	headers := http.Header{
		"Authorization": []string{"Bearer sometoken"},
		"AnotherHeader": []string{"Random", "Words"},
	}

	bToken, err := GetBearerToken(headers)
	if err != nil {
		t.Error(err)
	}

	if (bToken != "sometoken") {
		t.Errorf("bearer tokens do not match")
	}
}

func TestNotBearerToken(t *testing.T) {
	headers := http.Header{
		"Authorization": []string{"NotBearer sometoken"},
	}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Error(err)
	}
}
