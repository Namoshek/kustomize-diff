package kubernetes

import (
	"fmt"
	"io"
	"strings"

	"github.com/kylelemons/godebug/diff"
)

// Creates and prints the diff for two manifests.
func CreateAndPrintDiffForManifests(old Manifest, new Manifest, formatAsMarkdownCodeBlock bool, output io.Writer) {
	if formatAsMarkdownCodeBlock {
		fmt.Fprintln(output, "```diff")
	}

	if old.Content == new.Content {
		fmt.Fprintln(output, old.Content)

		return
	}

	diff := diff.Diff(old.Content, new.Content)

	// In case a new manifest is given, the old content is diffed as empty line, which we remove.
	diff = strings.TrimPrefix(diff, "-\n")

	// The diff may contain an empty line at the end, which we remove.
	diff = strings.TrimSuffix(diff, "\n ")

	// In case a removed manifest is given, the new content is diffed as empty line, which we remove.
	diff = strings.TrimSuffix(diff, "\n+")

	fmt.Fprintln(output, diff)

	if formatAsMarkdownCodeBlock {
		fmt.Fprintln(output, "```")
	}
}
