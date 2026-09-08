package main

import (
	"context"
	"fmt"
	"os"

	"github.com/liuxiaobopro/dashscope-go"
	"github.com/liuxiaobopro/dashscope-go/aigc"
)

func main() {
	resp, err := dashscope.VideoSynthesis.AsyncCall(context.Background(), &aigc.VideoSynthesisCallParams{
		APIKey: os.Getenv("DASHSCOPE_API_KEY"),
		Model:  aigc.Wanx21T2VTurbo,
		Prompt: "一只小猫在草地上奔跑",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
