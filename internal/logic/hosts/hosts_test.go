package hosts

import (
	"host-editor/internal/consts"
	"os"
	"path/filepath"
	"testing"
)

func newTestHosts(dir string) *sHosts {
	return &sHosts{dir: dir}
}

func TestValidateHostFileName(t *testing.T) {
	s := newTestHosts(t.TempDir())
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
			err := s.validateHostFileName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateHostFileName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestStoreInitCreatesDirAndDefaultFile(t *testing.T) {
	dir := t.TempDir()
	s := newTestHosts(dir)

	files, err := s.listHostFiles()
	if err != nil {
		t.Fatalf("listHostFiles: %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %d", len(files))
	}

	if _, err := s.createHostFile(consts.DefaultHostFile); err != nil {
		t.Fatalf("createHostFile: %v", err)
	}

	files, err = s.listHostFiles()
	if err != nil {
		t.Fatalf("listHostFiles after create: %v", err)
	}
	if len(files) != 1 || files[0].Name != consts.DefaultHostFile {
		t.Fatalf("expected [hosts], got %v", files)
	}
}

func TestReadSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := newTestHosts(dir)

	content := "# system hosts\n127.0.0.1 localhost\n"
	if err := s.saveHostFile(consts.DefaultHostFile, content); err != nil {
		t.Fatalf("saveHostFile: %v", err)
	}

	got, err := s.readHostFile(consts.DefaultHostFile)
	if err != nil {
		t.Fatalf("readHostFile: %v", err)
	}
	if got != content {
		t.Errorf("content mismatch:\ngot:  %q\nwant: %q", got, content)
	}
}

func TestSaveRejectsInvalidName(t *testing.T) {
	dir := t.TempDir()
	s := newTestHosts(dir)

	if err := s.saveHostFile("../evil", "data"); err == nil {
		t.Error("expected error for path traversal, got nil")
	}
}

func TestCreateRejectsDuplicate(t *testing.T) {
	dir := t.TempDir()
	s := newTestHosts(dir)

	if _, err := s.createHostFile("test"); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := s.createHostFile("test"); err == nil {
		t.Error("expected error for duplicate, got nil")
	}
}

func TestDeleteHostFile(t *testing.T) {
	dir := t.TempDir()
	s := newTestHosts(dir)

	if _, err := s.createHostFile("to-delete"); err != nil {
		t.Fatal(err)
	}

	if err := s.deleteHostFile("to-delete"); err != nil {
		t.Fatalf("deleteHostFile: %v", err)
	}

	files, _ := s.listHostFiles()
	for _, f := range files {
		if f.Name == "to-delete" {
			t.Error("file still exists after delete")
		}
	}
}

func TestHostFilePath(t *testing.T) {
	dir := t.TempDir()
	s := newTestHosts(dir)

	if got := s.hostFilePath("hosts"); got != filepath.Join(dir, "hosts") {
		t.Errorf("default file path: got %s", got)
	}
	if got := s.hostFilePath("my-hosts"); got != filepath.Join(dir, "my-hosts.hosts") {
		t.Errorf("custom file path: got %s", got)
	}
}

func TestInitPathComputation(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	expected := filepath.Join(home, consts.HostStoreDirName)
	if expected == "" {
		t.Error("expected non-empty dir")
	}
}
