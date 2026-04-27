package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type CreateFileTool struct {
}

func (t *CreateFileTool) Name() string {
	return "create_file"
}

func (t *CreateFileTool) Definition() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "create_file",
			Description: "Create a file and overwrite the written content.",
			Strict:      false,
			Parameters: json.RawMessage(`{
				"type":"object",
				"properties": {
					"path": {
						"type":"string",
						"description":"The path to the file to create or overwrite the written"
					},
					"content": {
						"type":"string",
						"description":"The written content"
					}
				},
				"required": ["path","content"]
			}`),
		},
	}
}

func (t *CreateFileTool) Execute(ctx context.Context, argsJSON string) (string, error) {
	var args struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}

	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments for %s: %w", t.Name(), err)
	}

	if err := CreateFile(args.Path, args.Content); err != nil {
		return "", err
	}

	return "File created successfully.", nil
}
