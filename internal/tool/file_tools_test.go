package tool

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListFileToolExecute(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	args, err := json.Marshal(map[string]string{"path": dir})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	tool := &ListFileTool{}
	result, err := tool.Execute(context.Background(), string(args))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(result, "test.txt") {
		t.Fatalf("Execute() result = %q, want it to contain %q", result, "test.txt")
	}
}

func TestCreateFileToolExecute(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "created.txt")

	args, err := json.Marshal(map[string]string{
		"path":    path,
		"content": "hello",
	})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	tool := &CreateFileTool{}
	result, err := tool.Execute(context.Background(), string(args))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result != "File created successfully." {
		t.Fatalf("Execute() result = %q, want %q", result, "File created successfully.")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(content) != "hello" {
		t.Fatalf("file content = %q, want %q", string(content), "hello")
	}
}

func TestCreateFileToolExecuteInvalidArguments(t *testing.T) {
	tool := &CreateFileTool{}

	result, err := tool.Execute(context.Background(), `{"path":`)
	if err == nil {
		t.Fatal("Execute() error = nil, want invalid arguments error")
	}
	if result != "" {
		t.Fatalf("Execute() result = %q, want empty result", result)
	}
	if !strings.Contains(err.Error(), "invalid arguments for create_file") {
		t.Fatalf("Execute() error = %q, want invalid arguments error", err)
	}
}

func TestEditFileToolExecute(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "edit.txt")
	if err := os.WriteFile(path, []byte("hello"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	args, err := json.Marshal(map[string]string{
		"path": path,
		"old":  "hello",
		"new":  "world",
	})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	tool := &EditFileTool{}
	result, err := tool.Execute(context.Background(), string(args))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result != "File edited successfully." {
		t.Fatalf("Execute() result = %q, want %q", result, "File edited successfully.")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(content) != "world" {
		t.Fatalf("file content = %q, want %q", string(content), "world")
	}
}

func TestRunCommandToolExecute(t *testing.T) {
	args, err := json.Marshal(map[string]interface{}{
		"command": "cmd",
		"args":    []string{"/c", "echo", "hello"},
	})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	tool := &RunCommandTool{}
	result, err := tool.Execute(context.Background(), string(args))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(result, `"exit_code":0`) {
		t.Fatalf("Execute() result = %q, want it to contain %q", result, `"exit_code":0`)
	}
}
