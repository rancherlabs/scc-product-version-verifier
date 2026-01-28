package curler

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	TenSecondDuration = 10 * time.Second
)

var (
	ProductQueryUrl        = "https://scc.suse.com/connect/subscriptions/products"
	ProductQueryUrlStaging = "https://stgscc.suse.com/connect/subscriptions/products"
)

func prepareClient() *http.Client {
	// 1. Clone the default transport and apply your TLS configuration
	// We must cast to *http.Transport to access Clone()
	defaultTransport := http.DefaultTransport.(*http.Transport)
	transport := defaultTransport.Clone()

	// Applying InsecureSkipVerify for testing/debugging purposes, as you did.
	// NOTE: Setting InsecureSkipVerify = true is generally UNSAFE for production code.
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// 2. Create the Client with the custom transport
	return &http.Client{
		Transport: transport,
		Timeout:   TenSecondDuration,
	}
}

func CurlVerify(name, version, arch, regCode string) error {
	return curlVerify(ProductQueryUrl, name, version, arch, regCode)
}

func CurlVerifyStaging(name, version, arch, regCode string) error {
	return curlVerify(ProductQueryUrlStaging, name, version, arch, regCode)
}

func curlVerify(apiURL, name, version, arch, regCode string) error {
	queryParams := fmt.Sprintf("?identifier=%s&version=%s&arch=%s", name, version, arch)
	fullURLWithQuery := apiURL + queryParams
	logrus.WithField("query_url", fullURLWithQuery).Info("Prepared URL to query SCC API for product info")

	// 1. Create the Client with the custom transport
	client := prepareClient()

	// 2. Create the Request
	req, err := http.NewRequest(http.MethodGet, fullURLWithQuery, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)

		return err
	}

	// 3. Force the identical headers as curl (CRITICAL STEP)
	req.Header.Set("Authorization", "Token token="+regCode)

	// This User-Agent is what curl used.
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0)")

	// Curl also included this, so we add it for parity.
	req.Header.Set("Accept", "*/*")

	// 4. Execute the Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error executing request: %v\n", err)

		return err
	}
	// Always close the body to reuse the connection
	defer resp.Body.Close()

	// 5. Read the Body (regardless of status code)
	bodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Printf("Error reading response body: %v\n", readErr)
		// We still continue to print status even if body read failed
	}

	// 6. Print Results
	fmt.Printf("✅ Request URL: %s\n", fullURLWithQuery)
	fmt.Printf("✅ Status Code: %d (%s)\n", resp.StatusCode, resp.Status)
	fmt.Printf("✅ Content-Length Header: %s\n", resp.Header.Get("Content-Length"))

	fmt.Println("--- Response Body ---")
	if len(bodyBytes) > 0 {
		fmt.Println(string(bodyBytes))
	} else {
		fmt.Println("Body is empty (zero length).")
		fmt.Println("Go's Response.Body was likely http.noBody or read returned EOF immediately.")
	}
	fmt.Println("---------------------")

	products := make([]interface{}, 1)
	json.Unmarshal(bodyBytes, &products)

	if len(products) == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
