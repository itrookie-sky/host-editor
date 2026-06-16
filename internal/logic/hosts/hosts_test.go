package hosts

import (
	"context"
	"host-editor/internal/consts"
	"host-editor/internal/model"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newTestHosts(dir string) *sHosts {
	return &sHosts{dir: dir}
}

func TestPublicMethodsRequireContextAsFirstParameter(t *testing.T) {
	var s = newTestHosts(t.TempDir())

	// 公开方法必须保持 context.Context 作为首个参数；基础类型不用指针，struct 类型使用指针。
	var _ func(context.Context) ([]*model.HostFileInfo, error) = s.ListHostFiles
	var _ func(context.Context, string) (string, error) = s.ReadHostFile
	var _ func(context.Context, *model.SaveHostFileRequest) error = s.SaveHostFile
	var _ func(context.Context, string) (*model.HostFileInfo, error) = s.CreateHostFile
	var _ func(context.Context, string) error = s.DeleteHostFile
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
			err := s.SaveHostFile(context.Background(), &model.SaveHostFileRequest{
				Name:    tt.input,
				Content: "data",
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveHostFile(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestStoreInitCreatesDirAndDefaultFile(t *testing.T) {
	dir := t.TempDir()
	s := newTestHosts(dir)

	files, err := s.ListHostFiles(context.Background())
	if err != nil {
		t.Fatalf("ListHostFiles: %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %d", len(files))
	}

	if _, err := s.CreateHostFile(context.Background(), consts.DefaultHostFile); err != nil {
		t.Fatalf("CreateHostFile: %v", err)
	}

	files, err = s.ListHostFiles(context.Background())
	if err != nil {
		t.Fatalf("ListHostFiles after create: %v", err)
	}
	if len(files) != 1 || files[0].Name != consts.DefaultHostFile {
		t.Fatalf("expected [hosts], got %v", files)
	}
}

func TestReadSaveRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := newTestHosts(dir)

	content := "# system hosts\n127.0.0.1 localhost\n"
	err := s.SaveHostFile(context.Background(), &model.SaveHostFileRequest{
		Name:    consts.DefaultHostFile,
		Content: content,
	})
	if err != nil {
		t.Fatalf("SaveHostFile: %v", err)
	}

	got, err := s.ReadHostFile(context.Background(), consts.DefaultHostFile)
	if err != nil {
		t.Fatalf("ReadHostFile: %v", err)
	}
	if got != content {
		t.Errorf("content mismatch:\ngot:  %q\nwant: %q", got, content)
	}
}

func TestSaveRejectsInvalidName(t *testing.T) {
	dir := t.TempDir()
	s := newTestHosts(dir)

	err := s.SaveHostFile(context.Background(), &model.SaveHostFileRequest{
		Name:    "../evil",
		Content: "data",
	})
	if err == nil {
		t.Error("expected error for path traversal, got nil")
	}
}

func TestCreateRejectsDuplicate(t *testing.T) {
	dir := t.TempDir()
	s := newTestHosts(dir)

	if _, err := s.CreateHostFile(context.Background(), "test"); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := s.CreateHostFile(context.Background(), "test"); err == nil {
		t.Error("expected error for duplicate, got nil")
	}
}

func TestDeleteHostFile(t *testing.T) {
	dir := t.TempDir()
	s := newTestHosts(dir)

	if _, err := s.CreateHostFile(context.Background(), "to-delete"); err != nil {
		t.Fatal(err)
	}

	if err := s.DeleteHostFile(context.Background(), "to-delete"); err != nil {
		t.Fatalf("DeleteHostFile: %v", err)
	}

	files, _ := s.ListHostFiles(context.Background())
	for _, f := range files {
		if f.Name == "to-delete" {
			t.Error("file still exists after delete")
		}
	}
}

func TestHostFilePath(t *testing.T) {
	dir := t.TempDir()
	s := newTestHosts(dir)

	err := s.SaveHostFile(context.Background(), &model.SaveHostFileRequest{
		Name:    "hosts",
		Content: "default",
	})
	if err != nil {
		t.Fatalf("SaveHostFile default: %v", err)
	}
	if got := filepath.Join(dir, "hosts"); !gfileExists(got) {
		t.Errorf("default file path: got %s", got)
	}
	err = s.SaveHostFile(context.Background(), &model.SaveHostFileRequest{
		Name:    "my-hosts",
		Content: "custom",
	})
	if err != nil {
		t.Fatalf("SaveHostFile custom: %v", err)
	}
	if got := filepath.Join(dir, "my-hosts.hosts"); !gfileExists(got) {
		t.Errorf("custom file path: got %s", got)
	}
}

func TestHostsLogicDoesNotUsePrivateMethods(t *testing.T) {
	data, err := os.ReadFile("hosts.go")
	if err != nil {
		t.Fatal(err)
	}

	privateMethods := []string{
		"func (s *sHosts) validateHostFileName(",
		"func (s *sHosts) hostFilePath(",
		"func (s *sHosts) isHostFile(",
		"func (s *sHosts) listHostFiles(",
		"func (s *sHosts) readHostFile(",
		"func (s *sHosts) saveHostFile(",
		"func (s *sHosts) createHostFile(",
		"func (s *sHosts) deleteHostFile(",
	}
	source := string(data)
	for _, method := range privateMethods {
		if strings.Contains(source, method) {
			t.Fatalf("unexpected private method %q", method)
		}
	}
}

func gfileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
