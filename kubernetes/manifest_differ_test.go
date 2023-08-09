package kubernetes

import (
	"bytes"
	"testing"
)

func TestCreateAndPrintDiffForManifestsReturnsContentIfManifestsAreIdentical(t *testing.T) {
	manifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: v1\nkind: Service\nmetadata:\n  name: backend\n  namespace: my-namespace",
	}

	output := new(bytes.Buffer)
	CreateAndPrintDiffForManifests(manifest, manifest, false, output)

	if output.String() != (manifest.Content + "\n") {
		t.Fatal("Diff should be the content if both manifests are identical. Diff:\n" + output.String())
	}
}

func TestCreateAndPrintDiffForManifestsReturnsCorrectResultIfChangesAreGiven(t *testing.T) {
	oldManifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 1",
	}
	newManifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 2",
	}

	output := new(bytes.Buffer)
	CreateAndPrintDiffForManifests(oldManifest, newManifest, false, output)

	expectedDiff := ` apiVersion: apps/v1
 kind: Deployment
 metadata:
   name: backend
   namespace: my-namespace
 spec:
-  replicas: 1
+  replicas: 2`

	if output.String() != (expectedDiff + "\n") {
		t.Fatal("Diff should be the content if both manifests are identical. Diff:\n" + output.String())
	}
}

func TestCreateAndPrintDiffForManifestsReturnsCorrectResultIfNewManifestIsGiven(t *testing.T) {
	oldManifest := Manifest{}
	newManifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 2",
	}

	output := new(bytes.Buffer)
	CreateAndPrintDiffForManifests(oldManifest, newManifest, false, output)

	expectedDiff := `+apiVersion: apps/v1
+kind: Deployment
+metadata:
+  name: backend
+  namespace: my-namespace
+spec:
+  replicas: 2`

	if output.String() != (expectedDiff + "\n") {
		t.Fatal("Diff should be the content if both manifests are identical. Diff:\n" + output.String())
	}
}

func TestCreateAndPrintDiffForManifestsReturnsCorrectResultIfRemovedManifestIsGiven(t *testing.T) {
	oldManifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 2",
	}
	newManifest := Manifest{}

	output := new(bytes.Buffer)
	CreateAndPrintDiffForManifests(oldManifest, newManifest, false, output)

	expectedDiff := `-apiVersion: apps/v1
-kind: Deployment
-metadata:
-  name: backend
-  namespace: my-namespace
-spec:
-  replicas: 2`

	if output.String() != (expectedDiff + "\n") {
		t.Fatal("Diff should be the content if both manifests are identical. Diff:\n" + output.String())
	}
}

func TestCreateAndPrintDiffForManifestsReturnsCorrectResultIfChangesAreGivenAndMarkdownFormatIsUsed(t *testing.T) {
	oldManifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 1",
	}
	newManifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 2",
	}

	output := new(bytes.Buffer)
	CreateAndPrintDiffForManifests(oldManifest, newManifest, true, output)

	expectedDiff := "```diff\n" + ` apiVersion: apps/v1
 kind: Deployment
 metadata:
   name: backend
   namespace: my-namespace
 spec:
-  replicas: 1
+  replicas: 2` + "\n```"

	if output.String() != (expectedDiff + "\n") {
		t.Fatal("Diff should be the content if both manifests are identical. Diff:\n" + output.String())
	}
}
