package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/ducnhd/lotusmarket/go/technical"
	"github.com/ducnhd/lotusmarket/go/types"
)

// ErrNoAPIKey is returned when CLAUDE_API_KEY is not configured.
var ErrNoAPIKey = fmt.Errorf("lotusmarket: CLAUDE_API_KEY not configured. Get your key at https://console.anthropic.com/settings/keys")

const (
	DefaultModel     = "claude-sonnet-4-6"
	DefaultMaxTokens = 4096
)

// Config holds Claude API configuration.
type Config struct {
	APIKey    string // Required. Anthropic API key.
	Model     string // Optional. Default: claude-sonnet-4-6
	MaxTokens int    // Optional. Default: 4096
}

// Client wraps the Anthropic SDK for stock analysis.
type Client struct {
	client    anthropic.Client
	model     string
	maxTokens int
}

// New creates a new AI client. Returns ErrNoAPIKey if APIKey is empty.
func New(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, ErrNoAPIKey
	}
	model := cfg.Model
	if model == "" {
		model = DefaultModel
	}
	maxTokens := cfg.MaxTokens
	if maxTokens <= 0 {
		maxTokens = DefaultMaxTokens
	}
	return &Client{
		client:    anthropic.NewClient(option.WithAPIKey(cfg.APIKey)),
		model:     model,
		maxTokens: maxTokens,
	}, nil
}

// AnalysisResult holds the AI analysis output.
type AnalysisResult struct {
	Text      string // The analysis text (Vietnamese)
	Model     string // Model used
	TokensIn  int    // Input tokens used
	TokensOut int    // Output tokens used
}

// Prompt sends a custom prompt to Claude and returns the response.
func (c *Client) Prompt(ctx context.Context, prompt string) (*AnalysisResult, error) {
	return c.call(ctx, prompt)
}

// AnalyzeTrend analyzes stock price trends with pre-computed technical indicators.
// It fetches nothing — you provide the price data, it builds the prompt and calls Claude.
func (c *Client) AnalyzeTrend(ctx context.Context, ticker string, prices []types.StockData) (*AnalysisResult, error) {
	if len(prices) < 5 {
		return nil, fmt.Errorf("need at least 5 price records, got %d", len(prices))
	}

	// Pre-compute technical indicators locally
	closes := make([]float64, len(prices))
	for i, p := range prices {
		closes[i] = p.Close
	}
	dash := technical.Dashboard(closes)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Bạn là chuyên gia phân tích kỹ thuật chứng khoán Việt Nam. Phân tích xu hướng giá cổ phiếu %s.\n\n", ticker))

	// Add pre-computed indicators
	sb.WriteString("CHỈ SỐ KỸ THUẬT (đã tính toán sẵn):\n")
	sb.WriteString(fmt.Sprintf("- RSI(14): %.1f\n", dash.RSI))
	if dash.MA20 != nil {
		sb.WriteString(fmt.Sprintf("- MA20: %.0f\n", *dash.MA20))
	}
	if dash.MA50 != nil {
		sb.WriteString(fmt.Sprintf("- MA50: %.0f\n", *dash.MA50))
	}
	if dash.MA200 != nil {
		sb.WriteString(fmt.Sprintf("- MA200: %.0f\n", *dash.MA200))
	}
	sb.WriteString(fmt.Sprintf("- Momentum(20): %.1f%%\n", dash.Momentum))
	sb.WriteString(fmt.Sprintf("- Tín hiệu: %s\n", dash.Signal))
	sb.WriteString(fmt.Sprintf("- Điểm: %.0f/100\n\n", dash.Score))

	// Add recent price data (last 20 days)
	sb.WriteString("DỮ LIỆU GIÁ GẦN ĐÂY:\n")
	sb.WriteString("Ngày | Mở | Cao | Thấp | Đóng | KL\n")
	start := 0
	if len(prices) > 20 {
		start = len(prices) - 20
	}
	for _, d := range prices[start:] {
		sb.WriteString(fmt.Sprintf("%s | %.0f | %.0f | %.0f | %.0f | %d\n",
			d.Date, d.Open, d.High, d.Low, d.Close, d.Volume))
	}

	sb.WriteString("\nHãy phân tích bằng tiếng Việt đơn giản:\n")
	sb.WriteString("1. Xu hướng giá ngắn hạn và trung hạn\n")
	sb.WriteString("2. Vùng hỗ trợ và kháng cự\n")
	sb.WriteString("3. Phân tích các đường MA\n")
	sb.WriteString("4. Đánh giá đà tăng/giảm dựa trên RSI và momentum\n")
	sb.WriteString("5. Mức độ rủi ro\n")
	sb.WriteString("\nQUY TẮC: Toàn bộ nội dung PHẢI viết bằng tiếng Việt, dùng từ ngữ dễ hiểu cho nhà đầu tư cá nhân.")

	return c.call(ctx, sb.String())
}

// AskQuestion answers a market-related question in Vietnamese.
func (c *Client) AskQuestion(ctx context.Context, question string) (*AnalysisResult, error) {
	prompt := fmt.Sprintf(`Bạn là chuyên gia tư vấn chứng khoán Việt Nam. Trả lời câu hỏi sau:

%s

Hãy trả lời rõ ràng, dễ hiểu bằng tiếng Việt. Nếu câu hỏi về cổ phiếu cụ thể, đưa ra các điểm dữ liệu liên quan.
QUY TẮC: Toàn bộ câu trả lời PHẢI bằng tiếng Việt, dùng từ ngữ đơn giản cho người mới đầu tư.`, question)

	return c.call(ctx, prompt)
}

// AnalyzeWithContext provides AI analysis with custom data context.
// Use this for advanced scenarios: pass pre-formatted data and analysis instructions.
func (c *Client) AnalyzeWithContext(ctx context.Context, dataContext string, instruction string) (*AnalysisResult, error) {
	prompt := fmt.Sprintf("%s\n\n%s\n\nQUY TẮC: Toàn bộ nội dung PHẢI viết bằng tiếng Việt, dùng từ ngữ dễ hiểu.", dataContext, instruction)
	return c.call(ctx, prompt)
}

func (c *Client) call(ctx context.Context, prompt string) (*AnalysisResult, error) {
	callCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	var lastErr error
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(1<<attempt) * time.Second)
		}
		message, err := c.client.Messages.New(callCtx, anthropic.MessageNewParams{
			Model:     anthropic.Model(c.model),
			MaxTokens: int64(c.maxTokens),
			Messages: []anthropic.MessageParam{
				anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
			},
		})
		if err != nil {
			lastErr = err
			continue
		}
		var text string
		for _, block := range message.Content {
			if block.Type == "text" {
				text += block.Text
			}
		}
		return &AnalysisResult{
			Text:      text,
			Model:     c.model,
			TokensIn:  int(message.Usage.InputTokens),
			TokensOut: int(message.Usage.OutputTokens),
		}, nil
	}
	return nil, fmt.Errorf("Claude API failed after 3 attempts: %w", lastErr)
}
