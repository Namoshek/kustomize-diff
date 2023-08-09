package kubernetes

import "testing"

func TestFilterUnchangedManifestsDoesNotFilterNewManifests(t *testing.T) {
	manifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: v1\nkind: Service\nmetadata:\n  name: backend\n  namespace: my-namespace",
	}
	manifestHash := manifest.CalculateHash()

	manifests := ManifestMap{
		manifestHash: manifest,
	}

	oldFilteredManifests, newFilteredManifests := FilterUnchangedManifests(make(ManifestMap), manifests)

	if len(oldFilteredManifests) != 0 || len(newFilteredManifests) != 1 {
		t.Fatal("New manifests should not be filtered from the manifest maps.")
	}
}

func TestFilterUnchangedManifestsDoesNotFilterRemovedManifests(t *testing.T) {
	manifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: v1\nkind: Service\nmetadata:\n  name: backend\n  namespace: my-namespace",
	}
	manifestHash := manifest.CalculateHash()

	manifests := ManifestMap{
		manifestHash: manifest,
	}

	oldFilteredManifests, newFilteredManifests := FilterUnchangedManifests(manifests, make(ManifestMap))

	if len(oldFilteredManifests) != 1 || len(newFilteredManifests) != 0 {
		t.Fatal("Removed manifests should not be filtered from the manifest maps.")
	}
}

func TestFilterUnchangedManifestsDoesNotFilterChangedManifests(t *testing.T) {
	oldManifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 1",
	}
	oldManifestHash := oldManifest.CalculateHash()

	oldManifests := ManifestMap{
		oldManifestHash: oldManifest,
	}

	newManifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 2",
	}
	newManifestHash := newManifest.CalculateHash()

	newManifests := ManifestMap{
		newManifestHash: newManifest,
	}

	oldFilteredManifests, newFilteredManifests := FilterUnchangedManifests(oldManifests, newManifests)

	if len(oldFilteredManifests) != 1 || len(newFilteredManifests) != 1 {
		t.Fatal("Changed manifests should not be filtered from the manifest maps.")
	}
}

func TestFilterUnchangedManifestsDoesFilterUnchangedManifests(t *testing.T) {
	oldManifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 1",
	}
	oldManifestHash := oldManifest.CalculateHash()

	oldManifests := ManifestMap{
		oldManifestHash: oldManifest,
	}

	newManifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 1",
	}
	newManifestHash := newManifest.CalculateHash()

	newManifests := ManifestMap{
		newManifestHash: newManifest,
	}

	oldFilteredManifests, newFilteredManifests := FilterUnchangedManifests(oldManifests, newManifests)

	if len(oldFilteredManifests) != 0 || len(newFilteredManifests) != 0 {
		t.Fatal("Unchanged manifests should be filtered from the manifest maps.")
	}
}

func TestGetUniqueManifestHashesReturnsHashesOfFirstMap(t *testing.T) {
	manifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
	}
	manifestHash := manifest.CalculateHash()

	manifests := ManifestMap{
		manifestHash: manifest,
	}

	uniqueManifestHashes := GetUniqueManifestHashes(manifests, make(ManifestMap))

	if len(uniqueManifestHashes) != 1 || uniqueManifestHashes[0] != manifestHash {
		t.Fatal("The unique manifest hashes from the first map should be returned.")
	}
}

func TestGetUniqueManifestHashesReturnsHashesOfSecondMap(t *testing.T) {
	manifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
	}
	manifestHash := manifest.CalculateHash()

	manifests := ManifestMap{
		manifestHash: manifest,
	}

	uniqueManifestHashes := GetUniqueManifestHashes(make(ManifestMap), manifests)

	if len(uniqueManifestHashes) != 1 || uniqueManifestHashes[0] != manifestHash {
		t.Fatal("The unique manifest hashes from the second map should be returned.")
	}
}

func TestGetUniqueManifestHashesReturnsHashesOfBothMapsWithoutDuplicates(t *testing.T) {
	manifest1 := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
	}
	manifest1Hash := manifest1.CalculateHash()

	manifest2 := Manifest{
		ApiVersion: "apps/v1",
		Kind:       "Deployment",
		Name:       "frontend",
		Namespace:  "my-namespace",
	}
	manifest2Hash := manifest2.CalculateHash()

	manifest3 := Manifest{
		ApiVersion: "apps/v1",
		Kind:       "StatefulSet",
		Name:       "database",
	}
	manifest3Hash := manifest3.CalculateHash()

	manifests1 := ManifestMap{
		manifest1Hash: manifest1,
		manifest2Hash: manifest2,
	}

	manifests2 := ManifestMap{
		manifest2Hash: manifest2,
		manifest3Hash: manifest3,
	}

	uniqueManifestHashes := GetUniqueManifestHashes(manifests1, manifests2)

	if len(uniqueManifestHashes) != 3 {
		t.Fatal("The unique manifest hashes should be returned from both maps without duplicates.")
	}
}
