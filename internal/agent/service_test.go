package agent

import (
	"context"
	"testing"

	"agent_demo/internal/session"

	"github.com/sashabaranov/go-openai"
)

type fakeLLM struct {
	responses []openai.ChatCompletionMessage
	calls     int
}

func (f *fakeLLM) CompleteStream(
	ctx context.Context,
	messages []openai.ChatCompletionMessage,
	tools []openai.Tool,
	onContent func(string),
) (openai.ChatCompletionMessage, error) {
	response := f.responses[f.calls]
	f.calls++
	if response.Content != "" {
		onContent(response.Content)
	}
	return response, nil
}

type fakeExecutor struct {
	calls int
}

func (f *fakeExecutor) Execute(ctx context.Context, name, argsJSON string) (string, error) {
	f.calls++
	return "tool result", nil
}

func TestRunWithoutToolCall(t *testing.T) {
	llm := &fakeLLM{
		responses: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "final answer",
			},
		},
	}
	executor := &fakeExecutor{}
	session := session.NewSession()
	agent := New(llm, executor, session, nil)

	var streamed string
	result, err := agent.Run(context.Background(), "hello", func(content string) {
		streamed += content
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result != "final answer" {
		t.Fatalf("Run() result = %q, want %q", result, "final answer")
	}
	if streamed != "final answer" {
		t.Fatalf("streamed content = %q, want %q", streamed, "final answer")
	}
	if executor.calls != 0 {
		t.Fatalf("executor calls = %d, want 0", executor.calls)
	}
	if len(session.Messages()) != 2 {
		t.Fatalf("session messages = %d, want 2", len(session.Messages()))
	}
}

func TestRunWithToolCall(t *testing.T) {
	llm := &fakeLLM{
		responses: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleAssistant,
				ToolCalls: []openai.ToolCall{
					{
						ID:   "call-1",
						Type: openai.ToolTypeFunction,
						Function: openai.FunctionCall{
							Name:      "read_file",
							Arguments: `{"path":"test.txt"}`,
						},
					},
				},
			},
			{
				Role:    openai.ChatMessageRoleAssistant,
				Content: "final answer",
			},
		},
	}
	executor := &fakeExecutor{}
	session := session.NewSession()
	agent := New(llm, executor, session, nil)

	result, err := agent.Run(context.Background(), "read file", func(content string) {})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result != "final answer" {
		t.Fatalf("Run() result = %q, want %q", result, "final answer")
	}
	if llm.calls != 2 {
		t.Fatalf("llm calls = %d, want 2", llm.calls)
	}
	if executor.calls != 1 {
		t.Fatalf("executor calls = %d, want 1", executor.calls)
	}

	messages := session.Messages()
	if len(messages) != 4 {
		t.Fatalf("session messages = %d, want 4", len(messages))
	}
	if messages[2].Role != openai.ChatMessageRoleTool {
		t.Fatalf("messages[2].Role = %q, want tool", messages[2].Role)
	}
	if messages[2].ToolCallID != "call-1" {
		t.Fatalf("messages[2].ToolCallID = %q, want call-1", messages[2].ToolCallID)
	}
}
