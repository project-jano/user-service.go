package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
)

const (
	testUser     = "test-user"
	testPassword = "test-pwd"
)

func TestMissingBasicAuth(t *testing.T) {
	request, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	if IsAuthenticated(request, testUser, testPassword, "/") {
		t.Fail()
	}
}

func TestIncorrectAuthMethod(t *testing.T) {
	request, err := http.NewRequest("GET", "http:/test/", nil)
	encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", testUser, testPassword)))
	request.Header["Authorization"] = []string{"Token " + encoded}

	if err != nil {
		t.Fatal(err)
	}

	if IsAuthenticated(request, testUser, testPassword, "/") {
		t.Fail()
	}
}

func TestBasicAuthOk(t *testing.T) {
	request, err := http.NewRequest("GET", "/test", nil)
	encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", testUser, testPassword)))
	request.Header["Authorization"] = []string{"Basic " + encoded}

	if err != nil {
		t.Fatal(err)
	}

	if !IsAuthenticated(request, testUser, testPassword, "/") {
		t.Log("BasicAuth test failed")
		t.Fail()
	}
}
