# DashScope Go SDK API

Godoc: `go doc github.com/liuxiaobopro/dashscope-go`

## Packages

| Package | Description |
|---------|-------------|
| `dashscope` | Root package: API key, Generation / Embedding / ReRank re-exports |
| `aigc` | Generation, Conversation, ImageSynthesis, VideoSynthesis, MultiModalConversation, CodeGeneration |
| `embeddings` | TextEmbedding, MultiModalEmbedding, BatchTextEmbedding |
| `rerank` | TextReRank |
| `app` | Application (Bailian app completion) |
| `audio/asr` | Transcription, Recognition, Vocabulary |
| `audio/tts` | SpeechSynthesizer (websocket) |
| `audio/http_tts` | HttpSpeechSynthesizer |
| `nlp` | Understanding |
| `assistants` | Assistants API |
| `threads` | Threads / Messages / Runs / Steps |
| `finetune` | FineTunes, Deployments |
| `tokenizers` | Tokenization HTTP + local Tokenizer |
| `agentstudio` | AgentStudio client |
| `multimodal` | MultiModalDialog, Tingwu |

## Generation.Call

Call generation model service.

Args:

- model (str): The requested model, such as qwen-turbo.
- prompt (Any): The input prompt.
- history (list): The user provided history, deprecated.
- api_key (str, optional): The api api_key, can be None.
- messages (list): The generation messages.
- plugins (Any): The plugin config, str or dict.
- workspace (str): The dashscope workspace id.
- stream (bool, optional): Enable streaming output.
- temperature (float, optional): Controls randomness, range [0, 2).
- top_p (float, optional): Nucleus sampling, range (0, 1.0].
- top_k (int, optional): Size of candidate token set for sampling.
- max_tokens (int, optional): Maximum output token count.
- seed (int, optional): Random seed for reproducibility.
- stop (str or list, optional): Stop sequences.
- repetition_penalty (float, optional): Penalizes repeated sequences. 1.0 means no penalty.
- presence_penalty (float, optional): Controls content repetition, range [-2.0, 2.0].
- result_format (str, optional): "message" or "text".
- incremental_output (bool, optional): In streaming mode, output only new tokens (True) vs. cumulative output (False).
- enable_search (bool, optional): Enable web search.
- tools (list, optional): Tool definitions for function calling.
- tool_choice (str or dict, optional): Tool selection strategy.
- enable_thinking (bool, optional): Enable thinking mode for hybrid thinking models.
- thinking_budget (int, optional): Maximum token budget for thinking mode.
- n (int, optional): Number of responses to generate (1-4).
- logprobs (bool, optional): Whether to return log probabilities of the output tokens.
- top_logprobs (int, optional): Number of most likely tokens to return at each token position when logprobs is enabled.
- search_options (dict, optional): Configuration options for web search feature.
- parallel_tool_calls (bool, optional): Enable parallel tool calls for function calling.
- response_format (dict, optional): Format constraint for response, e.g., {"type": "json_object"} for JSON mode.

Returns: `GenerationResponse`. If stream is True, use `CallStream`.
