package tool

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecutorReadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	args, err := json.Marshal(map[string]string{"path": path})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	executor := Executor{}
	result, err := executor.Execute(context.Background(), "read_file", string(args))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result != "hello" {
		t.Fatalf("Execute() result = %q, want %q", result, "hello")
	}
}

func TestExecutorUnknownTool(t *testing.T) {
	executor := Executor{}
	result, err := executor.Execute(context.Background(), "unknown_tool", `{}`)
	if err == nil {
		t.Fatal("Execute() error = nil, want unknown tool error")
	}
	if result != "" {
		t.Fatalf("Execute() result = %q, want empty result", result)
	}
	if !strings.Contains(err.Error(), "unknown tool") {
		t.Fatalf("Execute() error = %q, want unknown tool error", err)
	}
}
