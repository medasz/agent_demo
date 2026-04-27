package tool

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/sashabaranov/go-openai"
)

type fakeTool struct {
	name       string
	result     string
	executed   bool
	definition openai.Tool
}

func (f *fakeTool) Name() string {
	return f.name
}

func (f *fakeTool) Definition() openai.Tool {
	return f.definition
}

func (f *fakeTool) Execute(ctx context.Context, argsJSON string) (string, error) {
	f.executed = true
	return f.result, nil
}

func TestRegistryDefinitions(t *testing.T) {
	registry := NewRegistry()
	definitions := registry.Definitions()

	if len(definitions) != 5 {
		t.Fatalf("Definitions() length = %d, want %d", len(definitions), 5)
	}

	names := make(map[string]bool)
	for _, definition := range definitions {
		if definition.Function == nil {
			t.Fatal("Definitions() contains tool without function definition")
		}
		names[definition.Function.Name] = true
	}

	for _, name := range []string{"read_file", "list_file", "edit_file", "create_file", "run_command"} {
		if !names[name] {
			t.Fatalf("Definitions() missing %q", name)
		}
	}
}

func TestRegistryExecutePrefersRegisteredTool(t *testing.T) {
	registry := NewRegistry()
	tool := &fakeTool{
		name:   "read_file",
		result: "from registered tool",
	}
	registry.Register(tool)

	args, err := json.Marshal(map[string]string{"path": "ignored.txt"})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	result, err := registry.Execute(context.Background(), "read_file", string(args))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if result != "from registered tool" {
		t.Fatalf("Execute() result = %q, want %q", result, "from registered tool")
	}
	if !tool.executed {
		t.Fatal("registered tool was not executed")
	}
}

func TestRegistryExecuteUnknownTool(t *testing.T) {
	registry := NewRegistry()

	result, err := registry.Execute(context.Background(), "unknown_tool", `{}`)
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
