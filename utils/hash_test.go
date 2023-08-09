package utils

import "testing"

func TestCalculateMD5AsStringProducesCorrectOutput1(t *testing.T) {
	hash := CalculateMD5AsString("foobarbaz")

	if hash != "6df23dc03f9b54cc38a0fc1483df6e21" {
		t.Fatal("The hash for 'foobarbaz' should be '6df23dc03f9b54cc38a0fc1483df6e21'.")
	}
}

func TestCalculateMD5AsStringProducesCorrectOutput2(t *testing.T) {
	hash := CalculateMD5AsString("Hello, World!")

	if hash != "65a8e27d8879283831b664bd8b7f0ad4" {
		t.Fatal("The hash for 'Hello, World!' should be '65a8e27d8879283831b664bd8b7f0ad4'.")
	}
}

func TestCalculateMD5AsStringProducesCorrectOutput3(t *testing.T) {
	hash := CalculateMD5AsString("This is just some example sentence.")

	if hash != "9990057bcb3fb87187126cc6045df6f3" {
		t.Fatal("The hash for 'This is just some example sentence.' should be '9990057bcb3fb87187126cc6045df6f3'.")
	}
}
