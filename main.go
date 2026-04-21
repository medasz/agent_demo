package main

import (
	"agent_demo/internal/config"
	"agent_demo/internal/llm"
	"agent_demo/internal/tool"
	"agent_demo/pkg"
	"bufio"
	"context"
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
	executor := tool.Executor{}
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
					result, err := executor.Execute(context.Background(), toolCall.Function.Name, toolCall.Function.Arguments)
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
