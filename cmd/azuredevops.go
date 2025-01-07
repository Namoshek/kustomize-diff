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

type AzureDevOpsCommandFlags struct {
	CommentPerResource   bool
	HideDiffInSpoiler    bool
	PrependedCommentText string
	AppendedCommentText  string
}

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
	azuredevopsCmd.Flags().StringP("personal-access-token", "a", "", "The personal access token used for authentication")
	azuredevopsCmd.Flags().StringP("organization", "o", "", "The name of the organization (or collection in case of Azure DevOps Server)")
	azuredevopsCmd.Flags().StringP("project", "p", "", "The name of the project where the pull request is located")
	azuredevopsCmd.Flags().StringP("repository-id", "r", "", "The name or id of the repository where the pull request is located")
	azuredevopsCmd.Flags().IntP("pull-request-id", "u", 0, "The id of the pull request that should be decorated")
	azuredevopsCmd.Flags().Bool("comment-per-resource", false, "Create a separate comment for each resource with differences")
	azuredevopsCmd.Flags().Bool("hide-diff-in-spoiler", false, "Add a spoiler around diffs to prevent displaying large comments")
	azuredevopsCmd.Flags().String("prepended-comment-text", "", "Text to prepend to created pull request comments; it is added before and outside spoilers if enabled")
	azuredevopsCmd.Flags().String("appended-comment-text", "", "Text to append to created pull request comments; it is added after and outside spoilers if enabled")
}

func runAzuredevopsCommand(cmd *cobra.Command, args []string) {
	// Parse and validate the command flags.
	azureDevOpsParameters, azureDevOpsCommandFlags, err := parseAndValidateFlags(cmd)
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
	if err != nil {
		utils.Logger.Error("Building Kustomizations failed.", zap.Error(err))
		os.Exit(1)
	}

	// Create a diff of both Kustomizations.
	diffs, err := k8s.CreateDiffForManifestFiles(oldKustomization, newKustomization)
	if err != nil {
		utils.Logger.Error("Creating the diff failed.", zap.Error(err))
		os.Exit(1)
	}

	if len(diffs) == 0 {
		utils.Logger.Debug("No diff found, exiting.")
		os.Exit(0)
	}

	// Prepare diff slices to process depending on the command flags.
	var diffSlices [][]k8s.ManifestDiff
	if azureDevOpsCommandFlags.CommentPerResource {
		for _, diff := range diffs {
			diffSlices = append(diffSlices, []k8s.ManifestDiff{diff})
		}
	} else {
		diffSlices = append(diffSlices, diffs)
	}

	// Process the diff slices one-by-one.
	for _, diffSlice := range diffSlices {
		err = createPullRequestCommentForManifests(diffSlice, azureDevOpsParameters, azureDevOpsCommandFlags)
		if err != nil {
			utils.Logger.Error("Creating pull request comment for diff slice failed.", zap.Error(err))
			os.Exit(1)
		}
	}

	os.Exit(0)
}

// Creates a pull request comment with the given diffs.
func createPullRequestCommentForManifests(diffs []k8s.ManifestDiff, azureDevOpsParameters *ado.AzureDevOpsParameters, azureDevOpsCommandFlags *AzureDevOpsCommandFlags) error {
	// Iterate the diffs and print them into a buffer.
	diffBuffer := new(bytes.Buffer)
	for _, diff := range diffs {
		k8s.PrintDiff(&diff, true, diffBuffer)
	}

	diffContent := diffBuffer.String()
	if azureDevOpsCommandFlags.HideDiffInSpoiler {
		diffContent = wrapContentInSpoiler(diffContent)
	}

	// Use the buffer to create a comment on the pull request.
	contentBuffer := bytes.NewBufferString("")

	err := writeContentToBuffer(contentBuffer, diffContent, azureDevOpsCommandFlags)
	if err != nil {
		return errors.Join(errors.New("Writing content to buffer failed."), err)
	}

	err = ado.CreatePullRequestComment(azureDevOpsParameters, contentBuffer.String())
	if err != nil {
		return errors.Join(errors.New("Creating pull request comment failed."), err)
	}

	return nil
}

