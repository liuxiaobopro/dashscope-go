// Copyright (c) Alibaba, Inc. and its affiliates.

package apientities

import (
	"encoding/json"
	"net/http"
)

// DashScopeAPIResponse the response content.
//
//	request_id (str): The request id.
//	status_code (int): HTTP status code, 200 indicates that the
//	    request was successful, and others indicate an error.
//	code (str): Error code if error occurs, otherwise empty str.
//	message (str): Set to error message on error.
//	output (Any): The request output.
//	usage (Any): The request usage information.
type DashScopeAPIResponse struct {
	StatusCode int               `json:"status_code"`
	RequestID  string            `json:"request_id"`
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Output     any               `json:"output"`
	Usage      any               `json:"usage"`
	Headers    map[string]string `json:"-"`
}

func (r *DashScopeAPIResponse) String() string {
	b, _ := json.Marshal(map[string]any{
		"status_code": r.StatusCode,
		"request_id":  r.RequestID,
		"code":        r.Code,
		"message":     r.Message,
		"output":      r.Output,
		"usage":       r.Usage,
	})
	return string(b)
}

func (r *DashScopeAPIResponse) OutputMap() map[string]any {
	if r == nil || r.Output == nil {
		return nil
	}
	if m, ok := r.Output.(map[string]any); ok {
		return m
	}
	return nil
}

func (r *DashScopeAPIResponse) UsageMap() map[string]any {
	if r == nil || r.Usage == nil {
		return nil
	}
	if m, ok := r.Usage.(map[string]any); ok {
		return m
	}
	return nil
}

func (r *DashScopeAPIResponse) TaskID() string {
	m := r.OutputMap()
	if m == nil {
		return ""
	}
	if v, ok := m["task_id"].(string); ok {
		return v
	}
	return ""
}

func (r *DashScopeAPIResponse) TaskStatus() string {
	m := r.OutputMap()
	if m == nil {
		return ""
	}
	if v, ok := m["task_status"].(string); ok {
		return v
	}
	return ""
}

func okStatus(code int) bool {
	return code == http.StatusOK
}

// Role message role.
const (
	RoleUser       = "user"
	RoleSystem     = "system"
	RoleBot        = "bot"
	RoleAssistant  = "assistant"
	RoleAttachment = "attachment"
)

// Message chat message.
type Message struct {
	Role             string `json:"role"`
	Content          any    `json:"content,omitempty"`
	ReasoningContent string `json:"reasoning_content,omitempty"`
	Name             string `json:"name,omitempty"`
	ToolCallID       string `json:"tool_call_id,omitempty"`
	ToolCalls        []any  `json:"tool_calls,omitempty"`
	Extra            map[string]any `json:"-"`
}

// NewMessage create a message.
func NewMessage(role string, content any) Message {
	return Message{Role: role, Content: content}
}

// Choice generation choice.
type Choice struct {
	FinishReason string         `json:"finish_reason"`
	Message      *Message       `json:"message,omitempty"`
	Index        int            `json:"index,omitempty"`
	Logprobs     map[string]any `json:"logprobs,omitempty"`
}

// Audio multimodal audio.
type Audio struct {
	Data      string `json:"data,omitempty"`
	URL       string `json:"url,omitempty"`
	ID        string `json:"id,omitempty"`
	ExpiresAt int    `json:"expires_at,omitempty"`
}

// GenerationOutput generation output.
type GenerationOutput struct {
	Text         string   `json:"text,omitempty"`
	Choices      []Choice `json:"choices,omitempty"`
	FinishReason string   `json:"finish_reason,omitempty"`
	Raw          map[string]any `json:"-"`
}

// GenerationUsage generation usage.
type GenerationUsage struct {
	InputTokens  int            `json:"input_tokens"`
	OutputTokens int            `json:"output_tokens"`
	TotalTokens  int            `json:"total_tokens,omitempty"`
	Raw          map[string]any `json:"-"`
}

// GenerationResponse generation response.
type GenerationResponse struct {
	DashScopeAPIResponse
	Output *GenerationOutput `json:"output"`
	Usage  *GenerationUsage  `json:"usage"`
}

// GenerationResponseFromAPI convert DashScopeAPIResponse to GenerationResponse.
func GenerationResponseFromAPI(api *DashScopeAPIResponse) *GenerationResponse {
	r := &GenerationResponse{DashScopeAPIResponse: *api}
	if api.StatusCode == http.StatusOK && api.Output != nil {
		r.Output = parseGenerationOutput(api.Output)
		if api.Usage != nil {
			r.Usage = parseGenerationUsage(api.Usage)
		} else {
			r.Usage = &GenerationUsage{}
		}
	}
	return r
}

