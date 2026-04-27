package tool

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	registry := &Registry{
		tools: make(map[string]Tool),
	}

	registry.Register(&ReadFileTool{})
	registry.Register(&ListFileTool{})
	registry.Register(&CreateFileTool{})
	registry.Register(&EditFileTool{})
	registry.Register(&RunCommandTool{})

	return registry
}

func (r *Registry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
}

func (r *Registry) Definitions() []openai.Tool {
	definitions := make([]openai.Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		definitions = append(definitions, tool.Definition())
	}
	return definitions
}

func (r *Registry) Execute(ctx context.Context, name, argsJSON string) (string, error) {
	if tool, ok := r.tools[name]; ok {
		return tool.Execute(ctx, argsJSON)
	}
	return "", fmt.Errorf("unknown tool: %s", name)
}
