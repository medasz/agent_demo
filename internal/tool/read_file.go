package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type ReadFileTool struct {
}

func (t *ReadFileTool) Name() string {
	return "read_file"
}

func (t *ReadFileTool) Definition() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "read_file",
			Description: "Reads the content of a file.",
			Strict:      false,
			Parameters: json.RawMessage(`{
				"type":"object",
				"properties": {
					"path": {
						"type":"string",
						"description":"The path to the file to read"
					}
				},
				"required": ["path"]
			}`),
		},
	}
}

func (t *ReadFileTool) Execute(ctx context.Context, argsJSON string) (string, error) {
	var args struct {
		Path string `json:"path"`
	}

	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments for %s: %w", t.Name(), err)
	}

	return ReadFile(args.Path)
}
