# DashScope Go SDK

> [English](README.md) | **中文**

DashScope Go SDK 提供了访问[阿里云百炼（Model Studio）](https://help.aliyun.com/zh/model-studio/) API 的完整接口，覆盖文本生成、多模态理解、向量（Embedding）、重排（Rerank）、图像/视频生成、语音合成与识别等能力。

本仓库是 [dashscope-sdk-python](https://github.com/dashscope/dashscope-sdk-python) v1.27.4 的 Go 一比一移植（不含 CLI / acli）。

## 安装

```shell
go get github.com/liuxiaobopro/dashscope-go
```

需要 Go 1.23+。

## 快速开始

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

## API Key 鉴权

SDK 使用 API Key 进行鉴权。获取 API Key 的方法请参考[阿里云百炼官方文档（国内站）](https://help.aliyun.com/zh/model-studio/)和[阿里云百炼官方文档（国际站）](https://www.alibabacloud.com/help/en/model-studio/)。

### 使用 API Key

1. 通过代码设置 API Key

```go
dashscope.SetAPIKey("YOUR-DASHSCOPE-API-KEY")
// 或者通过代码指定 API Key 文件路径
// dashscope.SetAPIKeyFilePath("~/.dashscope/api_key")
```

2. 通过环境变量设置 API Key

a. 直接使用以下环境变量设置 API Key

```shell
export DASHSCOPE_API_KEY='YOUR-DASHSCOPE-API-KEY'
```

b. 通过环境变量指定 API Key 文件路径

```shell
export DASHSCOPE_API_KEY_FILE_PATH='~/.dashscope/api_key'
```

3. 将 API Key 保存到文件

```go
_ = dashscope.SaveAPIKey("YOUR-DASHSCOPE-API-KEY", "")
// 空路径保存到默认位置 "~/.dashscope/api_key"
```

## 支持的模型

| 类别 | 推荐模型 | SDK 类型 |
|----------|-------------------|-----------|
| 文本生成 | qwen3.8-max、qwen3.7-max、qwen3.7-plus、qwen3.6-flash | `Generation` |
| 多模态理解 | qwen3.5-omni-plus、qwen3.7-plus（视觉） | `MultiModalConversation` |
| 文本向量 | text-embedding-v4、text-embedding-v3 | `TextEmbedding` |
| 多模态向量 | tongyi-embedding-vision-plus、qwen3-vl-embedding | `MultiModalEmbedding` |
| 文本重排 | qwen3-rerank、gte-rerank-v2 | `TextReRank` |
| 图像生成 | wan2.7-image-pro、qwen-image-2.0-pro | `ImageSynthesis` |
| 视频生成 | wan2.7-t2v、wan2.7-i2v、happyhorse-1.0-t2v/i2v | `VideoSynthesis` |
| 语音合成（TTS） | cosyvoice-v3.5-plus、cosyvoice-v1 | `SpeechSynthesizer`、`HttpSpeechSynthesizer` |
| 语音识别（ASR） | fun-asr-realtime、fun-asr、paraformer-v1 | `Transcription` |
| 全模态（实时） | qwen3.5-omni-plus-realtime | `MultiModalConversation` |

最新模型列表请访问[百炼模型广场](https://bailian.console.aliyun.com/)。

## 日志

如需输出 DashScope 日志，请配置日志级别：

```shell
export DASHSCOPE_LOGGING_LEVEL='info'
```

## 输出

输出包含以下字段：

```
     request_id (str): 请求 ID。
     status_code (int): HTTP 状态码，200 表示请求成功，其他值表示错误。
     code (str): 出错时的错误码，否则为空字符串。
     message (str): 出错时设置为错误信息。
     output (Any): 请求输出。
     usage (Any): 请求用量信息。
```

## 许可证

本项目采用 Apache License (Version 2.0) 许可证。
