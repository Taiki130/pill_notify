package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/caarlos0/env/v10"
	"github.com/sashabaranov/go-openai"
)

type config struct {
	LineToken         string `env:"LINE_TOKEN,required,notEmpty"`
	OpenAIAPIKey      string `env:"OPENAI_API_KEY,required,notEmpty"`
	Prompt            string `env:"OPENAI_PROMPT,required,notEmpty"`
	ImageThumbnailURL string `env:"IMAGE_THUMBNAIL_URL,required,notEmpty"`
	ImageFullsizeURL  string `env:"IMAGE_FULLSIZE_URL,required,notEmpty"`
}

func main() {

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(fmt.Errorf("環境変数の設定に失敗しました。: %w", err))
	}

	message, err := generateMessage(cfg.OpenAIAPIKey, cfg.Prompt)
	if err != nil {
		log.Fatal(fmt.Errorf("メッセージの生成に失敗しました。: %w", err))
	}

	formData := map[string]string{
		"message":        message,
		"imageThumbnail": cfg.ImageThumbnailURL,
		"imageFullsize":  cfg.ImageFullsizeURL,
	}

	body, contentType, err := createFormData(formData)
	if err != nil {
		log.Fatal(fmt.Errorf("フォームデータの生成に失敗しました。: %w", err))
	}

	req, err := addHeader(cfg.LineToken, body, contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(fmt.Errorf("LINEへのメッセージ送信に失敗しました。: %w", err))
	}
	defer resp.Body.Close()
}

func addHeader(lineToken string, body io.Reader, contentType string) (req *http.Request, err error) {
	method := "POST"
	url := "https://notify-api.line.me/api/notify"

	authHeader := fmt.Sprintf("Bearer %s", lineToken)
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", authHeader)
	return
}

func generateMessage(openAIAPIKey string, prompt string) (message string, err error) {
	client := openai.NewClient(openAIAPIKey)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("OpenAI APIからのレスポンスに選択肢が含まれていません")
	}
	return resp.Choices[0].Message.Content, nil
}

func createFormData(formData map[string]string) (io.Reader, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for k, v := range formData {
		if err := writer.WriteField(k, v); err != nil {
			return nil, "", err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return body, writer.FormDataContentType(), nil
}
