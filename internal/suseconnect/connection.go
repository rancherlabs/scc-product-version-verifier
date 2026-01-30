package suseconnect

import (
	"github.com/SUSE/connect-ng/pkg/connection"

	"github.com/rancherlabs/scc-product-version-verifier/cmd/version"
)

func Options() connection.Options {
	return connection.DefaultOptions("scc-product-version-verifier", version.Version, "en_US")
}

func Connection(options *connection.Options) *connection.ApiConnection {
	if options == nil {
		defaultOpts := Options()
		options = &defaultOpts
	}

	return connection.New(*options, connection.NoCredentials{})
}
