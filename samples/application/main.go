package main

import (
	"context"
	"fmt"
	"os"

	"github.com/liuxiaobopro/dashscope-go"
	"github.com/liuxiaobopro/dashscope-go/app"
)

func main() {
	resp, err := dashscope.Application.Call(context.Background(), &app.ApplicationCallParams{
		APIKey: os.Getenv("DASHSCOPE_API_KEY"),
		AppID:  os.Getenv("DASHSCOPE_APP_ID"),
		Prompt: "你好",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
