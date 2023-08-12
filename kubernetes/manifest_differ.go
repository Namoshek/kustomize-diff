package kubernetes

import (
	"strings"

	"github.com/kylelemons/godebug/diff"
)

// The diff between two manifests.
type ManifestDiff struct {
	OldManifest *Manifest
	NewManifest *Manifest
	Diff        string
}

// Creates the diff for two manifest files, each containing multiple manifests separated by the YAML separator '---'.
func CreateDiffForManifestFiles(old *string, new *string) ([]ManifestDiff, error) {
	// Parse the Kustomizations into individual manifests for easier comparison.
	oldManifests, err := SplitKustomizationIntoManifests(old)
	if err != nil {
		return nil, err
	}

	newManifests, err := SplitKustomizationIntoManifests(new)
	if err != nil {
		return nil, err
	}

	// Remove all unchanged manifests as we do not need to process them further.
	oldManifests, newManifests = FilterUnchangedManifests(oldManifests, newManifests)

	// Retrieve all unique manifest hashes and iterate them to create the diff per manifest.
	manifestHashes := GetUniqueManifestHashes(oldManifests, newManifests)

	var diffs []ManifestDiff
	for _, hash := range *manifestHashes {
		oldManifest, newManifest := (*oldManifests)[hash], (*newManifests)[hash]

		diff := CreateDiffForManifests(&oldManifest, &newManifest)
		diffs = append(diffs, *diff)
	}

	return diffs, nil
}

// Creates the diff for two manifests.
func CreateDiffForManifests(old *Manifest, new *Manifest) *ManifestDiff {

	if old.Content == new.Content {
		return &ManifestDiff{
			OldManifest: old,
			NewManifest: new,
			Diff:        old.Content,
		}
	}

	diff := diff.Diff(old.Content, new.Content)

	// In case a new manifest is given, the old content is diffed as empty line, which we remove.
	diff = strings.TrimPrefix(diff, "-\n")

	// The diff may contain an empty line at the end, which we remove.
	diff = strings.TrimSuffix(diff, "\n ")

	// In case a removed manifest is given, the new content is diffed as empty line, which we remove.
	diff = strings.TrimSuffix(diff, "\n+")

	return &ManifestDiff{
		OldManifest: old,
		NewManifest: new,
		Diff:        diff,
	}
}
