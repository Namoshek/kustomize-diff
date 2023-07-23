package cmd

import (
	"os"

	"github.com/namoshek/kustomize-diff/utils"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kustomize-diff",
	Short: "Diff two Kustomizations (e.g. for PR review)",
	Long:  `This command can create a diff for two Kustomization directories.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")

		utils.InitializeLogger(verbose)
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
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Print verbose output during execution")
}
