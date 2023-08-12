package cmd

import (
	"os"

	k8s "github.com/namoshek/kustomize-diff/kubernetes"
	kustomize "github.com/namoshek/kustomize-diff/kustomize"
	utils "github.com/namoshek/kustomize-diff/utils"

	"github.com/spf13/cobra"

	"go.uber.org/zap"
)

var inlineCmd = NewInlineCmd()

func NewInlineCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "inline <pathToOldVersion> <pathToNewVersion>",
		Short: "Creates an inline diff of two Kustomizations",
		Long:  `Use this action for a quick inline diff of two Kustomizations.`,
		Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
		Run:   runInlineCommand,
	}
}

func init() {
	rootCmd.AddCommand(inlineCmd)
}

func runInlineCommand(cmd *cobra.Command, args []string) {
	// Attempt to create the Kustomizations of the provided directories.
	pathToOldVersion, pathToNewVersion := args[0], args[1]

	kustomizeExecutable, err := cmd.Flags().GetString("kustomize-executable")
	if err != nil {
		utils.Logger.Error("Reading --kustomize-executable option failed.")
		os.Exit(1)
	}

	oldKustomization, newKustomization, err := kustomize.BuildKustomizations(kustomizeExecutable, pathToOldVersion, pathToNewVersion)

	// Create a diff of both Kustomizations.
	diffs, err := k8s.CreateDiffForManifestFiles(oldKustomization, newKustomization)
	if err != nil {
		utils.Logger.Error("Creating the diff failed.", zap.Error(err))
		os.Exit(1)
	}

	// Iterate the diffs and print them to stdout.
	for _, diff := range diffs {
		k8s.PrintDiff(&diff, true, os.Stdout)
	}

	os.Exit(0)
}
