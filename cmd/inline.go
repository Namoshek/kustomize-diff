package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	k8s "github.com/namoshek/kustomize-diff/kubernetes"
	utils "github.com/namoshek/kustomize-diff/utils"

	"github.com/hashicorp/go-set"
	"github.com/kylelemons/godebug/diff"
	"github.com/spf13/cobra"

	"golang.org/x/exp/maps"

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

	oldKustomization, err := kustomizeBuild(kustomizeExecutable, pathToOldVersion)
	if err != nil {
		utils.Logger.Error("Building the Kustomization for '" + pathToOldVersion + "' failed.")
		os.Exit(2)
	}

	newKustomization, err := kustomizeBuild(kustomizeExecutable, pathToNewVersion)
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

func kustomizeBuild(kustomizeExecutable string, path string) (string, error) {
	utils.Logger.Debug("Building Kustomization for '" + path + "'.")

	out, err := exec.Command(kustomizeExecutable, "build", path).Output()

	if err != nil {
		return "", errors.Join(errors.New("Building Kustomization for '"+path+"' failed."), err)
	}

	return string(out), nil
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
	oldManifests, newManifests = filterUnchangedManifests(oldManifests, newManifests)

	// Retrieve all unique manifest hashes and iterate them to print the diff per manifest.
	manifestHashes := getUniqueManifestHashes(oldManifests, newManifests)

	for _, hash := range manifestHashes {
		oldManifest, newManifest := oldManifests[hash], newManifests[hash]

		createAndPrintDiffForManifest(oldManifest, newManifest)
	}

	return nil
}

func createAndPrintDiffForManifest(old k8s.Manifest, new k8s.Manifest) {
	header := old
	if header == (k8s.Manifest{}) {
		header = new
	}

	fmt.Println("```diff")

	diff := diff.Diff(old.Content, new.Content)
	diff = strings.TrimSuffix(diff, "\n ")
	fmt.Println(diff)

	fmt.Println("```")
}

func filterUnchangedManifests(oldManifests map[string]k8s.Manifest, newManifests map[string]k8s.Manifest) (map[string]k8s.Manifest, map[string]k8s.Manifest) {
	filteredOldManifests, filteredNewManifests := make(map[string]k8s.Manifest), make(map[string]k8s.Manifest)

	for key, element := range oldManifests {
		if manifest, exists := newManifests[key]; exists && manifest.Content == element.Content {
			continue
		}

		filteredOldManifests[key] = element
	}

	for key, element := range newManifests {
		if manifest, exists := oldManifests[key]; exists && manifest.Content == element.Content {
			continue
		}

		filteredNewManifests[key] = element
	}

	return filteredOldManifests, filteredNewManifests
}

func getUniqueManifestHashes(old map[string]k8s.Manifest, new map[string]k8s.Manifest) []string {
	oldKeys := maps.Keys(old)
	keys := append(oldKeys, maps.Keys(new)...)

	return set.From[string](keys).Slice()
}
