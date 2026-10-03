package ai

import (
	"context"

	"github.com/arvazvi/crawler/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAI struct {
	client *openai.Client
}

func (c *OpenAI) Embed(ctx context.Context, text string) ([]float64, error) {
	resp, err := c.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String("Your text goes here to get vectorized."),
		},
		Model: openai.EmbeddingModelTextEmbedding3Small,
	})

	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}

func NewOpenAIClient(conf *config.Config) *OpenAI {
	client := openai.NewClient(
		option.WithAPIKey(conf.AI.ApiKey),
	)
	return &OpenAI{client: &client}
}