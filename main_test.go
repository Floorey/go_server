package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func clearDB() {
	db.Exec("DELETE FROM vms")
}

func TestGetVMs(t *testing.T) {
	InitDB()
	clearDB() // Clear the database before running the test

	req, err := http.NewRequest("GET", "/api/vms", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Authenticate(http.HandlerFunc(GetVMs))

	// Create a test token
	token, err := createTestToken()
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", "token="+token)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var vms []VM
	err = json.NewDecoder(rr.Body).Decode(&vms)
	if err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if len(vms) != 0 {
		t.Errorf("handler returned unexpected body: got %v want an empty array", rr.Body.String())
	}
}

func TestCreateVM(t *testing.T) {
	InitDB()
	clearDB() // Clear the database before running the test

	vm := VM{Name: "test-vm", Image: "test-image", State: "running"}
	payload, err := json.Marshal(vm)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/api/vms", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := Authenticate(http.HandlerFunc(CreateVM))

	// Create a test token
	token, err := createTestToken()
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Cookie", "token="+token)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var createdVM VM
	err = json.NewDecoder(rr.Body).Decode(&createdVM)
	if err != nil {
		t.Fatal(err)
	}

	if createdVM.Name != vm.Name || createdVM.Image != vm.Image || createdVM.State != vm.State {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), payload)
	}
}

func TestLogin(t *testing.T) {
	creds := Credentials{Username: "user", Password: "password"}
	payload, err := json.Marshal(creds)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Login)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	cookies := rr.Result().Cookies()
	if len(cookies) == 0 {
		t.Errorf("handler did not return a token")
	} else {
		if token := cookies[0].Value; token == "" {
			t.Errorf("handler did not return a token")
		}
	}
}

// Helper function to create a test token
func createTestToken() (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}
