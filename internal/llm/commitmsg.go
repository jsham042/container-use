package llm

import (
	"context"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// GenerateCommitMessage takes a git diff and returns a semantic commit message
// using OpenAI's GPT-4 model.
func GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	// Get OpenAI client
	client, err := NewClient(ctx)
	if err != nil {
		return "", err
	}

	// Define system prompt with semantic commit rules
	systemPrompt := `You are an expert Git user who writes high quality semantic commit messages.
Format your commit messages according to these rules:

Format: <type>(<scope>): <subject>

<body>

<footer>

Types:
- feat: New feature
- fix: Bug fix
- docs: Documentation changes
- style: Code style changes (formatting, semicolons, etc.)
- refactor: Code refactoring without changing functionality
- perf: Performance improvements
- test: Adding or modifying tests
- build: Build system or dependency changes
- ci: CI/CD configuration changes
- chore: Routine tasks, maintenance

Subject Line Rules:
- Keep under 50 characters
- Use imperative mood ("add feature" not "added feature")
- Start with lowercase
- No period at the end

Body Guidelines:
- Line length: Wrap at 72 characters
- Explain what and why, not how
- Separate subject from body with a blank line
- Use bullet points with - or * for lists

Footer Elements:
- Breaking changes: Start with BREAKING CHANGE:
- Issue references: Use Fixes #123, Closes #456

OUTPUT ONLY THE COMMIT MESSAGE, nothing else.`

	// Create a ChatCompletion request
	req := openai.ChatCompletionRequest{
		Model: selectGPTModel(client),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: "Here is the git diff to summarize in a semantic commit message:\n\n" + diff,
			},
		},
		Temperature: 0.5, // Lower temperature for more deterministic output
	}

	// Make the API call
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	// Extract and trim the message from the response
	if len(resp.Choices) > 0 {
		commitMsg := strings.TrimSpace(resp.Choices[0].Message.Content)
		return commitMsg, nil
	}

	return "", nil
}

// selectGPTModel attempts to use GPT-4o if available, falling back to GPT-4
func selectGPTModel(client *openai.Client) string {
	// First preference is GPT-4o
	return "gpt-4o"
}