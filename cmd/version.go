package cmd

import (
	"fmt"

	"github.com/rancher-sandbox/scc-product-version-verifier/cmd/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of scc-product-version-verifier",
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion() {
	fmt.Printf("`scc-product-version-verifier` - version: %s, commit: %s, built at: %s\n", version.Version, version.GitCommit, version.Date)
}
