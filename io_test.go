package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFileSuccess(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "input.txt")
	content := []byte("hello\nworld")

	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	got, err := readFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got) != string(content) {
		t.Fatalf("got %q, want %q", string(got), string(content))
	}
}

func TestReadFileDirectoryError(t *testing.T) {
	dir := t.TempDir()

	_, err := readFile(dir)
	if err == nil {
		t.Fatal("expected error for directory path, got nil")
	}
}

func TestWriteNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "output.txt")
	content := []byte("saved content")

	if err := writeNewFile(path, content); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}

	if string(got) != string(content) {
		t.Fatalf("got %q, want %q", string(got), string(content))
	}
}

func TestBackupFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "original.txt")
	content := []byte("original content")

	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	backupPath, err := backupFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}

	if string(got) != string(content) {
		t.Fatalf("got %q, want %q", string(got), string(content))
	}
}

func TestOverwriteFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	if err := os.WriteFile(path, []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := overwriteFile(path, []byte("new content")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected read error: %v", err)
	}

	if string(got) != "new content" {
		t.Fatalf("got %q, want %q", string(got), "new content")
	}
}

func TestEnsureDirCreatesDirectory(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "files")

	if err := ensureDir(target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(target)
	if err != nil {
		t.Fatalf("expected directory to exist: %v", err)
	}

	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", target)
	}
}

func TestListFilesEmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	items, err := listFiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) != 0 {
		t.Fatalf("got %d items, want 0", len(items))
	}
}

func TestListFilesReturnsFilesOnly(t *testing.T) {
	dir := t.TempDir()

	file1 := filepath.Join(dir, "a.txt")
	file2 := filepath.Join(dir, "b.csv")
	subdir := filepath.Join(dir, "nested")

	if err := os.WriteFile(file1, []byte("one"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("two"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatal(err)
	}

	items, err := listFiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
}

func TestListFilesIncludesExpectedPaths(t *testing.T) {
	dir := t.TempDir()

	file1 := filepath.Join(dir, "first.txt")
	file2 := filepath.Join(dir, "second.txt")

	if err := os.WriteFile(file1, []byte("one"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("two"), 0644); err != nil {
		t.Fatal(err)
	}

	items, err := listFiles(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	gotPaths := map[string]bool{}
	for _, item := range items {
		fi, ok := item.(fileItem)
		if !ok {
			t.Fatalf("expected item to be fileItem, got %T", item)
		}
		gotPaths[fi.path] = true
	}

	if !gotPaths[file1] {
		t.Fatalf("missing expected path %q", file1)
	}
	if !gotPaths[file2] {
		t.Fatalf("missing expected path %q", file2)
	}
}
