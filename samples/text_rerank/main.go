package main

import (
	"context"
	"fmt"
	"os"

	"github.com/liuxiaobopro/dashscope-go"
	"github.com/liuxiaobopro/dashscope-go/rerank"
)

func main() {
	resp, err := dashscope.TextReRank.Call(context.Background(), &rerank.TextReRankCallParams{
		APIKey:    os.Getenv("DASHSCOPE_API_KEY"),
		Model:     rerank.GteRerankV2,
		Query:     "什么是文本排序模型",
		Documents: []string{"文本排序模型广泛用于搜索引擎和推荐系统中", "量子计算是计算科学的一个前沿领域"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
