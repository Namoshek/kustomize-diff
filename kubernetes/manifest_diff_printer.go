package kubernetes

import (
	"fmt"
	"io"
)

// Creates and prints the diff for two manifests.
func PrintDiff(diff *ManifestDiff, formatAsMarkdownCodeBlock bool, output io.Writer) {
	if formatAsMarkdownCodeBlock {
		fmt.Fprintln(output, "```diff")
	}

	fmt.Fprintln(output, diff.Diff)

	if formatAsMarkdownCodeBlock {
		fmt.Fprintln(output, "```")
	}
}
