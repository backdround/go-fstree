package fstree_test

import (
	"github.com/lithammer/dedent"
	"os"
	"strings"
)

////////////////////////////////////////////////////////////
// Utility functions

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func prepareYaml(data string) string {
	data = dedent.Dedent(data)
	data = strings.ReplaceAll(data, "\t", "  ")
	return data
}

func createRoot() (rootPath string, clean func()) {
	rootPath, err := os.MkdirTemp("", "go-fstree-test-*.d")
	assertNoError(err)

	clean = func() {
		err := os.RemoveAll(rootPath)
		assertNoError(err)
	}

	return
}
