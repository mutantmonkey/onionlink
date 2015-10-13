package main

import (
	"path/filepath"
	"testing"
)

func TestGenerateFilename(t *testing.T) {
	ext1 := filepath.Ext(generateFilename("test.jpg"))
	if ext1 != ".jpg" {
		t.Fatal("Filename for test.jpg did not end with .jpg")
	}

	file2 := generateFilename("test.tar.gz")
	if file2[len(file2)-7:] != ".tar.gz" {
		t.Fatal("Filename for test.tar.gz did not end with .tar.gz")
	}
}
