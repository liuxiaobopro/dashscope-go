# DashScope Go SDK

> **English** | [中文](README_zh.md)

The DashScope Go SDK provides a comprehensive interface to [Alibaba Cloud Model Studio (Bailian)](https://www.alibabacloud.com/help/en/model-studio/) APIs, covering text generation, multi-modal understanding, embeddings, reranking, image/video generation, speech synthesis & recognition, and more.

This repository is a 1:1 Go port of [dashscope-sdk-python](https://github.com/dashscope/dashscope-sdk-python) v1.27.4 (without CLI / acli).

## Installation

```shell
go get github.com/liuxiaobopro/dashscope-go
```

Requires Go 1.23+.

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/liuxiaobopro/dashscope-go"
	"github.com/liuxiaobopro/dashscope-go/aigc"
	"github.com/liuxiaobopro/dashscope-go/apientities"
)

func main() {
	resp, err := dashscope.Generation.Call(context.Background(), &aigc.GenerationCallParams{
		Model: "qwen-plus",
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
```

## API Key Authentication

The SDK uses API key for authentication. Please refer to [official documentation for alibabacloud china](https://www.alibabacloud.com/help/en/model-studio/) and [official documentation for alibabacloud international](https://www.alibabacloud.com/help/en/model-studio/) regarding how to obtain your api-key.

### Using the API Key

1. Set the API key via code

```go
dashscope.SetAPIKey("YOUR-DASHSCOPE-API-KEY")
// Or specify the API key file path via code
// dashscope.SetAPIKeyFilePath("~/.dashscope/api_key")
```

2. Set the API key via environment variables

a. Set the API key directly using the environment variable below

```shell
export DASHSCOPE_API_KEY='YOUR-DASHSCOPE-API-KEY'
```

b. Specify the API key file path via an environment variable

```shell
export DASHSCOPE_API_KEY_FILE_PATH='~/.dashscope/api_key'
```

3. Save the API key to a file

```go
_ = dashscope.SaveAPIKey("YOUR-DASHSCOPE-API-KEY", "")
// empty path saves to default location "~/.dashscope/api_key"
```

## Supported Models

| Category | Recommended Models | SDK Type |
|----------|-------------------|-----------|
| Text Generation | qwen3.8-max, qwen3.7-max, qwen3.7-plus, qwen3.6-flash | `Generation` |
| Multi-Modal Understanding | qwen3.5-omni-plus, qwen3.7-plus (vision) | `MultiModalConversation` |
| Text Embedding | text-embedding-v4, text-embedding-v3 | `TextEmbedding` |
| Multi-Modal Embedding | tongyi-embedding-vision-plus, qwen3-vl-embedding | `MultiModalEmbedding` |
| Text ReRank | qwen3-rerank, gte-rerank-v2 | `TextReRank` |
| Image Generation | wan2.7-image-pro, qwen-image-2.0-pro | `ImageSynthesis` |
| Video Generation | wan2.7-t2v, wan2.7-i2v, happyhorse-1.0-t2v/i2v | `VideoSynthesis` |
| Speech Synthesis (TTS) | cosyvoice-v3.5-plus, cosyvoice-v1 | `SpeechSynthesizer`, `HttpSpeechSynthesizer` |
| Speech Recognition (ASR) | fun-asr-realtime, fun-asr, paraformer-v1 | `Transcription` |
| Omni (Real-time) | qwen3.5-omni-plus-realtime | `MultiModalConversation` |

For the latest model list, visit [Bailian Model Plaza](https://bailian.console.aliyun.com/).

## Logging

To output Dashscope logs, you need to configure the logger.

```shell
export DASHSCOPE_LOGGING_LEVEL='info'
```

## Output

The output contains the following fields:

```
     request_id (str): The request id.
     status_code (int): HTTP status code, 200 indicates that the
         request was successful, others indicate an error.
     code (str): Error code if error occurs, otherwise empty str.
     message (str): Set to error message on error.
     output (Any): The request output.
     usage (Any): The request usage information.
```

## License

This project is licensed under the Apache License (Version 2.0).
