package cmd

import (
	"fmt"
	"log"

	"github.com/rancher-sandbox/scc-product-version-verifier/internal/curler"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RegCode = ""

const (
	minRequiredFlags = 2
	optionalFlags    = 3
)

// curlCmd represents the curl command.
var curlCmd = &cobra.Command{
	Use:   "curl-verify [product name] [product version] [product arch : optional]",
	Short: "Mimics using curl to check if SCC product exists",
	Args:  cobra.RangeArgs(minRequiredFlags, optionalFlags),
	Run: func(cmd *cobra.Command, args []string) {
		productName := args[0]
		productVersion := args[1]
		// Force arch to be optional with default to Rancher value
		productArch := "unknown"
		if len(args) == optionalFlags {
			productArch = args[2]
		}

		regCode := viper.GetString("regcode")

		logrus.WithFields(map[string]interface{}{
			"product_name":    productName,
			"product_version": productVersion,
			"product_arch":    productArch,
		}).Infof("Checking SCC for product info")
		logrus.WithFields(map[string]interface{}{
			"product_triplet": fmt.Sprintf("%s/%s/%s", productName, productVersion, productArch),
		}).Infof("Using product triplet to check SCC backend")

		logrus.WithField("reg_code", regCode).Trace("Authenticating to SCC API")

		if err := curler.CurlVerify(productName, productVersion, productArch, regCode); err != nil {
			log.Fatalf("Error verifying product: %v", err)
		}
	},
}

func init() {
	curlCmd.Flags().StringP("regcode", "R", "", "The SCC Registration Code used to auth for the API call. Can also be set with the SCC_REGCODE environment variable.")
	viper.BindPFlag("regcode", curlCmd.Flags().Lookup("regcode"))

	rootCmd.AddCommand(curlCmd)
}