func parseGenerationOutput(v any) *GenerationOutput {
	b, _ := json.Marshal(v)
	out := &GenerationOutput{}
	_ = json.Unmarshal(b, out)
	if m, ok := v.(map[string]any); ok {
		out.Raw = m
	}
	return out
}

func parseGenerationUsage(v any) *GenerationUsage {
	b, _ := json.Marshal(v)
	u := &GenerationUsage{}
	_ = json.Unmarshal(b, u)
	if m, ok := v.(map[string]any); ok {
		u.Raw = m
		if t, ok := m["total_tokens"].(float64); ok {
			u.TotalTokens = int(t)
		}
	}
	return u
}

// MultiModalConversationOutput multimodal conversation output.
type MultiModalConversationOutput struct {
	Text         string   `json:"text,omitempty"`
	FinishReason string   `json:"finish_reason,omitempty"`
	Choices      []Choice `json:"choices,omitempty"`
	Audio        *Audio   `json:"audio,omitempty"`
}

// MultiModalConversationUsage multimodal conversation usage.
type MultiModalConversationUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	Characters   int `json:"characters"`
}

// MultiModalConversationResponse multimodal conversation response.
type MultiModalConversationResponse struct {
	DashScopeAPIResponse
	Output *MultiModalConversationOutput `json:"output"`
	Usage  *MultiModalConversationUsage  `json:"usage"`
}

// MultiModalConversationResponseFromAPI convert API response.
func MultiModalConversationResponseFromAPI(api *DashScopeAPIResponse) *MultiModalConversationResponse {
	r := &MultiModalConversationResponse{DashScopeAPIResponse: *api}
	if api.StatusCode == http.StatusOK && api.Output != nil {
		b, _ := json.Marshal(api.Output)
		out := &MultiModalConversationOutput{}
		_ = json.Unmarshal(b, out)
		r.Output = out
		u := &MultiModalConversationUsage{}
		if api.Usage != nil {
			ub, _ := json.Marshal(api.Usage)
			_ = json.Unmarshal(ub, u)
		}
		r.Usage = u
	}
	return r
}

// ConversationResponse conversation response (alias of GenerationResponse).
type ConversationResponse = GenerationResponse

// TranscriptionOutput transcription output.
type TranscriptionOutput struct {
	TaskID     string `json:"task_id"`
	TaskStatus string `json:"task_status"`
}

// TranscriptionUsage transcription usage.
type TranscriptionUsage struct {
	Raw map[string]any `json:"-"`
}

// TranscriptionResponse transcription response.
type TranscriptionResponse struct {
	DashScopeAPIResponse
	Output *TranscriptionOutput `json:"output"`
	Usage  *TranscriptionUsage  `json:"usage"`
}

// TranscriptionResponseFromAPI convert API response.
func TranscriptionResponseFromAPI(api *DashScopeAPIResponse) *TranscriptionResponse {
	r := &TranscriptionResponse{DashScopeAPIResponse: *api}
	if api.StatusCode == http.StatusOK {
		if api.Output != nil {
			b, _ := json.Marshal(api.Output)
			out := &TranscriptionOutput{}
			_ = json.Unmarshal(b, out)
			r.Output = out
		}
		if api.Usage != nil {
			if m, ok := api.Usage.(map[string]any); ok {
				r.Usage = &TranscriptionUsage{Raw: m}
			} else {
				r.Usage = &TranscriptionUsage{}
			}
		}
	}
	return r
}

// RecognitionOutput recognition output.
type RecognitionOutput struct {
	Sentence any `json:"sentence"`
}

// RecognitionUsage recognition usage.
type RecognitionUsage struct {
	Duration int `json:"duration"`
}

// RecognitionResponse recognition response.
type RecognitionResponse struct {
	DashScopeAPIResponse
	Output *RecognitionOutput `json:"output"`
	Usage  *RecognitionUsage  `json:"usage"`
}

// RecognitionResponseFromAPI convert API response.
func RecognitionResponseFromAPI(api *DashScopeAPIResponse) *RecognitionResponse {
	r := &RecognitionResponse{DashScopeAPIResponse: *api}
	if api.StatusCode == http.StatusOK {
		if api.Output != nil {
			if m, ok := api.Output.(map[string]any); ok {
				if _, has := m["sentence"]; has {
					b, _ := json.Marshal(api.Output)
					out := &RecognitionOutput{}
					_ = json.Unmarshal(b, out)
					r.Output = out
				}
			}
		}
		if api.Usage != nil {
			b, _ := json.Marshal(api.Usage)
			u := &RecognitionUsage{}
			_ = json.Unmarshal(b, u)
			r.Usage = u
		}
	}
	return r
}

