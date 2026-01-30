package curler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rancherlabs/scc-product-version-verifier/internal/curler"
	"github.com/stretchr/testify/assert"
)

func TestCurlVerify(t *testing.T) {
	asserts := assert.New(t)

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the query parameters
		expectedQuery := "identifier=product&version=1.0&arch=x86_64"
		asserts.Equal(expectedQuery, r.URL.RawQuery)

		// Check headers
		asserts.Equal("Token token=test-reg-code", r.Header.Get("Authorization"))

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
	productList, err := curler.CurlVerify("product", "1.0", "x86_64", "test-reg-code")
	asserts.Nil(err)
	asserts.NotNil(productList)
}

func TestCurlVerify_NotFound(t *testing.T) {
	asserts := assert.New(t)
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
	_, err := curler.CurlVerify("product", "1.0", "x86_64", "test-reg-code")
	asserts.NotNil(err)
	asserts.EqualError(err, "product not found")
}
