package kubernetes

import (
	"fmt"
	"io"
	"strings"

	"github.com/kylelemons/godebug/diff"
)

// Creates and prints the diff for two manifest files, each containing multiple manifests separated by the YAML separator '---'.
func CreateAndPrintDiffForManifestFiles(old *string, new *string, formatAsMarkdownCodeBlock bool, output io.Writer) error {
	// Parse the Kustomizations into individual manifests for easier comparison.
	oldManifests, err := SplitKustomizationIntoManifests(old)
	if err != nil {
		return err
	}

	newManifests, err := SplitKustomizationIntoManifests(new)
	if err != nil {
		return err
	}

	// Remove all unchanged manifests as we do not need to process them further.
	oldManifests, newManifests = FilterUnchangedManifests(oldManifests, newManifests)

	// Retrieve all unique manifest hashes and iterate them to print the diff per manifest.
	manifestHashes := GetUniqueManifestHashes(oldManifests, newManifests)

	for _, hash := range *manifestHashes {
		oldManifest, newManifest := (*oldManifests)[hash], (*newManifests)[hash]

		CreateAndPrintDiffForManifests(oldManifest, newManifest, formatAsMarkdownCodeBlock, output)
	}

	return nil
}

// Creates and prints the diff for two manifests.
func CreateAndPrintDiffForManifests(old Manifest, new Manifest, formatAsMarkdownCodeBlock bool, output io.Writer) {
	if formatAsMarkdownCodeBlock {
		fmt.Fprintln(output, "```diff")
	}

	if old.Content == new.Content {
		fmt.Fprintln(output, old.Content)

		return
	}

	diff := diff.Diff(old.Content, new.Content)

	// In case a new manifest is given, the old content is diffed as empty line, which we remove.
	diff = strings.TrimPrefix(diff, "-\n")

	// The diff may contain an empty line at the end, which we remove.
	diff = strings.TrimSuffix(diff, "\n ")

	// In case a removed manifest is given, the new content is diffed as empty line, which we remove.
	diff = strings.TrimSuffix(diff, "\n+")

	fmt.Fprintln(output, diff)

	if formatAsMarkdownCodeBlock {
		fmt.Fprintln(output, "```")
	}
}
