package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OpenAIClient struct {
	apiKey string
	model  string
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error"`
}

type Choice struct {
	Message Message `json:"message"`
}

func NewOpenAIClient(apiKey, model string) *OpenAIClient {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if model == "" {
		model = "gpt-3.5-turbo"
	}
	return &OpenAIClient{
		apiKey: apiKey,
		model:  model,
	}
}

func (c *OpenAIClient) GenerateText(prompt string) (string, error) {
	if c.apiKey == "" {
		return "Mock generated document for: " + prompt, nil
	}

	reqBody := OpenAIRequest{
		Model: c.model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a legal document expert. Generate professional legal documents based on the provided prompts.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var openaiResp OpenAIResponse
	if err := json.Unmarshal(body, &openaiResp); err != nil {
		return "", err
	}

	if openaiResp.Error != nil {
		return "", fmt.Errorf("OpenAI error: %s", openaiResp.Error.Message)
	}

	if len(openaiResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return openaiResp.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) AnalyzeText(text string) (AnalysisResult, error) {
	prompt := fmt.Sprintf(`Analyze the following legal document for potential risks and issues:

%s

Provide:
1. A list of potential risks or concerns
2. Suggestions for improvements
3. A risk score from 0 to 100 (0 = no risk, 100 = high risk)

Format your response as JSON with keys: risks (array), suggestions (array), score (number)`, text)

	content, err := c.GenerateText(prompt)
	if err != nil {
		return AnalysisResult{}, err
	}

	var result AnalysisResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		result.Suggestions = []string{content}
	}

	return result, nil
}

type GroqClient struct {
	apiKey string
	model  string
}

type GroqRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type GroqResponse struct {
	Choices []Choice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func NewGroqClient(apiKey, model string) *GroqClient {
	if apiKey == "" {
		apiKey = os.Getenv("GROQ_API_KEY")
	}
	if model == "" {
		model = "mixtral-8x7b-32768"
	}
	return &GroqClient{
		apiKey: apiKey,
		model:  model,
	}
}

func (c *GroqClient) GenerateText(prompt string) (string, error) {
	if c.apiKey == "" {
		return "Mock generated document for: " + prompt, nil
	}

	reqBody := GroqRequest{
		Model: c.model,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a legal document expert. Generate professional legal documents based on the provided prompts.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var groqResp GroqResponse
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return "", err
	}

	if groqResp.Error != nil {
		return "", fmt.Errorf("Groq error: %s", groqResp.Error.Message)
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("no response from Groq")
	}

	return groqResp.Choices[0].Message.Content, nil
}

func (c *GroqClient) AnalyzeText(text string) (AnalysisResult, error) {
	prompt := fmt.Sprintf(`Analyze the following legal document for potential risks and issues:

%s

Provide:
1. A list of potential risks or concerns
2. Suggestions for improvements
3. A risk score from 0 to 100 (0 = no risk, 100 = high risk)

Format your response as JSON with keys: risks (array), suggestions (array), score (number)`, text)

	content, err := c.GenerateText(prompt)
	if err != nil {
		return AnalysisResult{}, err
	}

	var result AnalysisResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		result.Suggestions = []string{content}
	}

	return result, nil
}

func GetAIClient(provider string) AIClient {
	if provider == "groq" {
		return NewGroqClient("", "")
	}
	return NewOpenAIClient("", "gpt-3.5-turbo")
}
