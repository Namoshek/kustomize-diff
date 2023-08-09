package cmd

import (
	"os"

	k8s "github.com/namoshek/kustomize-diff/kubernetes"
	kustomize "github.com/namoshek/kustomize-diff/kustomize"
	utils "github.com/namoshek/kustomize-diff/utils"

	"github.com/spf13/cobra"

	"go.uber.org/zap"
)

// inlineCmd represents the inline command
var inlineCmd = NewInlineCmd()

func NewInlineCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "inline <pathToOldVersion> <pathToNewVersion>",
		Short: "Creates an inline diff of two Kustomizations",
		Long:  `Use this action for a quick inline diff of two Kustomizations.`,
		Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
		Run:   runCommand,
	}
}

func init() {
	rootCmd.AddCommand(inlineCmd)

	inlineCmd.Flags().StringP("kustomize-executable", "k", "kustomize", "Path to the kustomize binary")
}

func runCommand(cmd *cobra.Command, args []string) {
	// Ensure the given Kustomization directories exist.
	utils.Logger.Debug("Checking existence of given Kustomzation directories.")

	pathToOldVersion, pathToNewVersion := args[0], args[1]

	if _, err := os.Stat(pathToOldVersion); os.IsNotExist(err) {
		utils.Logger.Error("Directory '" + pathToOldVersion + "' does not exist.")
		os.Exit(1)
	}

	if _, err := os.Stat(pathToNewVersion); os.IsNotExist(err) {
		utils.Logger.Error("Directory '" + pathToNewVersion + "' does not exist.")
		os.Exit(1)
	}

	// Build the Kustomizations in a safe way.
	utils.Logger.Debug("Building Kustomizations for both version.")

	kustomizeExecutable, err := cmd.Flags().GetString("kustomize-executable")
	if err != nil {
		utils.Logger.Error("Reading --kustomize-executable option failed.")
		os.Exit(2)
	}

	oldKustomization, err := kustomize.BuildKustomization(kustomizeExecutable, pathToOldVersion)
	if err != nil {
		utils.Logger.Error("Building the Kustomization for '" + pathToOldVersion + "' failed.")
		os.Exit(2)
	}

	newKustomization, err := kustomize.BuildKustomization(kustomizeExecutable, pathToNewVersion)
	if err != nil {
		utils.Logger.Error("Building the Kustomization for '" + pathToNewVersion + "' failed.")
		os.Exit(2)
	}

	// Create a diff of both Kustomizations and print the results.
	if err := createAndPrintDiff(oldKustomization, newKustomization); err != nil {
		utils.Logger.Error("Creating the diff failed.", zap.Error(err))
		os.Exit(3)
	}

	os.Exit(0)
}

func createAndPrintDiff(old string, new string) error {
	// Parse the Kustomizations into individual manifests for easier comparison.
	oldManifests, err := k8s.SplitKustomizationIntoManifests(old)
	if err != nil {
		return err
	}

	newManifests, err := k8s.SplitKustomizationIntoManifests(new)
	if err != nil {
		return err
	}

	// Remove all unchanged manifests as we do not need to process them further.
	oldManifests, newManifests = k8s.FilterUnchangedManifests(oldManifests, newManifests)

	// Retrieve all unique manifest hashes and iterate them to print the diff per manifest.
	manifestHashes := k8s.GetUniqueManifestHashes(oldManifests, newManifests)

	for _, hash := range manifestHashes {
		oldManifest, newManifest := oldManifests[hash], newManifests[hash]

		k8s.CreateAndPrintDiffForManifests(oldManifest, newManifest, true, os.Stdout)
	}

	return nil
}
