package main

// source: https://www.cloudbees.com/blog/testing-http-handlers-go

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func doRequest[T any](t *testing.T, s *http.Server, method, path, token, body string, wantStatus int) T {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rr := httptest.NewRecorder()
	s.Handler.ServeHTTP(rr, req)

	if rr.Code != wantStatus {
		t.Fatalf("%s %s: got status %d, want %d (body: %s)",
			method, path, rr.Code, wantStatus, rr.Body.String())
	}

	var res T
	// return if empty struct type T
	if _, ok := any(res).(struct{}); ok {
		return res
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("%s %s: decode response %q: %v", method, path, rr.Body.String(), err)
	}
	return res
}

func TestHealthzEndpoint(t *testing.T) {
	cfg, err := createConfig()
	if err != nil {
		t.Fatal(err)
	}
	server := createServer(cfg)

	req := httptest.NewRequest("GET", "/api/healthz", nil)
	rr := httptest.NewRecorder()
	server.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("healthz handler returned wrong status code: got %v want %v",
			rr.Code, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Fatalf("healthz handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestRefreshEndpoint(t *testing.T) {
	// match the requirements of this lession:
	// https://www.boot.dev/lessons/f7285cef-5185-4b15-b5fc-9533ccaafe8a
	cfg, err := createConfig()
	if err != nil {
		t.Fatal(err)
	}
	server := createServer(cfg)

	// Reset database Test
	doRequest[struct{}](
		t,
		&server,
		"POST",
		"/admin/reset",
		"",
		"",
		http.StatusOK,
	)

	// Create saul user Test
	createSaulReq := `
{
	"email": "saul@bettercall.com",
	"password": "123456"
}`
	createSaulRes := doRequest[User](
		t,
		&server,
		"POST",
		"/api/users",
		"",
		createSaulReq,
		http.StatusCreated,
	)
	correctSaulEmail := "saul@bettercall.com"
	if createSaulRes.Email != correctSaulEmail {
		t.Fatalf("new user handler returned wrong email: got %v want %v",
			createSaulRes.Email, correctSaulEmail)
	}

	// Login saul Test
	loginSaulReq := `
{
	"email": "saul@bettercall.com",
	"password": "123456"
}`
	loginSaulRes := doRequest[loginUser](
		t,
		&server,
		"POST",
		"/api/login",
		"",
		loginSaulReq,
		http.StatusOK,
	)

	// Create chirp Test
	createChirpReq := `
{
	"body": "Let's just say I know a guy... who knows a guy... who knows another guy."
}`
	doRequest[Chirp](
		t,
		&server,
		"POST",
		"/api/chirps",
		loginSaulRes.Token,
		createChirpReq,
		http.StatusCreated,
	)

	// Refresh token Test
	refreshTokenRes := doRequest[refreshedToken](
		t,
		&server,
		"POST",
		"/api/refresh",
		loginSaulRes.RefreshToken,
		"",
		http.StatusOK,
	)

	// Create chirp 2 Test
	createChirpReq2 := `
{
	"body": "I'm the guy who's gonna win you this case."
}`
	doRequest[Chirp](
		t,
		&server,
		"POST",
		"/api/chirps",
		refreshTokenRes.Token,
		createChirpReq2,
		http.StatusCreated,
	)

	// Revoke token Test
	doRequest[struct{}](
		t,
		&server,
		"POST",
		"/api/revoke",
		loginSaulRes.RefreshToken,
		"",
		http.StatusNoContent,
	)

	// Refresh token 2 Test
	doRequest[refreshedToken](
		t,
		&server,
		"POST",
		"/api/refresh",
		loginSaulRes.RefreshToken,
		"",
		http.StatusUnauthorized,
	)
}
