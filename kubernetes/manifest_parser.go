package kubernetes

import (
	"errors"
	"strings"

	"github.com/namoshek/kustomize-diff/utils"
	
	"gopkg.in/yaml.v3"
)

func SplitKustomizationIntoManifests(kustomization string) (map[string]Manifest, error) {
	kustomization = strings.ReplaceAll(kustomization, "\r\n", "\n")
	parts := strings.Split(kustomization, "---\n")

	result := make(map[string]Manifest)
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

func parseManifest(content string) (Manifest, error) {
	var data map[string]interface{}
	err := yaml.Unmarshal([]byte(content), &data)

	if err != nil {
		return Manifest{}, errors.Join(errors.New("Parsing manifest to retrieve headers failed."), err)
	}

	return Manifest{
		ApiVersion: utils.GetMapValueOrDefault(data, "apiVersion", "").(string),
		Kind:       utils.GetMapValueOrDefault(data, "kind", "").(string),
		Name:       utils.GetMapValueOrDefault(utils.GetMapValueOrDefault(data, "metadata", make(map[string]interface{})).(map[string]interface{}), "name", "").(string),
		Namespace:  utils.GetMapValueOrDefault(utils.GetMapValueOrDefault(data, "metadata", make(map[string]interface{})).(map[string]interface{}), "namespace", "").(string),
		Content:    content,
	}, nil
}
