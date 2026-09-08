package main

import (
	"context"
	"fmt"
	"os"

	"github.com/liuxiaobopro/dashscope-go"
	"github.com/liuxiaobopro/dashscope-go/embeddings"
)

func main() {
	resp, err := dashscope.MultiModalEmbedding.Call(context.Background(), &embeddings.MultiModalEmbeddingCallParams{
		APIKey: os.Getenv("DASHSCOPE_API_KEY"),
		Model:  embeddings.TongyiEmbeddingVisionPlus,
		Input: []embeddings.MultiModalEmbeddingItemBase{
			embeddings.NewMultiModalEmbeddingItemText("hello", 1),
			embeddings.NewMultiModalEmbeddingItemImage("https://dashscope.oss-cn-beijing.aliyuncs.com/images/dog_and_girl.jpeg", 1),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
