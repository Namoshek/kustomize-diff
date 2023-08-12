package cmd

import (
	"bytes"
	"errors"
	"net/url"
	"os"

	ado "github.com/namoshek/kustomize-diff/azuredevops"
	k8s "github.com/namoshek/kustomize-diff/kubernetes"
	kustomize "github.com/namoshek/kustomize-diff/kustomize"
	utils "github.com/namoshek/kustomize-diff/utils"

	"github.com/spf13/cobra"

	"go.uber.org/zap"
)

var azuredevopsCmd = NewAzuredevopsCmd()

func NewAzuredevopsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "azuredevops <pathToOldVersion> <pathToNewVersion>",
		Short: "Creates a diff of two Kustomizations and posts it as new comment thread on an Azure DevOps pull request",
		Long:  `Use this action to create a diff of two Kustomizations which should be posted as comment thread on an Azure DevOps pull request.`,
		Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
		Run:   runAzuredevopsCommand,
	}
}

func init() {
	rootCmd.AddCommand(azuredevopsCmd)

	azuredevopsCmd.Flags().StringP("instance", "i", "https://dev.azure.com", "Base URI to the Azure DevOps services or on-premise server")
	azuredevopsCmd.Flags().StringP("organization", "o", "", "The name of the organization (or collection in case of Azure DevOps Server)")
	azuredevopsCmd.Flags().StringP("personal-access-token", "a", "", "The personal access token used for authentication")
	azuredevopsCmd.Flags().StringP("project", "p", "", "The name of the project where the pull request is located")
	azuredevopsCmd.Flags().StringP("repository-id", "r", "", "The name or id of the repository where the pull request is located")
	azuredevopsCmd.Flags().IntP("pull-request-id", "u", 0, "The id of the pull request that should be decorated")
}

func runAzuredevopsCommand(cmd *cobra.Command, args []string) {
	// Parse and validate the command flags.
	azureDevOpsParameters, err := parseAndValidateFlags(cmd)
	if err != nil {
		utils.Logger.Error("Flag validation failed.", zap.Error(err))
		os.Exit(1)
	}

	// Attempt to create the Kustomizations of the provided directories.
	pathToOldVersion, pathToNewVersion := args[0], args[1]

	kustomizeExecutable, err := cmd.Flags().GetString("kustomize-executable")
	if err != nil {
		utils.Logger.Error("Reading --kustomize-executable option failed.")
		os.Exit(1)
	}

	oldKustomization, newKustomization, err := kustomize.BuildKustomizations(kustomizeExecutable, pathToOldVersion, pathToNewVersion)

	// Create a diff of both Kustomizations and print the results.
	buffer := new(bytes.Buffer)
	if err := k8s.CreateAndPrintDiffForManifestFiles(oldKustomization, newKustomization, true, buffer); err != nil {
		utils.Logger.Error("Creating the diff failed.", zap.Error(err))
		os.Exit(1)
	}

	content := buffer.String()
	ado.CreatePullRequestComment(azureDevOpsParameters, content)

	os.Exit(0)
}

// Parses and validates the command flags according to our requirements.
// The flags are only validated structurally and requests may still fail if improper credentials are passed.
func parseAndValidateFlags(cmd *cobra.Command) (*ado.AzureDevOpsParameters, error) {
	// Ensure the instance is a valid URI.
	instance, err := cmd.Flags().GetString("instance")
	if err != nil || instance == "" {
		return nil, errors.New("The provided instance URL is invalid.")
	}
	instanceUrl, err := url.ParseRequestURI(instance)
	if err != nil || instanceUrl == nil {
		return nil, errors.New("The provided instance URL is invalid.")
	}

	// Ensure the organization is filled.
	organization, err := cmd.Flags().GetString("organization")
	if err != nil || organization == "" {
		return nil, errors.New("The provided organization is invalid.")
	}

	// Ensure the personal access token is filled.
	personalAccessToken, err := cmd.Flags().GetString("personal-access-token")
	if err != nil || personalAccessToken == "" {
		return nil, errors.New("The provided personal access token (PAT) is invalid.")
	}

	// Ensure the project is filled.
	project, err := cmd.Flags().GetString("project")
	if err != nil || project == "" {
		return nil, errors.New("The provided project is invalid.")
	}

	// Ensure the repository id is filled.
	repositoryId, err := cmd.Flags().GetString("repository-id")
	if err != nil || repositoryId == "" {
		return nil, errors.New("The provided repository id is invalid.")
	}

	// Ensure the pull request id has the expected format.
	pullRequestId, err := cmd.Flags().GetInt("pull-request-id")
	if err != nil || pullRequestId < 1 {
		return nil, errors.New("The provided pull-request-id is invalid: must be an integer > 0.")
	}

	return &ado.AzureDevOpsParameters{
		Instance:            instance,
		Organization:        organization,
		PersonalAccessToken: personalAccessToken,
		Project:             project,
		PullRequestId:       pullRequestId,
		RepositoryId:        repositoryId,
	}, nil
}
