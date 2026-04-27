package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type ListFileTool struct {
}

func (t *ListFileTool) Name() string {
	return "list_file"
}

func (t *ListFileTool) Definition() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "list_file",
			Description: "List files in a directory.",
			Strict:      false,
			Parameters: json.RawMessage(`{
				"type":"object",
				"properties": {
					"path": {
						"type":"string",
						"description":"The directory path to list files from."
					}
				},
				"required": ["path"]
			}`),
		},
	}
}

func (t *ListFileTool) Execute(ctx context.Context, argsJSON string) (string, error) {
	var args struct {
		Path string `json:"path"`
	}

	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments for %s: %w", t.Name(), err)
	}

	return ListFiles(args.Path)
}
