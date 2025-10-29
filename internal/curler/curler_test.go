package curler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rancher-sandbox/scc-product-version-verifier/internal/curler"
)

func TestCurlVerify(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the query parameters
		expectedQuery := "/?identifier=product&version=1.0&arch=x86_64"
		if r.URL.String() != expectedQuery {
			t.Errorf("unexpected query: got %v want %v", r.URL.String(), expectedQuery)
		}

		// Check headers
		if r.Header.Get("Authorization") != "Token token=test-reg-code" {
			t.Errorf("unexpected Authorization header: got %v", r.Header.Get("Authorization"))
		}

		// Send a response
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[{"id": 1, "name": "product"}]`)
	}))
	defer server.Close()

	// Override the ProductQueryUrl to use the mock server
	originalURL := curler.ProductQueryUrl
	curler.ProductQueryUrl = server.URL
	defer func() { curler.ProductQueryUrl = originalURL }()

	// Call the function to be tested
	err := curler.CurlVerify("product", "1.0", "x86_64", "test-reg-code")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCurlVerify_NotFound(t *testing.T) {
	// Create a mock server that returns an empty array
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[]`)
	}))
	defer server.Close()

	// Override the ProductQueryUrl to use the mock server
	originalURL := curler.ProductQueryUrl
	curler.ProductQueryUrl = server.URL
	defer func() { curler.ProductQueryUrl = originalURL }()

	// Call the function to be tested
	err := curler.CurlVerify("product", "1.0", "x86_64", "test-reg-code")
	if err == nil {
		t.Errorf("expected an error, but got nil")
	}

	expectedError := "product not found"
	if err.Error() != expectedError {
		t.Errorf("unexpected error message: got %q want %q", err.Error(), expectedError)
	}
}
