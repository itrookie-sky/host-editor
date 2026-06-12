package view

import (
	"host-editor/internal/consts"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateHostFileName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "my-hosts", false},
		{"valid default", "hosts", false},
		{"empty", "", true},
		{"path traversal", "../hosts", true},
		{"slash", "foo/bar", true},
		{"backslash", "foo\\bar", true},
		{"colon", "foo:bar", true},
		{"null byte", "foo\x00bar", true},
		{"dot", ".", true},
		{"double dot", "..", true},
		{"too long", string(make([]byte, 65)), true},
		{"max length", "abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz01", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateHostFileName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHostFileName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestStoreInitCreatesDirAndDefaultFile(t *testing.T) {
	dir := t.TempDir()
	storeDir := filepath.Join(dir, ".host-editor")

	if err := ensureStoreDir(storeDir); err != nil {
		t.Fatalf("ensureStoreDir: %v", err)
	}

	files, err := listHostFiles(storeDir)
	if err != nil {
		t.Fatalf("listHostFiles: %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %d", len(files))
	}

	if _, err := createHostFile(storeDir, consts.DefaultHostFile); err != nil {
		t.Fatalf("createHostFile: %v", err)
	}

	files, err = listHostFiles(storeDir)
	if err != nil {
		t.Fatalf("listHostFiles after create: %v", err)
	}
	if len(files) != 1 || files[0].Name != consts.DefaultHostFile {
		t.Fatalf("expected [hosts], got %v", files)
	}
}

func TestEnsureDefaultHostFileCopiesSystemHosts(t *testing.T) {
	dir := t.TempDir()
	sourcePath := filepath.Join(dir, "system-hosts")
	storeDir := filepath.Join(dir, ".host-editor")
	content := "# system hosts\n127.0.0.1 localhost\n"

	err := os.WriteFile(sourcePath, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("write source hosts: %v", err)
	}

	err = ensureStoreDir(storeDir)
	if err != nil {
		t.Fatalf("ensureStoreDir: %v", err)
	}

	err = ensureDefaultHostFile(storeDir, sourcePath)
	if err != nil {
		t.Fatalf("ensureDefaultHostFile: %v", err)
	}

	got, err := readHostFile(storeDir, consts.DefaultHostFile)
	if err != nil {
		t.Fatalf("readHostFile: %v", err)
	}
	if got != content {
		t.Errorf("default hosts content mismatch:\ngot:  %q\nwant: %q", got, content)
	}
}

func TestReadSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if err := ensureStoreDir(dir); err != nil {
		t.Fatal(err)
	}

	content := "# test hosts\n127.0.0.1 localhost\n"
	if err := saveHostFile(dir, consts.DefaultHostFile, content); err != nil {
		t.Fatalf("saveHostFile: %v", err)
	}

	got, err := readHostFile(dir, consts.DefaultHostFile)
	if err != nil {
		t.Fatalf("readHostFile: %v", err)
	}
	if got != content {
		t.Errorf("content mismatch:\ngot:  %q\nwant: %q", got, content)
	}
}

func TestSaveRejectsInvalidName(t *testing.T) {
	dir := t.TempDir()
	if err := ensureStoreDir(dir); err != nil {
		t.Fatal(err)
	}

	if err := saveHostFile(dir, "../evil", "data"); err == nil {
		t.Error("expected error for path traversal, got nil")
	}
}

func TestCreateRejectsDuplicate(t *testing.T) {
	dir := t.TempDir()
	if err := ensureStoreDir(dir); err != nil {
		t.Fatal(err)
	}

	if _, err := createHostFile(dir, "test"); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := createHostFile(dir, "test"); err == nil {
		t.Error("expected error for duplicate, got nil")
	}
}

func TestDeleteHostFile(t *testing.T) {
	dir := t.TempDir()
	if err := ensureStoreDir(dir); err != nil {
		t.Fatal(err)
	}

	if _, err := createHostFile(dir, "to-delete"); err != nil {
		t.Fatal(err)
	}

	if err := deleteHostFile(dir, "to-delete"); err != nil {
		t.Fatalf("deleteHostFile: %v", err)
	}

	files, _ := listHostFiles(dir)
	for _, f := range files {
		if f.Name == "to-delete" {
			t.Error("file still exists after delete")
		}
	}
}

func TestHostFilePath(t *testing.T) {
	dir := "/tmp/test"

	if got := hostFilePath(dir, "hosts"); got != filepath.Join(dir, "hosts") {
		t.Errorf("default file path: got %s", got)
	}
	if got := hostFilePath(dir, "my-hosts"); got != filepath.Join(dir, "my-hosts.hosts") {
		t.Errorf("custom file path: got %s", got)
	}
}

func TestRealStoreDir(t *testing.T) {
	dir, err := hostStoreDir()
	if err != nil {
		t.Fatalf("hostStoreDir: %v", err)
	}
	if dir == "" {
		t.Error("expected non-empty dir")
	}
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, consts.HostStoreDirName)
	if dir != expected {
		t.Errorf("got %s, want %s", dir, expected)
	}
}
