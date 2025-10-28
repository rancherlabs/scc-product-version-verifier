package cmd

import (
	"log"

	"github.com/rancher-sandbox/scc-product-version-verifier/internal/curler"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var RegCode = ""

// curlCmd represents the curl command
var curlCmd = &cobra.Command{
	Use:   "curl-verify [product name] [product version] [product arch : optional]",
	Short: "Mimics using curl to check if SCC product exists",
	Args:  cobra.RangeArgs(2, 3),
	Run: func(cmd *cobra.Command, args []string) {
		productName := args[0]
		productVersion := args[1]
		// Force arch to be optional with default to Rancher value
		productArch := "unknown"
		if len(args) == 3 {
			productArch = args[2]
		}

		logrus.Infof("Verifying product: %s, version: %s, arch: %s", productName, productVersion, productArch)
		logrus.Infof("using product triplet: %s/%s/%s", productName, productVersion, productArch)

		logrus.Infof("Using RegCode: %s", RegCode)

		if err := curler.CurlVerify(productName, productVersion, productArch, RegCode); err != nil {
			log.Fatalf("Error verifying product: %v", err)
		}
	},
}

func init() {
	curlCmd.Flags().StringVarP(&RegCode, "reg-code", "R", "", "The SCC Registration Code used to auth for the API call")

	rootCmd.AddCommand(curlCmd)
}
