package main

import (
	"agent_demo/internal/agent"
	"agent_demo/internal/config"
	"agent_demo/internal/llm"
	"agent_demo/internal/session"
	"agent_demo/internal/tool"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	registry := tool.NewRegistry()
	tools := registry.Definitions()
	executor := tool.Executor{}
	client := llm.NewClient(conf.APIKey, conf.BaseURL, conf.ModelName)
	sessioner := session.NewSession()
	reader := bufio.NewReader(os.Stdin)
	agenter := agent.New(client, &executor, sessioner, tools)
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
		_, err = agenter.Run(context.Background(), line, func(content string) {
			fmt.Print(content)
		})
		if err != nil {
			log.Printf("agent run: %v\n", err)
			continue
		}
		fmt.Println()
	}
}
