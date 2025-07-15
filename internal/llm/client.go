package llm

import (
	"context"
	"errors"
	"os"

	"github.com/sashabaranov/go-openai"
)

// NewClient creates and returns an initialized OpenAI client.
// It reads the OpenAI API base URL and key from environment variables.
func NewClient(ctx context.Context) (*openai.Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, errors.New("OPENAI_API_KEY environment variable is not set")
	}

	config := openai.DefaultConfig(apiKey)

	// Set custom API base URL if provided
	apiBase := os.Getenv("OPENAI_API_BASE")
	if apiBase != "" {
		config.BaseURL = apiBase
	}

	client := openai.NewClientWithConfig(config)
	return client, nil
}