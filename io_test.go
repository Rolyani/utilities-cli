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
