package main

import (
	"context"
	"fmt"
	"os"

	"github.com/liuxiaobopro/dashscope-go"
	"github.com/liuxiaobopro/dashscope-go/aigc"
)

func main() {
	resp, err := dashscope.ImageSynthesis.AsyncCall(context.Background(), &aigc.ImageSynthesisCallParams{
		APIKey: os.Getenv("DASHSCOPE_API_KEY"),
		Model:  aigc.WanxV1,
		Prompt: "一只在草地上奔跑的小狗",
		N:      dashscope.Ptr(1),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	waited, err := dashscope.ImageSynthesis.Wait(context.Background(), resp, os.Getenv("DASHSCOPE_API_KEY"), "", -1)
	if err != nil {
		panic(err)
	}
	fmt.Println(waited)
}
