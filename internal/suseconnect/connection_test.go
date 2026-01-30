package suseconnect_test

import (
	"testing"

	"github.com/rancherlabs/scc-product-version-verifier/internal/suseconnect"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	asserts := assert.New(t)
	defaultOptions := suseconnect.Options()
	asserts.NotNil(defaultOptions)
}

func TestConnection(t *testing.T) {
	asserts := assert.New(t)
	defaultConnection := suseconnect.Connection(nil)
	asserts.NotNil(defaultConnection)
}
