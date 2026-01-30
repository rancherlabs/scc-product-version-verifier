package cmd

import (
	"fmt"
	"log"

	"github.com/rancherlabs/scc-product-version-verifier/internal/suseconnect"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// verifyCmd represents the verify command.
var verifyCmd = &cobra.Command{
	Use:   "verify [product name] [product version] [product arch : optional]",
	Short: "Uses connect-ng to check if SCC product exists",
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
		useStaging := viper.GetBool("staging")
		logrus.Infof("use staging set to: %v", useStaging)

		logrus.WithFields(map[string]interface{}{
			"product_name":    productName,
			"product_version": productVersion,
			"product_arch":    productArch,
		}).Infof("Checking SCC for product info")
		logrus.WithFields(map[string]interface{}{
			"product_triplet": fmt.Sprintf("%s/%s/%s", productName, productVersion, productArch),
		}).Infof("Using product triplet to check SCC backend")

		logrus.WithField("reg_code", regCode).Trace("Authenticating to SCC API")

		options := suseconnect.Options()
		if useStaging {
			options.URL = "https://stgscc.suse.com"
		}
		conn := suseconnect.Connection(&options)

		productList, err := suseconnect.Verify(conn, productName, productVersion, productArch, regCode)
		if err != nil {
			log.Fatalf("Error verifying product: %v", err)
		}

		if len(productList) > 0 {
			logrus.Info("--- Response Products ---")
			for _, product := range productList {
				logrus.WithFields(logrus.Fields{
					"name":       product.Name,
					"identifier": product.Identifier,
					"version":    product.Version,
					"arch":       product.Arch,
				}).Info("✅ Product found")
			}
			logrus.Info("-------------------------")
		} else {
			logrus.Warn("⚠️  No products found")
		}
	},
}

func init() {
	verifyCmd.Flags().BoolP("staging", "S", false, "Use the SCC Staging API instead of Production")
	viper.BindPFlag("staging", verifyCmd.Flags().Lookup("staging"))
	verifyCmd.Flags().StringP("regcode", "R", "", "The SCC Registration Code used to auth for the API call. Can also be set with the SCC_REGCODE environment variable.")
	viper.BindPFlag("regcode", verifyCmd.Flags().Lookup("regcode"))

	rootCmd.AddCommand(verifyCmd)
}