// Parses and validates the command flags according to our requirements.
// The flags are only validated structurally and requests may still fail if improper credentials are passed.
func parseAndValidateFlags(cmd *cobra.Command) (*ado.AzureDevOpsParameters, *AzureDevOpsCommandFlags, error) {
	// Ensure the instance is a valid URI.
	instance, err := cmd.Flags().GetString("instance")
	if err != nil || instance == "" {
		return nil, nil, errors.New("The provided instance URL is invalid.")
	}
	instanceUrl, err := url.ParseRequestURI(instance)
	if err != nil || instanceUrl == nil {
		return nil, nil, errors.New("The provided instance URL is invalid.")
	}

	// Ensure the organization is filled.
	organization, err := cmd.Flags().GetString("organization")
	if err != nil || organization == "" {
		return nil, nil, errors.New("The provided organization is invalid.")
	}

	// Ensure the personal access token is filled.
	personalAccessToken, err := cmd.Flags().GetString("personal-access-token")
	if err != nil || personalAccessToken == "" {
		return nil, nil, errors.New("The provided personal access token (PAT) is invalid.")
	}

	// Ensure the project is filled.
	project, err := cmd.Flags().GetString("project")
	if err != nil || project == "" {
		return nil, nil, errors.New("The provided project is invalid.")
	}

	// Ensure the repository id is filled.
	repositoryId, err := cmd.Flags().GetString("repository-id")
	if err != nil || repositoryId == "" {
		return nil, nil, errors.New("The provided repository id is invalid.")
	}

	// Ensure the pull request id has the expected format.
	pullRequestId, err := cmd.Flags().GetInt("pull-request-id")
	if err != nil || pullRequestId < 1 {
		return nil, nil, errors.New("The provided pull-request-id is invalid: must be an integer > 0.")
	}

	// Ensure boolean flags are valid.
	commentPerResource, err := cmd.Flags().GetBool("comment-per-resource")
	if err != nil {
		return nil, nil, errors.New("The provided comment-per-resource is invalid.")
	}

	hideDiffInSpoiler, err := cmd.Flags().GetBool("hide-diff-in-spoiler")
	if err != nil {
		return nil, nil, errors.New("The provided hide-diff-in-spoiler is invalid.")
	}

	// Ensure optional texts were passed successfully.
	prependedCommentText, err := cmd.Flags().GetString("prepended-comment-text")
	if err != nil {
		return nil, nil, errors.New("The provided prepended-comment-text is invalid.")
	}

	appendedCommentText, err := cmd.Flags().GetString("appended-comment-text")
	if err != nil {
		return nil, nil, errors.New("The provided appended-comment-text is invalid.")
	}

	return &ado.AzureDevOpsParameters{
			Instance:            instance,
			Organization:        organization,
			PersonalAccessToken: personalAccessToken,
			Project:             project,
			PullRequestId:       pullRequestId,
			RepositoryId:        repositoryId,
		},
		&AzureDevOpsCommandFlags{
			CommentPerResource:   commentPerResource,
			HideDiffInSpoiler:    hideDiffInSpoiler,
			PrependedCommentText: prependedCommentText,
			AppendedCommentText:  appendedCommentText,
		},
		nil
}

func writeContentToBuffer(buffer *bytes.Buffer, content string, azureDevOpsCommandFlags *AzureDevOpsCommandFlags) error {
	if azureDevOpsCommandFlags.PrependedCommentText != "" {
		_, err := buffer.WriteString(azureDevOpsCommandFlags.PrependedCommentText + "\n")
		if err != nil {
			return errors.Join(errors.New("Adding prepended-comment-text to buffer failed."), err)
		}
	}

	_, err := buffer.WriteString(content)
	if err != nil {
		return errors.Join(errors.New("Adding content to buffer failed."), err)
	}

	if azureDevOpsCommandFlags.AppendedCommentText != "" {
		_, err := buffer.WriteString("\n" + azureDevOpsCommandFlags.AppendedCommentText)
		if err != nil {
			return errors.Join(errors.New("Adding appended-comment-text to buffer failed."), err)
		}
	}

	return nil
}

func wrapContentInSpoiler(content string) string {
	return "<details>\n<summary>Show Diff</summary>\n\n" + content + "\n</details>"
}
