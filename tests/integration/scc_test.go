package integration

import (
	"os"
	"testing"

	"github.com/rancherlabs/scc-product-version-verifier/internal/curler"
	"github.com/rancherlabs/scc-product-version-verifier/internal/suseconnect"
	"github.com/stretchr/testify/assert"
)

func Test_CurlerWithSCC(t *testing.T) {
	regCode := os.Getenv("SCC_REGCODE")
	if regCode == "" {
		t.Skip("Skipping integration test: SCC_REGCODE environment variable not set")
	}

	err := curler.CurlVerify("rancher-prime", "v2.13.1", "other", regCode)
	assert.Nil(t, err)
}

func Test_VerifyWithSCC(t *testing.T) {
	regCode := os.Getenv("SCC_REGCODE")
	if regCode == "" {
		t.Skip("Skipping integration test: SCC_REGCODE environment variable not set")
	}

	options := suseconnect.Options()
	conn := suseconnect.Connection(&options)

	products, err := suseconnect.Verify(conn, "rancher-prime", "v2.13.1", "other", regCode)
	assert.Nil(t, err)
	assert.NotNil(t, products)
	assert.GreaterOrEqual(t, len(products), 1)
}
