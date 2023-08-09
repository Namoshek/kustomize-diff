package kubernetes

import "testing"

func TestParsingSingleWordAsManifestFails(t *testing.T) {
	manifest, err := parseManifest("mytwocents")

	if manifest != (Manifest{}) || err == nil {
		t.Fatal("'mytwocents' should not be parsed as manifest successfully.")
	}
}

func TestParsingMultipleWordAsManifestFails(t *testing.T) {
	manifest, err := parseManifest("mytwocents and more")

	if manifest != (Manifest{}) || err == nil {
		t.Fatal("'mytwocents and more' should not be parsed as manifest successfully.")
	}
}

func TestParsingJsonAsManifestFails(t *testing.T) {
	manifest, err := parseManifest("{\"prop\":\"value\"}")

	if manifest != (Manifest{}) || err == nil {
		t.Fatal("JSON should not be parsed as manifest successfully.")
	}
}

func TestParsingNonManifestYamlAsManifestFails(t *testing.T) {
	manifest, err := parseManifest("foo: bar\nbaz: true")

	if manifest != (Manifest{}) || err == nil {
		t.Fatal("Non-manifest YAML should not be parsed as manifest successfully.")
	}
}

func TestParsingIncompleteManifestYamlWithMissingNameAsManifestFails(t *testing.T) {
	manifest, err := parseManifest("kind: Service\napiVersion: v1\nmetadata:\n  namespace: my-namespace")

	if manifest != (Manifest{}) || err == nil {
		t.Fatal("Incomplete manifest YAML should not be parsed as manifest successfully.")
	}
}

func TestParsingIncompleteManifestYamlWithMissingNamespaceAsManifestSucceeds(t *testing.T) {
	manifest, err := parseManifest("kind: Service\napiVersion: v1\nmetadata:\n  name: my-service")

	if manifest == (Manifest{}) || err != nil {
		t.Fatal("Manifest YAML without 'metadata.namespace' should be parsed as manifest successfully.")
	}
}

func TestParsingCompleteManifestYamlAsManifestSucceeds(t *testing.T) {
	manifest, err := parseManifest("kind: Service\napiVersion: v1\nmetadata:\n  name: my-service\n  namespace: my-namespace")

	if manifest == (Manifest{}) || err != nil {
		t.Fatal("Complete manifest YAML should be parsed as manifest successfully.")
	}
}
