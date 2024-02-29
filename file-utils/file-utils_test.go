/*
 * @kordax (Dmitry Morozov)
 * dmorozov@valoru-software.com
 * Copyright (c) 2024.
 */

package fileutils_test

import (
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	fileutils "github.com/kordax/basic-utils/file-utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExists(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // clean up

	exists := fileutils.Exists(tmpfile.Name())
	assert.True(t, exists, "File should exist")

	exists = fileutils.Exists("non_existent_file.xyz")
	assert.False(t, exists, "File should not exist")
}

func TestMustRead(t *testing.T) {
	content := "Hello, World!"
	tmpfile, err := ioutil.TempFile("", "example")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // clean up

	_, err = tmpfile.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, tmpfile.Close())

	readContent := fileutils.MustRead(tmpfile.Name())
	assert.Equal(t, content, string(readContent), "Content should match")
}

func TestCreateFile(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	require.NoError(t, err)
	tmpfilePath := tmpfile.Name()
	require.NoError(t, tmpfile.Close())
	defer os.Remove(tmpfilePath) // clean up

	// Test: Create a file with content
	testContent := "Test Content"
	err = fileutils.CreateFile(tmpfilePath, &testContent)
	require.NoError(t, err, "File should be created without error")

	// Verify content
	readContent, err := ioutil.ReadFile(tmpfilePath)
	require.NoError(t, err, "Should read file without error")
	assert.Equal(t, testContent, string(readContent), "File content should match")
}

func TestListFiles(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "exampledir")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir) // clean up

	expectedFiles := []string{"file1.txt", "file2.log", "file3.txt"}
	for _, fname := range expectedFiles {
		tmpfilePath := tmpdir + "/" + fname
		err := ioutil.WriteFile(tmpfilePath, []byte("content"), 0644)
		require.NoError(t, err, "Should create file without error")
	}

	files, err := fileutils.ListFiles(tmpdir, nil)
	require.NoError(t, err, "Should list files without error")
	assert.Len(t, files, 3, "Should list all files")

	regex := regexp.MustCompile(`\.txt$`).String() // Filter for .txt files
	files, err = fileutils.ListFiles(tmpdir, &regex)
	require.NoError(t, err, "Should list files without error")
	assert.Len(t, files, 2, "Should list only .txt files")
}
