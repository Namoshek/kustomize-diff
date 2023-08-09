package kustomize

import (
	"errors"
	"os/exec"

	"github.com/namoshek/kustomize-diff/utils"
)

// Builds the Kustomization for the given path using the given Kustomize executable.
func BuildKustomization(kustomizeExecutable string, path string) (string, error) {
	utils.Logger.Debug("Building Kustomization for '" + path + "'.")

	out, err := exec.Command(kustomizeExecutable, "build", path).Output()

	if err != nil {
		return "", errors.Join(errors.New("Building Kustomization for '"+path+"' failed."), err)
	}

	return string(out), nil
}
