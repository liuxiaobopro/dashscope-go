package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/liuxiaobopro/dashscope-go"
	"github.com/liuxiaobopro/dashscope-go/aigc"
	"github.com/liuxiaobopro/dashscope-go/apientities"
)

func main() {
	resp, err := dashscope.Generation.Call(context.Background(), &aigc.GenerationCallParams{
		APIKey: os.Getenv("DASHSCOPE_API_KEY"),
		Model:  "qwen-plus",
		Messages: []apientities.Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Who are you?"},
		},
		ResultFormat: dashscope.Ptr("message"),
	})
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == http.StatusOK {
		fmt.Println(resp.Output.Choices[0].Message.Content)
	} else {
		fmt.Printf("Error: %s - %s\n", resp.Code, resp.Message)
	}
}
