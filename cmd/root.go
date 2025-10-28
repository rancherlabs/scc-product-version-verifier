package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	versionFlag bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scc-product-version-verifier",
	Short: "A CLI tool to verify product versions in SUSE Customer Center.",
	Long: `scc-product-version-verifier is a command-line tool that interacts with the
SUSE Customer Center (SCC) API to verify if a specific product, version,
and architecture combination exists.`,
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			printVersion()
			os.Exit(0)
		}

		if err := cmd.Help(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "V", false, "Print version information and exit")
}
