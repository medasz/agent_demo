package tool

import (
	"agent_demo/pkg"
	"context"
	"encoding/json"
	"fmt"
)

type Executor struct {
}

func (e *Executor) Execute(ctx context.Context, name, argsJSON string) (string, error) {
	var result string
	var err error

	switch name {
	case "read_file":
		var args struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for %s: %w", name, err)
		}
		result, err = pkg.ReadFile(args.Path)
	case "list_file":
		var args struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for %s: %w", name, err)
		}
		result, err = pkg.ListFiles(args.Path)
	case "edit_file":
		var args struct {
			Path string `json:"path"`
			Old  string `json:"old"`
			New  string `json:"new"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for %s: %w", name, err)
		}
		err = pkg.EditFile(args.Path, args.Old, args.New)
		if err == nil {
			result = "File edited successfully."
		}

	case "create_file":
		var args struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for %s: %w", name, err)
		}
		err = pkg.CreateFile(args.Path, args.Content)
		if err == nil {
			result = "File created successfully."
		}
	case "run_command":
		var args struct {
			Command string   `json:"command"`
			Args    []string `json:"args"`
			Workdir string   `json:"workdir"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments for %s: %w", name, err)
		}

		out, errOut, exitCode := pkg.RunCommand(args.Command, args.Workdir, args.Args)
		resMap := map[string]interface{}{
			"stdout":    out,
			"stderr":    errOut,
			"exit_code": exitCode,
		}

		b, _ := json.Marshal(resMap)

		result = string(b)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
	return result, err
}
