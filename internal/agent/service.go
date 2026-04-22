package agent

import (
	"context"
	"fmt"

	"agent_demo/internal/session"

	"github.com/sashabaranov/go-openai"
)

const maxToolIterations = 8

type LLMClient interface {
	CompleteStream(
		ctx context.Context,
		messages []openai.ChatCompletionMessage,
		tools []openai.Tool,
		onContent func(string),
	) (openai.ChatCompletionMessage, error)
}

type ToolExecutor interface {
	Execute(ctx context.Context, name, argsJSON string) (string, error)
}

type Agent struct {
	llm      LLMClient
	executor ToolExecutor
	session  *session.Session
	tools    []openai.Tool
}

func New(
	client LLMClient,
	executor ToolExecutor,
	session *session.Session,
	tools []openai.Tool,
) *Agent {
	return &Agent{
		llm:      client,
		executor: executor,
		session:  session,
		tools:    append([]openai.Tool(nil), tools...),
	}
}

func (a *Agent) Run(ctx context.Context, input string, onContent func(string)) (string, error) {
	a.session.Append(openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: input,
	})

	for iteration := 0; iteration < maxToolIterations; iteration++ {
		aiResponse, err := a.llm.CompleteStream(ctx, a.session.Messages(), a.tools, onContent)
		if err != nil {
			return "", err
		}

		a.session.Append(aiResponse)
		if len(aiResponse.ToolCalls) == 0 {
			return aiResponse.Content, nil
		}

		for _, toolCall := range aiResponse.ToolCalls {
			result, err := a.executor.Execute(ctx, toolCall.Function.Name, toolCall.Function.Arguments)
			if err != nil {
				result = fmt.Sprintf("tool %s failed: %s", toolCall.Function.Name, err.Error())
			}
			a.session.Append(openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    result,
				Name:       toolCall.Function.Name,
				ToolCallID: toolCall.ID,
			})
		}
	}

	return "", fmt.Errorf("agent exceeded max tool iterations")
}
