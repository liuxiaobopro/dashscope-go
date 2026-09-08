package main

import (
	"context"
	"fmt"
	"os"

	"github.com/liuxiaobopro/dashscope-go"
	"github.com/liuxiaobopro/dashscope-go/aigc"
	"github.com/liuxiaobopro/dashscope-go/apientities"
)

func main() {
	for resp, err := range dashscope.Generation.CallStream(context.Background(), &aigc.GenerationCallParams{
		APIKey: os.Getenv("DASHSCOPE_API_KEY"),
		Model:  "qwen-plus",
		Messages: []apientities.Message{
			{Role: "user", Content: "hello"},
		},
		ResultFormat:      dashscope.Ptr("message"),
		IncrementalOutput: dashscope.Ptr(true),
		Stream:            true,
	}) {
		if err != nil {
			panic(err)
		}
		fmt.Println(resp)
	}
}
