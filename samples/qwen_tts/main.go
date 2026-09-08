package main

import (
	"context"
	"fmt"
	"os"

	"github.com/liuxiaobopro/dashscope-go/audio/qwen_tts"
)

func main() {
	resp, err := qwentts.SpeechSynthesizer.Call(context.Background(), &qwentts.SpeechSynthesizerCallParams{
		APIKey: os.Getenv("DASHSCOPE_API_KEY"),
		Model:  "qwen-tts",
		Text:   "你好，我是通义千问",
		Voice:  "Cherry",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
