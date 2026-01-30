package suseconnect

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SUSE/connect-ng/pkg/connection"
	"github.com/stretchr/testify/assert"
)

func TestLocalServer(t *testing.T) {
	asserts := assert.New(t)
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the query parameters
		asserts.Equal(fmt.Sprintf("%s?identifier=product&version=1.0&arch=x86_64", ProductsQueryPath), r.URL.String())

		// Check headers
		asserts.Equal("Token token=test-reg-code", r.Header.Get("Authorization"))

		// Send a response
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[{"id": 1, "name": "product"}]`)
	}))
	defer server.Close()

	// Override base option URL
	options := Options()
	options.URL = server.URL
	conn := Connection(&options)

	path := preparePathAndQuery("product", "1.0", "x86_64")
	request, buildErr := conn.BuildRequestRaw("GET", path, nil)
	asserts.Nil(buildErr)

	connection.AddRegcodeAuth(request, "test-reg-code")
	response, doErr := conn.Do(request)
	asserts.Nil(doErr)
	asserts.NotNil(response)
}

func TestVerify(t *testing.T) {
	asserts := assert.New(t)

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check the query parameters
		asserts.Equal(fmt.Sprintf("%s?identifier=rancher&version=2.13.2&arch=unknown", ProductsQueryPath), r.URL.String())

		// Check headers
		asserts.Equal("Token token=test-reg-code", r.Header.Get("Authorization"))

		// Send a response with a product
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `[{"identifier": "rancher", "name": "Rancher", "version": "2.13.2", "arch": "unknown"}]`)
	}))
	defer server.Close()

	// Override base option URL
	options := Options()
	options.URL = server.URL
	conn := Connection(&options)

	products, err := Verify(conn, "rancher", "2.13.2", "unknown", "test-reg-code")
	asserts.Nil(err)
	asserts.NotNil(products)
	asserts.Equal(1, len(products))
	asserts.Equal("Rancher", products[0].Name)
	asserts.Equal("rancher", products[0].Identifier)
}
