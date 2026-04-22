package tool

import (
	"context"
	"encoding/json"

	"github.com/sashabaranov/go-openai"
)

type Registry struct {
	tools []openai.Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: defaultTools(),
	}
}

func (r *Registry) Definitions() []openai.Tool {
	return append([]openai.Tool(nil), r.tools...)
}

func (r *Registry) Execute(ctx context.Context, name, argsJSON string) (string, error) {
	e := Executor{}
	return e.Execute(ctx, name, argsJSON)
}
func defaultTools() []openai.Tool {
	return []openai.Tool{
		{
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
		}, {
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
		}, {
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
		}, {
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
		}, {
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "run_command",
				Description: "Execute system commands.",
				Strict:      false,
				Parameters: json.RawMessage(`{
							"type":"object",
							"properties": {
								"command": {
									"type":"string",
									"description":"The system command to execute (e.g., 'ls', 'pwd', 'cat'). Do not include arguments here."
								},
								"workdir": {
									"type":"string",
									"description":"Optional working directory where the command should be executed (e.g., '/home/user')."
								},
								"args": {
									"type":"array",
									"items": {
										"type":"string",
										"description":"Arguments to pass to the command."		
									},
									"description":"A list of arguments for the command. Each argument should be a separate string (e.g., ['-l', '/tmp'])."
								}
							},
							"required": ["command"]
						}`),
			},
		},
	}
}
