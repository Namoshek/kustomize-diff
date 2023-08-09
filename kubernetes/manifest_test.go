package kubernetes

import "testing"

func TestCalculateHashForServiceYieldsCorrectResult(t *testing.T) {
	manifest := Manifest{
		ApiVersion: "v1",
		Kind:       "Service",
		Name:       "backend",
		Namespace:  "myapp-dev",
	}
	hash := manifest.CalculateHash()

	if hash != "1716ec9caa244a38d55db16b3d88878c" {
		t.Fatal("The hash for a Service manifest should be calculated correctly.")
	}
}

func TestCalculateHashForDeploymentYieldsCorrectResult(t *testing.T) {
	manifest := Manifest{
		ApiVersion: "apps/v1",
		Kind:       "Deployment",
		Name:       "frontend",
		Namespace:  "foo-bar-baz",
	}
	hash := manifest.CalculateHash()

	if hash != "0ebb1202e9ab7b9cb44b317fb984d6bd" {
		t.Fatal("The hash for a Deployment manifest should be calculated correctly.")
	}
}
