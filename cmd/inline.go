package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-set"
	"github.com/kylelemons/godebug/diff"
	"github.com/spf13/cobra"

	"golang.org/x/exp/maps"

	"gopkg.in/yaml.v3"
)

type manifest struct {
	apiVersion string
	kind       string
	name       string
	namespace  string
	content    string
}

// inlineCmd represents the inline command
var inlineCmd = &cobra.Command{
	Use:   "inline <pathToOldVersion> <pathToNewVersion>",
	Short: "Creates an inline diff of two Kustomizations",
	Long:  `Use this action for a quick inline diff of two Kustomizations.`,
	Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	Run: runCommand,
}

func init() {
	rootCmd.AddCommand(inlineCmd)

	inlineCmd.Flags().StringP("kustomize-executable", "k", "kustomize", "Path to the kustomize binary")
}

func runCommand(cmd *cobra.Command, args []string) {
	// Ensure the given Kustomization directories exist.
	printVerbose(cmd, "Checking existence of given Kustomzation directories.")

	pathToOldVersion, pathToNewVersion := args[0], args[1]

	if _, err := os.Stat(pathToOldVersion); os.IsNotExist(err) {
		fmt.Println("Directory '" + pathToOldVersion + "' does not exist.")
		os.Exit(1)
	}

	if _, err := os.Stat(pathToNewVersion); os.IsNotExist(err) {
		fmt.Println("Directory '" + pathToNewVersion + "' does not exist.")
		os.Exit(1)
	}

	// Build the Kustomizations in a safe way.
	printVerbose(cmd, "Building Kustomizations for both version.")

	oldKustomization, err := kustomizeBuild(cmd, pathToOldVersion)
	if err != nil {
		fmt.Println("Building the Kustomization for '" + pathToOldVersion + "' failed.")
		os.Exit(2)
	}

	newKustomization, err := kustomizeBuild(cmd, pathToNewVersion)
	if err != nil {
		fmt.Println("Building the Kustomization for '" + pathToNewVersion + "' failed.")
		os.Exit(2)
	}

	// Create a diff of both Kustomizations and print the results.
	if err := createAndPrintDiff(oldKustomization, newKustomization); err != nil {
		fmt.Println("Creating the diff failed:")
		fmt.Println(err)
		os.Exit(3)
	}

	os.Exit(0)
}

func kustomizeBuild(cmd *cobra.Command, path string) (string, error) {
	printVerbose(cmd, "Building Kustomization for '"+path+"'.")

	kustomizeExecutable, _ := cmd.Flags().GetString("kustomize-executable")
	out, err := exec.Command(kustomizeExecutable, "build", path).Output()

	if err != nil {
		return "", errors.Join(errors.New("Building Kustomization for '"+path+"' failed."), err)
	}

	return string(out), nil
}

func createAndPrintDiff(old string, new string) error {
	// Parse the Kustomizations into individual manifests for easier comparison.
	oldManifests, err := splitKustomizationIntoManifests(old)
	if err != nil {
		return err
	}

	newManifests, err := splitKustomizationIntoManifests(new)
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

func createAndPrintDiffForManifest(old manifest, new manifest) {
	header := old
	if header == (manifest{}) {
		header = new
	}

	fmt.Println("```diff")

	diff := diff.Diff(old.content, new.content)
	diff = strings.TrimSuffix(diff, "\n ")
	fmt.Println(diff)

	fmt.Println("```")
}

func splitKustomizationIntoManifests(kustomization string) (map[string]manifest, error) {
	kustomization = strings.ReplaceAll(kustomization, "\r\n", "\n")
	parts := strings.Split(kustomization, "---\n")

	result := make(map[string]manifest)
	for i := range parts {
		manifest, err := parseManifest(parts[i])
		if err != nil {
			return nil, errors.Join(errors.New("Parsing manifest failed."), err)
		}

		hash := calculateHash(manifest)

		result[hash] = manifest
	}
	return result, nil
}

func parseManifest(content string) (manifest, error) {
	var data map[string]interface{}
	err := yaml.Unmarshal([]byte(content), &data)

	if err != nil {
		return manifest{}, errors.Join(errors.New("Parsing manifest to retrieve headers failed."), err)
	}

	return manifest{
		apiVersion: getMapValueOrDefault(data, "apiVersion", "").(string),
		kind:       getMapValueOrDefault(data, "kind", "").(string),
		name:       getMapValueOrDefault(getMapValueOrDefault(data, "metadata", make(map[string]interface{})).(map[string]interface{}), "name", "").(string),
		namespace:  getMapValueOrDefault(getMapValueOrDefault(data, "metadata", make(map[string]interface{})).(map[string]interface{}), "namespace", "").(string),
		content:    content,
	}, nil
}

func calculateHash(manifest manifest) string {
	input := fmt.Sprintf("apiVersion: '%s', kind: '%s', name: '%s', namespace: '%s'", manifest.apiVersion, manifest.kind, manifest.name, manifest.namespace)
	return calculateMd5Hash(input)
}

func filterUnchangedManifests(oldManifests map[string]manifest, newManifests map[string]manifest) (map[string]manifest, map[string]manifest) {
	filteredOldManifests, filteredNewManifests := make(map[string]manifest), make(map[string]manifest)

	for key, element := range oldManifests {
		if manifest, exists := newManifests[key]; exists && manifest.content == element.content {
			continue
		}

		filteredOldManifests[key] = element
	}

	for key, element := range newManifests {
		if manifest, exists := oldManifests[key]; exists && manifest.content == element.content {
			continue
		}

		filteredNewManifests[key] = element
	}

	return filteredOldManifests, filteredNewManifests
}

func getUniqueManifestHashes(old map[string]manifest, new map[string]manifest) []string {
	oldKeys := maps.Keys(old)
	keys := append(oldKeys, maps.Keys(new)...)

	return set.From[string](keys).Slice()
}

func printVerbose(cmd *cobra.Command, text string) {
	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		fmt.Println(text)
	}
}

func getMapValueOrDefault(dict map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if x, found := dict[key]; found {
		return x
	}

	return defaultValue
}

func calculateMd5Hash(text string) string {
	hash := md5.Sum([]byte(text))

	return hex.EncodeToString(hash[:])
}
