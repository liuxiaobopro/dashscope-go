package main

import (
	"context"
	"fmt"
	"os"

	"github.com/liuxiaobopro/dashscope-go"
	"github.com/liuxiaobopro/dashscope-go/embeddings"
)

func main() {
	resp, err := dashscope.TextEmbedding.Call(context.Background(), &embeddings.TextEmbeddingCallParams{
		APIKey: os.Getenv("DASHSCOPE_API_KEY"),
		Model:  embeddings.TextEmbeddingV4,
		Input:  []string{"风急天高猿啸哀", "渚清沙白鸟飞回"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
