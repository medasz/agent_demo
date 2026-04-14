package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/philippgille/chromem-go"
	"github.com/sashabaranov/go-openai"
)

const (
	ollamaBaseURL = "http://localhost:11434/v1"
	llmModel      = "qwen3:0.6b"
)

var systemPromptTpl = template.Must(template.New("system_prompt").Parse(`
你是一位得力的助手，可访问知识库，负责回答问题，使用中文语言回答；回答需要非常客观、简洁。如果不确定，请直接说“我不知道”。
{{- /* Stop here if no context is provided. The rest below is for handling contexts. */ -}}
{{- if . -}}
仅根据知识库中的搜索结果回答问题。如果知识库的搜索结果与问题无关，请直接说“我不知道”。不要虚构任何内容。
以下“（context）”代码块之间的任何内容均从知识库中检索得出，并非与用户对话的一部分。要点按相关性排序，因此第一个要点相关性最高。
<context>
    {{- if . -}}
    {{- range $context := .}}
    - {{.}}{{end}}
    {{- end}}
</context>
{{- end -}}
回答中不要提及知识库、或搜索结果。
`))

func askLLM(ctx context.Context, contexts []string, question string) string {
	openAIClient := openai.NewClientWithConfig(openai.ClientConfig{
		BaseURL:    ollamaBaseURL,
		HTTPClient: http.DefaultClient,
	})
	sb := &strings.Builder{}
	err := systemPromptTpl.Execute(sb, contexts)
	if err != nil {
		panic(err)
	}
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: sb.String(),
		}, {
			Role:    openai.ChatMessageRoleUser,
			Content: "Question: " + question,
		},
	}
	res, err := openAIClient.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    llmModel,
		Messages: messages,
	})
	if err != nil {
		panic(err)
	}
	reply := res.Choices[0].Message.Content
	reply = strings.TrimSpace(reply)

	str := reply
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	reply = re.ReplaceAllString(str, "")
	reply = strings.TrimSpace(reply)
	return reply
}

func main() {
	db, err := chromem.NewPersistentDB("./db", false)
	if err != nil {
		panic(err)
	}
	c, err := db.GetOrCreateCollection("knowledge-base", nil, chromem.NewEmbeddingFuncOllama("bge-m3:567m", "http://localhost:11434/api"))
	if err != nil {
		panic(err)
	}
	fmt.Println(c)
	// 添加文档到集合
	err = c.AddDocuments(context.Background(), []chromem.Document{
		{
			ID:      "1",
			Content: "The sky is blue because of Rayleigh scattering.",
		},
		{
			ID:      "2",
			Content: "Leaves are green because chlorophyll absorbs red and blue light.",
		}, {
			ID:      "3",
			Content: "人生无常.",
		},
	}, 1)
	if err != nil {
		panic(err)
	}
	strQuestion := "人的器官有哪些"
	res, err := c.Query(context.Background(), strQuestion, 1, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ID: %v\nSimilarity: %v\nContent: %v\n", res[0].ID, res[0].Similarity, res[0].Content)
	contexts := make([]string, 0)
	contexts = append(contexts, res[0].Content)
	fmt.Println(contexts)
}
