package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func listFiles(path string) (string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fmt.Sprintf("%v", fileNames), nil
}

func editFile(path string, old, new string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	newData := strings.ReplaceAll(string(data), old, new)

	return os.WriteFile(path, []byte(newData), 0644)
}

func createFile(path, content string) error {
	// 1. 创建文件
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	// 2. 确保在函数退出时关闭文件
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func runCommand(command, workdir string, args []string) (string, string, int) {
	cmd := exec.Command(command, args...)
	if workdir != "" {
		cmd.Dir = workdir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}
	return stdout.String(), stderr.String(), exitCode
}

func main() {
	config := openai.DefaultConfig(os.Getenv("API_KEY"))
	config.BaseURL = os.Getenv("BASE_URL")
	client := openai.NewClientWithConfig(config)
	modelName := os.Getenv("MODEL_NAME")
	messages := []openai.ChatCompletionMessage{}
	reader := bufio.NewReader(os.Stdin)

	tools := []openai.Tool{
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
	for true {
		fmt.Print("You > ")
		line, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if strings.TrimSpace(line) == "exit" || strings.TrimSpace(line) == "quit" {
			fmt.Println("Bye!")
			break
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: line,
		})
		for {
			resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
				Model:    modelName,
				Messages: messages,
				Tools:    tools,
			})
			if err != nil {
				log.Printf("CreateChatCompletion error: %v\n", err)
				continue
			}
			aiResponse := resp.Choices[0].Message
			messages = append(messages, aiResponse)

			if aiResponse.ToolCalls != nil {
				fmt.Println("Agent > 正在调用工具...")
				for _, toolCall := range aiResponse.ToolCalls {
					var result string
					var err error

					switch toolCall.Function.Name {
					case "read_file":
						var args struct {
							Path string `json:"path"`
						}
						json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
						result, err = readFile(args.Path)
					case "list_file":
						var args struct {
							Path string `json:"path"`
						}
						json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
						result, err = listFiles(args.Path)
					case "edit_file":
						var args struct {
							Path string `json:"path"`
							Old  string `json:"old"`
							New  string `json:"new"`
						}
						json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
						err = editFile(args.Path, args.Old, args.New)
						if err == nil {
							result = "File edited successfully."
						}

					case "create_file":
						var args struct {
							Path    string `json:"path"`
							Content string `json:"content"`
						}
						json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
						err = createFile(args.Path, args.Content)
						if err == nil {
							result = "File created successfully."
						}
					case "run_command":
						var args struct {
							Command string   `json:"command"`
							Args    []string `json:"args"`
							Workdir string   `json:"workdir"`
						}
						json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

						out, errOut, exitCode := runCommand(args.Command, args.Workdir, args.Args)
						resMap := map[string]interface{}{
							"stdout":    out,
							"stderr":    errOut,
							"exit_code": exitCode,
						}

						b, _ := json.Marshal(resMap)

						result = string(b)
					}
					if err != nil {
						result = err.Error()
					}
					messages = append(messages, openai.ChatCompletionMessage{
						Role:       openai.ChatMessageRoleTool,
						Content:    result,
						Name:       toolCall.Function.Name,
						ToolCallID: toolCall.ID,
					})
				}
			} else {
				fmt.Println("Agent > ", aiResponse.Content)
				break
			}
		}
	}
}