// IsSentenceEnd determine whether the speech recognition result is the end of a sentence.
func IsSentenceEnd(sentence map[string]any) bool {
	if sentence == nil {
		return false
	}
	v, ok := sentence["end_time"]
	return ok && v != nil
}

// SpeechSynthesisOutput speech synthesis output.
type SpeechSynthesisOutput struct {
	Sentence map[string]any `json:"sentence"`
}

// SpeechSynthesisUsage speech synthesis usage.
type SpeechSynthesisUsage struct {
	Characters int `json:"characters"`
}

// SpeechSynthesisResponse speech synthesis response.
type SpeechSynthesisResponse struct {
	DashScopeAPIResponse
	Output *SpeechSynthesisOutput `json:"output"`
	Usage  *SpeechSynthesisUsage  `json:"usage"`
}

// SpeechSynthesisResponseFromAPI convert API response.
func SpeechSynthesisResponseFromAPI(api *DashScopeAPIResponse) *SpeechSynthesisResponse {
	r := &SpeechSynthesisResponse{DashScopeAPIResponse: *api}
	if api.StatusCode == http.StatusOK {
		if api.Output != nil {
			b, _ := json.Marshal(api.Output)
			out := &SpeechSynthesisOutput{}
			_ = json.Unmarshal(b, out)
			r.Output = out
		}
		if api.Usage != nil {
			b, _ := json.Marshal(api.Usage)
			u := &SpeechSynthesisUsage{}
			_ = json.Unmarshal(b, u)
			r.Usage = u
		}
	}
	return r
}

// ImageSynthesisResult image synthesis result.
type ImageSynthesisResult struct {
	URL string `json:"url"`
}

// ImageSynthesisOutput image synthesis output.
type ImageSynthesisOutput struct {
	TaskID     string                  `json:"task_id"`
	TaskStatus string                  `json:"task_status"`
	Results    []ImageSynthesisResult  `json:"results"`
}

// VideoSynthesisOutput video synthesis output.
type VideoSynthesisOutput struct {
	TaskID     string `json:"task_id"`
	TaskStatus string `json:"task_status"`
	VideoURL   string `json:"video_url"`
}

// ImageSynthesisUsage image synthesis usage.
type ImageSynthesisUsage struct {
	ImageCount int `json:"image_count"`
}

// VideoSynthesisUsage video synthesis usage.
type VideoSynthesisUsage struct {
	VideoCount    int    `json:"video_count"`
	VideoDuration int    `json:"video_duration"`
	VideoRatio    string `json:"video_ratio"`
}

// ImageSynthesisResponse image synthesis response.
type ImageSynthesisResponse struct {
	DashScopeAPIResponse
	Output *ImageSynthesisOutput `json:"output"`
	Usage  *ImageSynthesisUsage  `json:"usage"`
}

// ImageSynthesisResponseFromAPI convert API response.
func ImageSynthesisResponseFromAPI(api *DashScopeAPIResponse) *ImageSynthesisResponse {
	r := &ImageSynthesisResponse{DashScopeAPIResponse: *api}
	if api.StatusCode == http.StatusOK {
		if api.Output != nil {
			b, _ := json.Marshal(api.Output)
			out := &ImageSynthesisOutput{}
			_ = json.Unmarshal(b, out)
			r.Output = out
		}
		if api.Usage != nil {
			b, _ := json.Marshal(api.Usage)
			u := &ImageSynthesisUsage{}
			_ = json.Unmarshal(b, u)
			r.Usage = u
		}
	}
	return r
}

// VideoSynthesisResponse video synthesis response.
type VideoSynthesisResponse struct {
	DashScopeAPIResponse
	Output *VideoSynthesisOutput `json:"output"`
	Usage  *VideoSynthesisUsage  `json:"usage"`
}

