package azuredevops

import (
	"context"
	"errors"
	"strings"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/namoshek/kustomize-diff/utils"

	"go.uber.org/zap"
)

// Common parameters required to work with the Azure DevOps API.
type AzureDevOpsParameters struct {
	Instance            string
	Organization        string
	PersonalAccessToken string
	Project             string
	PullRequestId       int
	RepositoryId        string
}

// Retrieves the organization base URL from the parameters based on the instance and organization.
func (a AzureDevOpsParameters) GetOrganizationUri() string {
	return strings.TrimSuffix(a.Instance, "/") + "/" + a.Organization
}

// Creates a comment with the given content on a pull request specified by the given parameters.
func CreatePullRequestComment(azureDevOpsParameters *AzureDevOpsParameters, content string) (error) {
	// Create a connection to your organization or collection.
	connection := azuredevops.NewPatConnection(azureDevOpsParameters.GetOrganizationUri(), azureDevOpsParameters.PersonalAccessToken)

	ctx := context.Background()

	// Create a HTTP client to interact with the core API of Azure DevOps.
	httpClient, err := git.NewClient(ctx, connection)
	if err != nil {
		return errors.Join(errors.New("Creating HTTP client for Azure DevOps 'git' section failed."), err)
	}

	// Create the comment.
	createThreadRequest := git.CreateThreadArgs{
		CommentThread: &git.GitPullRequestCommentThread{
			Comments: &[]git.Comment{
				{
					CommentType: &git.CommentTypeValues.CodeChange,
					Content:     &content,
				},
			},
			Status: &git.CommentThreadStatusValues.Active,
		},
		Project:       &azureDevOpsParameters.Project,
		PullRequestId: &azureDevOpsParameters.PullRequestId,
		RepositoryId:  &azureDevOpsParameters.RepositoryId,
	}

	responseValue, err := httpClient.CreateThread(ctx, createThreadRequest)
	if err != nil {
		errors.Join(errors.New("Creating new comment thread on pull request failed."), err)
	}

	utils.Logger.Info("Created comment thread successfully.", zap.Int("threadId", *responseValue.Id))
	utils.Logger.Debug("The created comment thread:", zap.Any("createdThread", responseValue))

	return nil
}
