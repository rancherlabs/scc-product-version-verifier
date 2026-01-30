package suseconnect

import (
	"encoding/json"
	"fmt"

	"github.com/SUSE/connect-ng/pkg/connection"
	"github.com/SUSE/connect-ng/pkg/registration"
	"github.com/sirupsen/logrus"
)

const (
	ProductsQueryPath = "/connect/subscriptions/products/"
)

func preparePathAndQuery(name, version, arch string) string {
	return fmt.Sprintf("%s?identifier=%s&version=%s&arch=%s", ProductsQueryPath, name, version, arch)
}

func Verify(
	conn *connection.ApiConnection,
	productName,
	version,
	arch,
	regCode string,
) ([]registration.Product, error) {
	path := preparePathAndQuery(productName, version, arch)

	logrus.WithFields(logrus.Fields{
		"product_name": productName,
		"version":      version,
		"arch":         arch,
		"path":         path,
	}).Info("🔨 Building request to query SCC API")

	request, buildErr := conn.BuildRequestRaw("GET", path, nil)
	if buildErr != nil {
		logrus.WithError(buildErr).Error("❌ Error building request")

		return nil, buildErr
	}

	connection.AddRegcodeAuth(request, regCode)

	logrus.WithField("url", request.URL.String()).Info("🚀 Executing request to SCC API")

	responseData, doErr := conn.Do(request)
	if doErr != nil {
		logrus.WithError(doErr).Error("❌ Error executing request")

		return nil, doErr
	}

	logrus.WithField("response_size", len(responseData)).Info("📦 Received response from SCC API")

	products := make([]registration.Product, 0)
	err := json.Unmarshal(responseData, &products)
	if err != nil {
		logrus.WithError(err).Error("❌ Error unmarshalling response")

		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(products) == 0 {
		logrus.Warn("⚠️  No products found matching criteria")

		return nil, fmt.Errorf("product not found")
	}

	logrus.WithField("product_count", len(products)).Info("✅ Successfully retrieved products")

	return products, nil
}
