package kubernetes

import (
	"testing"
)

func TestCreateDiffForManifestsReturnsContentIfManifestsAreIdentical(t *testing.T) {
	manifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: v1\nkind: Service\nmetadata:\n  name: backend\n  namespace: my-namespace",
	}

	diff := CreateDiffForManifests(&manifest, &manifest)

	if diff.Diff != (manifest.Content) {
		t.Fatal("Diff should be the content if both manifests are identical. Diff:\n" + diff.Diff)
	}
}

func TestCreateDiffForManifestsReturnsCorrectResultIfChangesAreGiven(t *testing.T) {
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

	diff := CreateDiffForManifests(&oldManifest, &newManifest)

	expectedDiff := ` apiVersion: apps/v1
 kind: Deployment
 metadata:
   name: backend
   namespace: my-namespace
 spec:
-  replicas: 1
+  replicas: 2`

	if diff.Diff != expectedDiff {
		t.Fatal("Diff should show changes if manifest was altered. Diff:\n" + diff.Diff)
	}
}

func TestCreateDiffForManifestsReturnsCorrectResultIfNewManifestIsGiven(t *testing.T) {
	oldManifest := Manifest{}
	newManifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 2",
	}

	diff := CreateDiffForManifests(&oldManifest, &newManifest)

	expectedDiff := `+apiVersion: apps/v1
+kind: Deployment
+metadata:
+  name: backend
+  namespace: my-namespace
+spec:
+  replicas: 2`

	if diff.Diff != expectedDiff {
		t.Fatal("Diff should show all lines as new if manifest was added. Diff:\n" + diff.Diff)
	}
}

func TestCreateDiffForManifestsReturnsCorrectResultIfRemovedManifestIsGiven(t *testing.T) {
	oldManifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "my-namespace",
		Content:    "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: backend\n  namespace: my-namespace\nspec:\n  replicas: 2",
	}
	newManifest := Manifest{}

	diff := CreateDiffForManifests(&oldManifest, &newManifest)

	expectedDiff := `-apiVersion: apps/v1
-kind: Deployment
-metadata:
-  name: backend
-  namespace: my-namespace
-spec:
-  replicas: 2`

	if diff.Diff != expectedDiff {
		t.Fatal("Diff should show all lines as removed if manifest was deleted. Diff:\n" + diff.Diff)
	}
}
