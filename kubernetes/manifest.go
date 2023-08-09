package kubernetes

import (
	"fmt"

	"github.com/namoshek/kustomize-diff/utils"
)

type Manifest struct {
	ApiVersion string
	Kind       string
	Name       string
	Namespace  string
	Content    string
}

func (m Manifest) CalculateHash() string {
	input := fmt.Sprintf("apiVersion: '%s', kind: '%s', name: '%s', namespace: '%s'", m.ApiVersion, m.Kind, m.Name, m.Namespace)
	return utils.CalculateMD5AsString(input)
}
