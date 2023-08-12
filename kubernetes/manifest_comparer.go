package kubernetes

import (
	"github.com/hashicorp/go-set"
	"golang.org/x/exp/maps"
)

// Filters all common manifests with identical content from the given manifest maps.
// Returns two new maps which contain only the changed (new, altered or removed) manifests.
func FilterUnchangedManifests(oldManifests *ManifestMap, newManifests *ManifestMap) (*ManifestMap, *ManifestMap) {
	filteredOldManifests, filteredNewManifests := make(ManifestMap), make(ManifestMap)

	for key, element := range *oldManifests {
		if manifest, exists := (*newManifests)[key]; exists && manifest.Content == element.Content {
			continue
		}

		filteredOldManifests[key] = element
	}

	for key, element := range *newManifests {
		if manifest, exists := (*oldManifests)[key]; exists && manifest.Content == element.Content {
			continue
		}

		filteredNewManifests[key] = element
	}

	return &filteredOldManifests, &filteredNewManifests
}

// Returns a list of MD5 hashes for all manifests in the given input maps.
// If a manifest with identical MD5 hash is present in both maps, the hash is returned only once.
func GetUniqueManifestHashes(old *ManifestMap, new *ManifestMap) *[]string {
	oldKeys := maps.Keys(*old)
	keys := append(oldKeys, maps.Keys(*new)...)
	slice := set.From[string](keys).Slice()

	return &slice
}
