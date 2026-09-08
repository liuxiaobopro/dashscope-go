package main

import (
	"context"
	"fmt"
	"os"

	"github.com/liuxiaobopro/dashscope-go"
	"github.com/liuxiaobopro/dashscope-go/aigc"
)

func main() {
	resp, err := dashscope.MultiModalConversation.Call(context.Background(), &aigc.MultiModalConversationCallParams{
		APIKey: os.Getenv("DASHSCOPE_API_KEY"),
		Model:  "qwen-vl-plus",
		Messages: []any{
			map[string]any{
				"role": "user",
				"content": []any{
					map[string]any{"image": "https://dashscope.oss-cn-beijing.aliyuncs.com/images/dog_and_girl.jpeg"},
					map[string]any{"text": "这是什么?"},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
