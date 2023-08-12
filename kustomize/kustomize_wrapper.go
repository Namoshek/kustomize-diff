package kustomize

import (
	"errors"
	"os"
	"os/exec"

	"github.com/namoshek/kustomize-diff/utils"
)

// Builds the Kustomizations for the given paths using the given Kustomize executable.
func BuildKustomizations(kustomizeExecutable string, pathToOldVersion string, pathToNewVersion string) (*string, *string, error) {
	// Ensure the given Kustomization directories exist.
	utils.Logger.Debug("Checking existence of given Kustomzation directories.")

	if _, err := os.Stat(pathToOldVersion); os.IsNotExist(err) {
		return nil, nil, errors.Join(errors.New("Directory '"+pathToOldVersion+"' does not exist."), err)
	}

	if _, err := os.Stat(pathToNewVersion); os.IsNotExist(err) {
		return nil, nil, errors.Join(errors.New("Directory '"+pathToNewVersion+"' does not exist."), err)
	}

	// Build the Kustomizations in a safe way.
	utils.Logger.Debug("Building Kustomizations for both version.")

	oldKustomization, err := buildKustomization(kustomizeExecutable, pathToOldVersion)
	if err != nil {
		return nil, nil, errors.Join(errors.New("Building the Kustomization for '"+pathToOldVersion+"' failed."), err)
	}

	newKustomization, err := buildKustomization(kustomizeExecutable, pathToNewVersion)
	if err != nil {
		return nil, nil, errors.Join(errors.New("Building the Kustomization for '"+pathToNewVersion+"' failed."), err)
	}

	return &oldKustomization, &newKustomization, nil
}

// Builds the Kustomization for the given path using the given Kustomize executable.
func buildKustomization(kustomizeExecutable string, path string) (string, error) {
	utils.Logger.Debug("Building Kustomization for '" + path + "'.")

	out, err := exec.Command(kustomizeExecutable, "build", path).Output()

	if err != nil {
		return "", errors.Join(errors.New("Building Kustomization for '"+path+"' failed."), err)
	}

	return string(out), nil
}
