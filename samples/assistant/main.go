package main

import (
	"context"
	"fmt"
	"os"

	"github.com/liuxiaobopro/dashscope-go"
)

func main() {
	resp, err := dashscope.Assistants.Create(context.Background(), map[string]any{
		"model":        "qwen-plus",
		"name":         "helper",
		"instructions": "You are a helpful assistant.",
	}, os.Getenv("DASHSCOPE_API_KEY"), "")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
