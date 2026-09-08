// Copyright (c) Alibaba, Inc. and its affiliates.

package utils

import (
	"github.com/liuxiaobopro/dashscope-go/apientities"
)

type accData struct {
	Content          any
	ReasoningContent string
	ToolCalls        []any
	LogprobsContent  []any
	Finished         bool
	FinishReason     string
	AllChoicesSent   bool
	Role             string
}

func getAcc(m map[any]any, idx any) *accData {
	if v, ok := m[idx]; ok {
		if a, ok := v.(*accData); ok {
			return a
		}
	}
	a := &accData{Content: "", ToolCalls: []any{}, LogprobsContent: []any{}}
	m[idx] = a
	return a
}

func asString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func choiceIndex(c apientities.Choice, enumIdx int) int {
	if c.Index != 0 {
		return c.Index
	}
	return enumIdx
}

func accumulateContent(acc *accData, content any) any {
	if content == nil {
		return acc.Content
	}
	switch cur := content.(type) {
	case string:
		if s, ok := acc.Content.(string); ok {
			acc.Content = s + cur
		} else {
			acc.Content = cur
		}
		return acc.Content
	case []any:
		list, ok := acc.Content.([]any)
		if !ok {
			list = []any{}
		}
		for i, item := range cur {
			for len(list) <= i {
				list = append(list, map[string]any{"text": ""})
			}
			if m, ok := item.(map[string]any); ok {
				if t, ok := m["text"].(string); ok && t != "" {
					if lm, ok := list[i].(map[string]any); ok {
						prev, _ := lm["text"].(string)
						lm["text"] = prev + t
						list[i] = lm
					}
				}
			}
		}
		acc.Content = list
		return list
	default:
		return acc.Content
	}
}

func accumulateToolCalls(acc *accData, calls []any) {
	for _, currentCall := range calls {
		m, ok := currentCall.(map[string]any)
		if !ok {
			continue
		}
		idx, ok := m["index"]
		if !ok {
			acc.ToolCalls = append(acc.ToolCalls, currentCall)
			continue
		}
		var existing map[string]any
		for _, accCall := range acc.ToolCalls {
			if am, ok := accCall.(map[string]any); ok {
				if am["index"] == idx {
					existing = am
					break
				}
			}
		}
		if existing == nil {
			copied := map[string]any{}
			for k, v := range m {
				copied[k] = v
			}
			acc.ToolCalls = append(acc.ToolCalls, copied)
			continue
		}
		if fn, ok := m["function"].(map[string]any); ok {
			efn, _ := existing["function"].(map[string]any)
			if efn == nil {
				efn = map[string]any{}
				existing["function"] = efn
			}
			if name, ok := fn["name"].(string); ok {
				prev, _ := efn["name"].(string)
				efn["name"] = prev + name
			}
			if args, ok := fn["arguments"].(string); ok {
				prev, _ := efn["arguments"].(string)
				efn["arguments"] = prev + args
			}
		}
	}
}

