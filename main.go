package main

import (
	"context"
	"log"

	"agent_demo/internal/adapter"
	"agent_demo/internal/adapter/terminal"
	"agent_demo/internal/agent"
	"agent_demo/internal/config"
	"agent_demo/internal/llm"
	"agent_demo/internal/session"
	"agent_demo/internal/tool"
)

func main() {
	// 1. 初始化依赖
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	registry := tool.NewRegistry()
	tools := registry.Definitions()
	client := llm.NewClient(conf.APIKey, conf.BaseURL, conf.ModelName)
	sessioner := session.NewSession()
	agenter := agent.New(client, registry, sessioner, tools)

	// 2. 创建事件总线
	inputBus := make(chan adapter.Message, 10)

	// 3. 启动后台 Agent 调度核心
	go func() {
		for msg := range inputBus {
			// 通知渠道进入思考状态
			msg.ReplyChan <- adapter.ReplyChunk{IsThinking: true}

			// 调用大模型
			_, err := agenter.Run(context.Background(), msg.Content, func(content string) {
				msg.ReplyChan <- adapter.ReplyChunk{
					IsThinking: false,
					Chunk:      content,
				}
			})

			// 结束当前回复
			if err != nil {
				msg.ReplyChan <- adapter.ReplyChunk{Error: err, IsDone: true}
			} else {
				msg.ReplyChan <- adapter.ReplyChunk{IsDone: true}
			}
		}
	}()

	// 4. 在主线程启动终端渠道适配器
	termAdapter := terminal.New()
	if err := termAdapter.Start(inputBus); err != nil {
		log.Fatalf("Terminal error: %v", err)
	}
}
