package llm

import (
	"context"
	"io"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	client    *openai.Client
	modelName string
}

func NewClient(apiKey, baseURL, modelName string) *Client {
	conf := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		conf.BaseURL = baseURL
	}

	return &Client{
		client:    openai.NewClientWithConfig(conf),
		modelName: modelName,
	}
}

func (c *Client) CompleteStream(
	ctx context.Context,
	messages []openai.ChatCompletionMessage,
	tools []openai.Tool,
	onContent func(string),
) (openai.ChatCompletionMessage, error) {
	stream, err := c.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:    c.modelName,
		Messages: messages,
		Tools:    tools,
	})
	if err != nil {
		return openai.ChatCompletionMessage{}, err
	}

	var fullContent strings.Builder
	var reasoningContent strings.Builder
	var toolCalls []openai.ToolCall

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return openai.ChatCompletionMessage{}, err
		}
		if len(response.Choices) == 0 {
			continue
		}

		delta := response.Choices[0].Delta
		if delta.Content != "" {
			onContent(delta.Content)
			fullContent.WriteString(delta.Content)
		}
		if delta.ReasoningContent != "" {
			reasoningContent.WriteString(delta.ReasoningContent)
		}

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

	message := openai.ChatCompletionMessage{
		Role:             openai.ChatMessageRoleAssistant,
		Content:          fullContent.String(),
		ReasoningContent: reasoningContent.String(),
	}
	if len(toolCalls) > 0 {
		message.ToolCalls = toolCalls
	}

	return message, nil
}
