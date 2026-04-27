package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type EditFileTool struct {
}

func (t *EditFileTool) Name() string {
	return "edit_file"
}

func (t *EditFileTool) Definition() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "edit_file",
			Description: "Edits a file by replacing a string.",
			Strict:      false,
			Parameters: json.RawMessage(`{
				"type":"object",
				"properties": {
					"path": {
						"type":"string",
						"description":"The path to the file to edit"
					},
					"old": {
						"type":"string",
						"description":"The string to be replaced"
					},
					"new": {
						"type":"string",
						"description":"The new string to replace with"
					}
				},
				"required": ["path","old","new"]
			}`),
		},
	}
}

func (t *EditFileTool) Execute(ctx context.Context, argsJSON string) (string, error) {
	var args struct {
		Path string `json:"path"`
		Old  string `json:"old"`
		New  string `json:"new"`
	}

	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments for %s: %w", t.Name(), err)
	}

	if err := EditFile(args.Path, args.Old, args.New); err != nil {
		return "", err
	}

	return "File edited successfully.", nil
}
