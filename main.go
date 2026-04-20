package main

import (
	"agent_demo/internal/config"
	"agent_demo/internal/llm"
	"agent_demo/pkg"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	client := llm.NewClient(conf.APIKey, conf.BaseURL, conf.ModelName)
	messages := []openai.ChatCompletionMessage{}
	reader := bufio.NewReader(os.Stdin)

	for {
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
			aiResponse, err := client.CompleteStream(context.Background(), messages, pkg.Tools, func(content string) {
				fmt.Print(content)
			})
			if err != nil {
				log.Printf("CreateChatCompletion error: %v\n", err)
				continue
			}
			fmt.Println()

			messages = append(messages, aiResponse)

			if aiResponse.ToolCalls != nil {
				fmt.Println("Agent > calling tools...")
				for _, toolCall := range aiResponse.ToolCalls {
					var result string
					var err error

					switch toolCall.Function.Name {
					case "read_file":
						var args struct {
							Path string `json:"path"`
						}
						json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
						result, err = pkg.ReadFile(args.Path)
					case "list_file":
						var args struct {
							Path string `json:"path"`
						}
						json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
						result, err = pkg.ListFiles(args.Path)
					case "edit_file":
						var args struct {
							Path string `json:"path"`
							Old  string `json:"old"`
							New  string `json:"new"`
						}
						json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
						err = pkg.EditFile(args.Path, args.Old, args.New)
						if err == nil {
							result = "File edited successfully."
						}

					case "create_file":
						var args struct {
							Path    string `json:"path"`
							Content string `json:"content"`
						}
						json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
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
						json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

						out, errOut, exitCode := pkg.RunCommand(args.Command, args.Workdir, args.Args)
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
				break
			}
		}
	}
}
