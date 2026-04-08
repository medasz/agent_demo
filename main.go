package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	// 2. 确保在函数退出时关闭文件
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

func main() {
	config := openai.DefaultConfig(os.Getenv("API_KEY"))
	config.BaseURL = os.Getenv("BASE_URL")
	client := openai.NewClientWithConfig(config)

	messages := []openai.ChatCompletionMessage{}
	reader := bufio.NewReader(os.Stdin)

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

		resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model:    "deepseek-chat",
			Messages: messages,
			Tools: []openai.Tool{
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
						Name:        "crate_file",
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
				},
			},
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

				case "crate_file":
					var args struct {
						Path    string `json:"path"`
						Content string `json:"content"`
					}
					json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
					err = createFile(args.Path, args.Content)
					if err == nil {
						result = "File created successfully."
					}
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

			resp, err = client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
				Model:    "deepseek-chat",
				Messages: messages,
			})

			if err != nil {
				log.Printf("Tool Result ChatCompletion error: %v\n", err)
				continue
			}

			fmt.Println("Agent > ", resp.Choices[0].Message.Content)

			messages = append(messages, resp.Choices[0].Message)
		} else {
			fmt.Println("Agent > ", aiResponse.Content)
		}
	}
}
