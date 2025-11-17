package llm

import (
    "context"
    "fmt"
    "time"

    openai "github.com/sashabaranov/go-openai"

    "cliq-hub-backend/internal/config"
)

type Client interface {
    GenerateCliqfileFromPrompt(ctx context.Context, req GenerateRequest) (string, error)
}

type client struct {
    oa    *openai.Client
    model string
}

type GenerateRequest struct {
    CommandExample string
    Name           string
    Description    string
    Author         string
}

func NewClient(cfg *config.Config) (Client, error) {
    conf := openai.DefaultConfig(cfg.LLMAPIKey)
    if cfg.LLMBaseURL != "" {
        conf.BaseURL = cfg.LLMBaseURL
    }
    c := openai.NewClientWithConfig(conf)
    return &client{oa: c, model: cfg.LLMModel}, nil
}

func (c *client) GenerateCliqfileFromPrompt(ctx context.Context, req GenerateRequest) (string, error) {
    sys := "You generate ONLY valid cliqfile YAML per schema. No prose. No markdown fences."
    user := fmt.Sprintf(`Given a CLI command example and optional metadata, generate a complete cliqfile YAML.\n\nRequirements:\n- Fields: name, description, version (\"1.0\"), author, cliq_template_version (\"1.0\"), cmds (with name, description, command, variables).\n- Use Go template placeholders {{variable}} in the command string.\n- Return RAW YAML ONLY (no code fences, no extra text).\n\nInput:\ncommand_example: %s\nname: %s\ndescription: %s\nauthor: %s\n`, req.CommandExample, req.Name, req.Description, req.Author)

    ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
    defer cancel()

    resp, err := c.oa.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
        Model: c.model,
        Messages: []openai.ChatCompletionMessage{
            {Role: openai.ChatMessageRoleSystem, Content: sys},
            {Role: openai.ChatMessageRoleUser, Content: user},
        },
        Temperature: 0,
    })
    if err != nil {
        return "", err
    }
    if len(resp.Choices) == 0 {
        return "", fmt.Errorf("empty LLM response")
    }
    return resp.Choices[0].Message.Content, nil
}