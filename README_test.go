package basicutils_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadmeMinimumGoVersionMatchesGoModToolchain(t *testing.T) {
	readme, err := os.ReadFile("README.md")
	require.NoError(t, err)

	mod, err := os.ReadFile("go.mod")
	require.NoError(t, err)

	readmeVersion := regexp.MustCompile(`at least Go ([0-9]+\.[0-9]+\.[0-9]+)`).FindSubmatch(readme)
	require.Len(t, readmeVersion, 2, "README should mention the minimum Go patch version")

	toolchainVersion := regexp.MustCompile(`(?m)^toolchain go([0-9]+\.[0-9]+\.[0-9]+)$`).FindSubmatch(mod)
	require.Len(t, toolchainVersion, 2, "go.mod should declare the Go toolchain version")

	require.Equal(t, string(toolchainVersion[1]), string(readmeVersion[1]))
}
