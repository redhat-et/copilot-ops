package cmd

import (
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

// NewAskCmd creates a new ask command which communicates with the ChatGPT API.
func NewAskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     CommandAsk,
		Short:   "Ask ChatGPT a question and get an answer",
		Long:    "Ask talks to the ChatGPT API and returns the response in your terminal.",
		Example: "	copilot-ops ask 'Write a BASH script that checks the weather once every 5 minutes and sends an email if it's raining.'",
		RunE:    RunAsk,
		Args:    cobra.ExactArgs(1),
	}

	// add flags
	// cmd.Flags().String()
	return cmd
}

// RunAsk Runs the command to talk with the ChatGPT API and return the response.
func RunAsk(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()
	var client *openai.Client
	// request
	request := args[0]
	if request == "" {
		return fmt.Errorf("no request provided")
	}
	// grab envs
	openaiKey := os.Getenv("OPENAI_API_KEY")
	openaiOrgID := os.Getenv("OPENAI_ORG_ID")
	// create OpenAI client & send request
	if openaiOrgID != "" {
		openaiConfig := openai.DefaultConfig(openaiKey)
		openaiConfig.OrgID = openaiOrgID
		client = openai.NewClientWithConfig(openaiConfig)
	} else {
		client = openai.NewClient(openaiKey)
	}
	// create chat response
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: "You are a large language model trained by OpenAI.",
			},
			{
				Role:    "user",
				Content: request,
			},
		},
	})
	if err != nil {
		return fmt.Errorf("could not create chat completion: %w", err)
	}
	// print response to STDOUT
	if len(resp.Choices) == 0 {
		return fmt.Errorf("no response from API")
	}
	chatResponse := resp.Choices[0]
	os.Stderr.WriteString(chatResponse.Message.Content + "\n")
	return nil
}
