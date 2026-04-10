package main

import (
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
	config := openai.DefaultConfig(os.Getenv("API_KEY"))
	config.BaseURL = os.Getenv("BASE_URL")
	client := openai.NewClientWithConfig(config)
	modelName := os.Getenv("MODEL_NAME")
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
			stream, err := client.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
				Model:    modelName,
				Messages: messages,
				Tools:    pkg.Tools,
			})
			if err != nil {
				log.Printf("CreateChatCompletion error: %v\n", err)
				continue
			}

			var fullContent strings.Builder
			var toolCalls []openai.ToolCall

			//流式输出
			for {
				response, err := stream.Recv()
				if err != nil {
					break
				}

				delta := response.Choices[0].Delta

				// 1. 处理文字内容
				if delta.Content != "" {
					fmt.Print(delta.Content)
					fullContent.WriteString(delta.Content)
				}

				// 2. 处理工具调用
				for _, tcDelta := range delta.ToolCalls {
					if tcDelta.Index == nil {
						continue
					}
					idx := *tcDelta.Index

					for len(toolCalls) <= idx {
						toolCalls = append(toolCalls, openai.ToolCall{
							Index: &idx,
							ID:    tcDelta.ID,
						})
					}

					if tcDelta.Type != "" {
						toolCalls[idx].Type = tcDelta.Type
					}

					if tcDelta.Function.Name != "" {
						toolCalls[idx].Function.Name = tcDelta.Function.Name
					}
					if tcDelta.Function.Arguments != "" {
						toolCalls[idx].Function.Arguments += tcDelta.Function.Arguments
					}
				}
			}
			fmt.Println() // 换行

			aiResponse := openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: fullContent.String(),
			}
			if len(toolCalls) > 0 {
				aiResponse.ToolCalls = toolCalls
			}
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
				//fmt.Println("Agent > ", aiResponse.Content)
				break
			}
		}
	}
}
