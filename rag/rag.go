package main

import (
	"context"
	"fmt"

	"github.com/philippgille/chromem-go"
)

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
