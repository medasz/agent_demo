package session

import "github.com/sashabaranov/go-openai"

type Session struct {
	messages []openai.ChatCompletionMessage
}

func NewSession() *Session {
	return &Session{messages: make([]openai.ChatCompletionMessage, 0)}
}

func (session *Session) Append(message openai.ChatCompletionMessage) {
	session.messages = append(session.messages, message)
}

func (session *Session) Messages() []openai.ChatCompletionMessage {
	return append([]openai.ChatCompletionMessage(nil), session.messages...)
}
