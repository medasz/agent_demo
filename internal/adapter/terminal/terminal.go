package terminal

import (
	"agent_demo/internal/adapter"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type TerminalChannel struct{}

func New() *TerminalChannel {
	return &TerminalChannel{}
}

func (t *TerminalChannel) Start(inputBus chan<- adapter.Message) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You > ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		
		text := strings.TrimSpace(line)
		if text == "exit" || text == "quit" {
			fmt.Println("Bye!")
			break
		}
		if text == "" {
			continue
		}

		replyChan := make(chan adapter.ReplyChunk)
		inputBus <- adapter.Message{
			SessionID: "terminal",
			Content:   text,
			ReplyChan: replyChan,
		}

		// 等待并打印回复
		for reply := range replyChan {
			if reply.IsThinking {
				// 终端下暂时忽略思考状态
				continue
			}
			if reply.Error != nil {
				fmt.Printf("\nError: %v\n", reply.Error)
				break
			}
			if reply.IsDone {
				fmt.Println()
				break
			}
			fmt.Print(reply.Chunk)
		}
	}
	return nil
}
