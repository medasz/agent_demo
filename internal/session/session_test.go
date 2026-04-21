package session

import (
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestSessionAppendAndMessages(t *testing.T) {
	session := NewSession()

	session.Append(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: "hello",
	})
	session.Append(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: "hi",
	})

	messages := session.Messages()
	if len(messages) != 2 {
		t.Fatalf("Messages() length = %d, want %d", len(messages), 2)
	}
	if messages[0].Role != openai.ChatMessageRoleUser || messages[0].Content != "hello" {
		t.Fatalf("Messages()[0] = %#v, want user hello", messages[0])
	}
	if messages[1].Role != openai.ChatMessageRoleAssistant || messages[1].Content != "hi" {
		t.Fatalf("Messages()[1] = %#v, want assistant hi", messages[1])
	}
}

func TestSessionMessagesReturnsCopy(t *testing.T) {
	session := NewSession()
	session.Append(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: "original",
	})

	messages := session.Messages()
	messages[0] = openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: "changed",
	}

	got := session.Messages()
	if got[0].Role != openai.ChatMessageRoleUser || got[0].Content != "original" {
		t.Fatalf("Session message was mutated through returned slice: %#v", got[0])
	}
}
