package kubernetes

import (
	"bytes"
	"testing"
)

func TestPrintDiffWithoutMarkdownFormatting(t *testing.T) {
	
	manifestDiff := ManifestDiff{
		OldManifest: nil,
		NewManifest: nil,
		Diff: "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 2",
	}

	output := new(bytes.Buffer)
	PrintDiff(&manifestDiff, false, output)

	if output.String() != (manifestDiff.Diff + "\n") {
		t.Fatal("Diff should be the content if both manifests are identical. Diff:\n" + output.String())
	}
}

func TestPrintDiffWithMarkdownFormatting(t *testing.T) {
	
	manifestDiff := ManifestDiff{
		OldManifest: nil,
		NewManifest: nil,
		Diff: "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 2",
	}

	output := new(bytes.Buffer)
	PrintDiff(&manifestDiff, true, output)

	if output.String() != ("```diff\n" + manifestDiff.Diff + "\n```\n") {
		t.Fatal("Diff should be the content if both manifests are identical. Diff:\n" + output.String())
	}
}