// MergeSingleResponse merge a single response chunk with accumulated data.
// Returns true to yield, false to skip. Extra responses for n>1 are appended to extra.
func MergeSingleResponse(parsed *apientities.GenerationResponse, accumulated map[any]any, n int) (yield bool, extra []*apientities.GenerationResponse) {
	if parsed == nil {
		return false, nil
	}
	if n > 1 && len(accumulated) > 0 {
		allSent := true
		hasFlag := false
		for _, v := range accumulated {
			if a, ok := v.(*accData); ok {
				hasFlag = true
				if !a.AllChoicesSent {
					allSent = false
				}
			}
		}
		if hasFlag && allSent {
			return false, nil
		}
	}

	if parsed.Output != nil && parsed.Output.Text != "" && len(parsed.Output.Choices) == 0 {
		acc := getAcc(accumulated, 0)
		if s, ok := acc.Content.(string); ok {
			acc.Content = s + parsed.Output.Text
		} else {
			acc.Content = parsed.Output.Text
		}
		parsed.Output.Text = asString(acc.Content)
		return true, nil
	}

	if parsed.Output == nil || len(parsed.Output.Choices) == 0 {
		return true, nil
	}

	for i := range parsed.Output.Choices {
		choice := &parsed.Output.Choices[i]
		idx := choiceIndex(*choice, i)
		acc := getAcc(accumulated, idx)
		if choice.Message == nil {
			choice.Message = &apientities.Message{Role: "assistant", Content: acc.Content}
			if acc.Role != "" {
				choice.Message.Role = acc.Role
			}
			if acc.ReasoningContent != "" {
				choice.Message.ReasoningContent = acc.ReasoningContent
			}
			if len(acc.ToolCalls) > 0 {
				choice.Message.ToolCalls = acc.ToolCalls
			}
		} else {
			if choice.Message.Role != "" {
				acc.Role = choice.Message.Role
			}
			if choice.Message.Content != nil {
				choice.Message.Content = accumulateContent(acc, choice.Message.Content)
			} else if acc.Content != nil {
				choice.Message.Content = acc.Content
			}
			if choice.Message.ReasoningContent != "" {
				acc.ReasoningContent += choice.Message.ReasoningContent
			}
			if acc.ReasoningContent != "" {
				choice.Message.ReasoningContent = acc.ReasoningContent
			}
			if len(choice.Message.ToolCalls) > 0 {
				accumulateToolCalls(acc, choice.Message.ToolCalls)
				choice.Message.ToolCalls = acc.ToolCalls
			} else if len(acc.ToolCalls) > 0 {
				choice.Message.ToolCalls = acc.ToolCalls
			}
			if acc.Role != "" && choice.Message.Role == "" {
				choice.Message.Role = acc.Role
			}
		}
		if n > 1 && choice.FinishReason != "" && choice.FinishReason != "null" {
			acc.FinishReason = choice.FinishReason
			acc.Finished = true
		}
	}

	if n > 1 {
		finishedCount := 0
		for _, v := range accumulated {
			if a, ok := v.(*accData); ok && a.Finished {
				finishedCount++
			}
		}
		var finishedInPacket []apientities.Choice
		for _, c := range parsed.Output.Choices {
			if c.FinishReason != "" && c.FinishReason != "null" {
				finishedInPacket = append(finishedInPacket, c)
			}
		}
		if len(finishedInPacket) == 0 {
			return true, nil
		}
		first := finishedInPacket[0].FinishReason
		if first == "stop" {
			if finishedCount < n {
				for i := range parsed.Output.Choices {
					if parsed.Output.Choices[i].FinishReason != "" && parsed.Output.Choices[i].FinishReason != "null" {
						parsed.Output.Choices[i].FinishReason = "null"
					}
				}
			}
			return true, nil
		}
		return true, nil
	}
	return true, nil
}

// MergeMultimodalSingleResponse merge multimodal incremental chunks.
func MergeMultimodalSingleResponse(parsed *apientities.MultiModalConversationResponse, accumulated map[any]any, n int) (bool, []*apientities.MultiModalConversationResponse) {
	if parsed == nil || parsed.Output == nil {
		return true, nil
	}
	gen := &apientities.GenerationResponse{
		DashScopeAPIResponse: parsed.DashScopeAPIResponse,
		Output: &apientities.GenerationOutput{
			Text:         parsed.Output.Text,
			Choices:      parsed.Output.Choices,
			FinishReason: parsed.Output.FinishReason,
		},
	}
	yield, _ := MergeSingleResponse(gen, accumulated, n)
	if gen.Output != nil {
		parsed.Output.Text = gen.Output.Text
		parsed.Output.Choices = gen.Output.Choices
		parsed.Output.FinishReason = gen.Output.FinishReason
	}
	return yield, nil
}
