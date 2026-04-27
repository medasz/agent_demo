package tool

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type Tool interface {
	Name() string
	Definition() openai.Tool
	Execute(ctx context.Context, argsJSON string) (string, error)
}