// VideoSynthesisResponseFromAPI convert API response.
func VideoSynthesisResponseFromAPI(api *DashScopeAPIResponse) *VideoSynthesisResponse {
	r := &VideoSynthesisResponse{DashScopeAPIResponse: *api}
	if api.StatusCode == http.StatusOK {
		if api.Output != nil {
			b, _ := json.Marshal(api.Output)
			out := &VideoSynthesisOutput{}
			_ = json.Unmarshal(b, out)
			r.Output = out
		}
		if api.Usage != nil {
			b, _ := json.Marshal(api.Usage)
			u := &VideoSynthesisUsage{}
			_ = json.Unmarshal(b, u)
			r.Usage = u
		}
	}
	return r
}

// ReRankResult rerank result.
type ReRankResult struct {
	Index          int            `json:"index"`
	RelevanceScore float64        `json:"relevance_score"`
	Document       map[string]any `json:"document,omitempty"`
}

// ReRankOutput rerank output.
type ReRankOutput struct {
	Results []ReRankResult `json:"results"`
}

// ReRankUsage rerank usage.
type ReRankUsage struct {
	TotalTokens int `json:"total_tokens"`
}

// ReRankResponse rerank response.
type ReRankResponse struct {
	DashScopeAPIResponse
	Output *ReRankOutput `json:"output"`
	Usage  *ReRankUsage  `json:"usage"`
}

// ReRankResponseFromAPI convert API response.
func ReRankResponseFromAPI(api *DashScopeAPIResponse) *ReRankResponse {
	r := &ReRankResponse{DashScopeAPIResponse: *api}
	if api.StatusCode == http.StatusOK && api.Output != nil {
		b, _ := json.Marshal(api.Output)
		out := &ReRankOutput{}
		_ = json.Unmarshal(b, out)
		r.Output = out
		u := &ReRankUsage{}
		if api.Usage != nil {
			ub, _ := json.Marshal(api.Usage)
			_ = json.Unmarshal(ub, u)
		}
		r.Usage = u
	}
	return r
}

// TextToSpeechAudio TTS audio.
type TextToSpeechAudio struct {
	ExpiresAt int    `json:"expires_at"`
	ID        string `json:"id"`
	Data      string `json:"data,omitempty"`
	URL       string `json:"url,omitempty"`
}

// TextToSpeechOutput TTS output.
type TextToSpeechOutput struct {
	FinishReason string            `json:"finish_reason,omitempty"`
	Audio        *TextToSpeechAudio `json:"audio,omitempty"`
}

// TextToSpeechResponse TTS response.
type TextToSpeechResponse struct {
	DashScopeAPIResponse
	Output *TextToSpeechOutput             `json:"output"`
	Usage  *MultiModalConversationUsage    `json:"usage"`
}

// TextToSpeechResponseFromAPI convert API response.
func TextToSpeechResponseFromAPI(api *DashScopeAPIResponse) *TextToSpeechResponse {
	r := &TextToSpeechResponse{DashScopeAPIResponse: *api}
	if api.StatusCode == http.StatusOK && api.Output != nil {
		b, _ := json.Marshal(api.Output)
		out := &TextToSpeechOutput{}
		_ = json.Unmarshal(b, out)
		r.Output = out
		u := &MultiModalConversationUsage{}
		if api.Usage != nil {
			ub, _ := json.Marshal(api.Usage)
			_ = json.Unmarshal(ub, u)
		}
		r.Usage = u
	}
	return r
}

// ImageGenerationOutput image generation output.
type ImageGenerationOutput struct {
	Text         string   `json:"text,omitempty"`
	FinishReason string   `json:"finish_reason,omitempty"`
	Choices      []Choice `json:"choices,omitempty"`
	Audio        *Audio   `json:"audio,omitempty"`
}

// ImageGenerationUsage image generation usage.
type ImageGenerationUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	Characters   int `json:"characters"`
}

// ImageGenerationResponse image generation response.
type ImageGenerationResponse struct {
	DashScopeAPIResponse
	Output *ImageGenerationOutput `json:"output"`
	Usage  *ImageGenerationUsage  `json:"usage"`
}

// ImageGenerationResponseFromAPI convert API response.
func ImageGenerationResponseFromAPI(api *DashScopeAPIResponse) *ImageGenerationResponse {
	r := &ImageGenerationResponse{DashScopeAPIResponse: *api}
	if api.StatusCode == http.StatusOK && api.Output != nil {
		b, _ := json.Marshal(api.Output)
		out := &ImageGenerationOutput{}
		_ = json.Unmarshal(b, out)
		r.Output = out
		u := &ImageGenerationUsage{}
		if api.Usage != nil {
			ub, _ := json.Marshal(api.Usage)
			_ = json.Unmarshal(ub, u)
		}
		r.Usage = u
	}
	return r
}

var _ = okStatus
