package fstree_test

import (
	"testing"

	"github.com/backdround/go-fstree/v2"
	"github.com/stretchr/testify/require"
)

func TestMutual(t *testing.T) {
	root, clean := createRoot()
	defer clean()

	yamlData := prepareYaml(`
		new-directory:
			file.txt:
				type: file
			link1:
				type: link
				path: ./file.txt
			subdirectory:
	`)

	err := fstree.MakeOverOSFS(root, yamlData)
	require.NoError(t, err)

	difference, err := fstree.CheckOverOSFS(root, yamlData)
	require.NoError(t, err)
	require.Nil(t, difference)
}
