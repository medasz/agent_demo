package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type RunCommandTool struct {
}

func (t *RunCommandTool) Name() string {
	return "run_command"
}

func (t *RunCommandTool) Definition() openai.Tool {
	return openai.Tool{
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
	}
}

func (t *RunCommandTool) Execute(ctx context.Context, argsJSON string) (string, error) {
	var args struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
		Workdir string   `json:"workdir"`
	}

	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid arguments for %s: %w", t.Name(), err)
	}

	out, errOut, exitCode := RunCommand(args.Command, args.Workdir, args.Args)
	resMap := map[string]interface{}{
		"stdout":    out,
		"stderr":    errOut,
		"exit_code": exitCode,
	}

	result, err := json.Marshal(resMap)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
