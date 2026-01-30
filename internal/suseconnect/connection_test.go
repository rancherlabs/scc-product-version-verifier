package suseconnect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	asserts := assert.New(t)
	defaultOptions := Options()
	asserts.NotNil(defaultOptions)
}

func TestConnection(t *testing.T) {
	asserts := assert.New(t)
	defaultConnection := Connection(nil)
	asserts.NotNil(defaultConnection)
}
