package main

import (
	"context"
	"fmt"
	"os"

	"github.com/liuxiaobopro/dashscope-go/agentstudio"
)

func main() {
	c := agentstudio.NewClient(os.Getenv("DASHSCOPE_API_KEY"), "")
	resp, err := c.AgentsList(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
