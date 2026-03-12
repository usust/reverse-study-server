package prompt

import (
	"context"
	"errors"
	"reverse-study-server/internal/bootstrap"
	promptrepo "reverse-study-server/internal/repository/prompt"
	chatapi "reverse-study-server/volcengine/chat"
	"reverse-study-server/volcengine/repository"
	"strings"

	"gorm.io/gorm"
)

const createCCodePromptName = "创建C代码提示词"

// GenerateCCode 使用火山引擎模型生成一段 C 源代码。
func GenerateCCode(ctx context.Context, req chatapi.PromptChatRequest) (chatapi.PromptChatResponse, error) {
	if strings.TrimSpace(req.Prompt) == "" {
		dbPrompt, err := promptrepo.GetByName(ctx, bootstrap.GormDB, createCCodePromptName)
		if err == nil && strings.TrimSpace(dbPrompt.Content) != "" {
			req.Prompt = dbPrompt.Content
		} else if err == nil {
			req.Prompt = ``
		} else if errors.Is(err, gorm.ErrRecordNotFound) {

		} else {
			return chatapi.PromptChatResponse{}, err
		}
	}

	apiInfo, err := repository.GetAPIByID(ctx, bootstrap.GormDB, req.ModelAPIID)
	if err != nil {
		return chatapi.PromptChatResponse{}, err
	}

	chatResp, err := chatapi.PromptChat(ctx, chatapi.PromptChatRequest{
		Config: apiInfo,
		Prompt: req.Prompt,
	})
	if err != nil {
		return chatapi.PromptChatResponse{}, err
	}

	cCode := stripMarkdownCodeFence(chatResp.Response.Content)

	return chatapi.PromptChatResponse{
		Request: chatapi.PromptChatRequestInfo{
			ID:           req.ModelAPIID,
			Name:         req.Config.Name,
			Provider:     req.Config.Provider,
			Model:        req.Config.APIModel,
			DisplayModel: req.Config.Model,
		},
		Response: chatapi.PromptChatResult{
			Content:          cCode,
			ReasoningContent: chatResp.Response.ReasoningContent,
			FinishReason:     chatResp.Response.FinishReason,
		},
		Info: chatapi.PromptChatInfo{
			RequestID:   chatResp.Info.RequestID,
			Object:      chatResp.Info.Object,
			Created:     chatResp.Info.Created,
			ServiceTier: chatResp.Info.ServiceTier,
		},
		Usage: chatapi.ExecutePromptUsage{
			PromptTokens:         chatResp.Usage.PromptTokens,
			CompletionTokens:     chatResp.Usage.CompletionTokens,
			TotalTokens:          chatResp.Usage.TotalTokens,
			ReasoningTokens:      chatResp.Usage.ReasoningTokens,
			CachedPromptTokens:   chatResp.Usage.CachedPromptTokens,
			ProvisionedTokensIn:  chatResp.Usage.ProvisionedTokensIn,
			ProvisionedTokensOut: chatResp.Usage.ProvisionedTokensOut,
		}}, nil
}

// stripMarkdownCodeFence 尝试去掉模型偶尔返回的 Markdown 代码块包裹。
func stripMarkdownCodeFence(content string) string {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "```") {
		return trimmed
	}

	lines := strings.Split(trimmed, "\n")
	if len(lines) == 0 {
		return trimmed
	}

	if strings.HasPrefix(lines[0], "```") {
		lines = lines[1:]
	}
	if len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "```" {
		lines = lines[:len(lines)-1]
	}

	return strings.TrimSpace(strings.Join(lines, "\n"))
}
