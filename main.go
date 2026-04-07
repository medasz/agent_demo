package main

import (
	"bufio"
	"context"
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

		fmt.Println("Agent > 你说的是：", line)
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: line,
		})

		resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
			Model:    "deepseek-chat",
			Messages: messages,
		})
		if err != nil {
			log.Printf("CreateChatCompletion error: %v\n", err)
			continue
		}
		aiResponse := resp.Choices[0].Message
		messages = append(messages, aiResponse)
	}
}
