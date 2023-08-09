package kubernetes

import (
	"fmt"

	"github.com/namoshek/kustomize-diff/utils"
)

// A manifest describes a Kubernetes object with the most important parameters and its content.
type Manifest struct {
	ApiVersion string
	Kind       string
	Name       string
	Namespace  string
	Content    string
}

// This type contains the manifest hash as key and the manifest itself as value.
type ManifestMap map[string]Manifest

// Calculates the MD5 hash for a manifest header (apiVersion, kind, name and namespace).
// Can be used to compare manifests.
func (m Manifest) CalculateHash() string {
	input := fmt.Sprintf("apiVersion: '%s', kind: '%s', name: '%s', namespace: '%s'", m.ApiVersion, m.Kind, m.Name, m.Namespace)
	return utils.CalculateMD5AsString(input)
}
