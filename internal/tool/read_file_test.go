package tool

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadFileToolExecute(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	args, err := json.Marshal(map[string]string{"path": path})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	tool := &ReadFileTool{}
	result, err := tool.Execute(context.Background(), string(args))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result != "hello" {
		t.Fatalf("Execute() result = %q, want %q", result, "hello")
	}
}

func TestReadFileToolExecuteInvalidArguments(t *testing.T) {
	tool := &ReadFileTool{}

	result, err := tool.Execute(context.Background(), `{"path":`)
	if err == nil {
		t.Fatal("Execute() error = nil, want invalid arguments error")
	}
	if result != "" {
		t.Fatalf("Execute() result = %q, want empty result", result)
	}
	if !strings.Contains(err.Error(), "invalid arguments for read_file") {
		t.Fatalf("Execute() error = %q, want invalid arguments error", err)
	}
}
