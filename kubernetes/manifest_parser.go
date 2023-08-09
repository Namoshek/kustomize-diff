package kubernetes

import (
	"errors"
	"strings"

	"gopkg.in/yaml.v3"
)

// This type is a convenience layer on top of a generic map and represents a YAML object.
type YamlObject map[string]interface{}

// Splits the given Kustomization into individual manifests per object.
func SplitKustomizationIntoManifests(kustomization string) (ManifestMap, error) {
	kustomization = strings.ReplaceAll(kustomization, "\r\n", "\n")
	parts := strings.Split(kustomization, "---\n")

	result := make(ManifestMap)
	for i := range parts {
		manifest, err := parseManifest(parts[i])
		if err != nil {
			return nil, errors.Join(errors.New("Parsing manifest failed."), err)
		}

		hash := manifest.CalculateHash()

		result[hash] = manifest
	}
	return result, nil
}

// Parses the given string as Kubernetes manifest. The parsing will fail if apiVersion, kind or metadata.name is missing.
func parseManifest(content string) (Manifest, error) {
	var data YamlObject
	err := yaml.Unmarshal([]byte(content), &data)

	if err != nil {
		return Manifest{}, errors.Join(errors.New("Parsing manifest to retrieve headers failed."), err)
	}

	apiVersion := data.getMapValueOrDefault("apiVersion", "").(string)
	if apiVersion == "" {
		return Manifest{}, errors.New("No valid 'apiVersion' found in manifest.")
	}

	kind := data.getMapValueOrDefault("kind", "").(string)
	if kind == "" {
		return Manifest{}, errors.New("No valid 'kind' found in manifest.")
	}

	name := data.getMapValueOrDefault("metadata", make(YamlObject)).(YamlObject).getMapValueOrDefault("name", "").(string)
	if name == "" {
		return Manifest{}, errors.New("No valid 'metadata.name' found in manifest.")
	}

	namespace := data.getMapValueOrDefault("metadata", make(YamlObject)).(YamlObject).getMapValueOrDefault("namespace", "").(string)

	return Manifest{
		ApiVersion: apiVersion,
		Kind:       kind,
		Name:       name,
		Namespace:  namespace,
		Content:    content,
	}, nil
}

// Retrieves the value for a given key from the given YAML object.
func (o YamlObject) getMapValueOrDefault(key string, defaultValue interface{}) interface{} {
	if x, found := o[key]; found {
		return x
	}

	return defaultValue
}
