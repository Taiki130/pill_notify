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
	"os"

	"github.com/sashabaranov/go-openai"
)

func main() {
	lineToken := os.Getenv("LINE_TOKEN")
	prompt := os.Getenv("OPENAI_PROMPT")
	openAPIKey := os.Getenv("OPENAI_API_KEY")
	imageThumbnailURL := os.Getenv("IMAGE_THUMBNAIL_URL")
	imageFullsizeURL := os.Getenv("IMAGE_FULLSIZE_URL")

	message, err := generateMessage(openAPIKey, prompt)
	if err != nil {
		log.Fatal(fmt.Errorf("メッセージの生成に失敗しました。: %w", err))
	}

	formData := map[string]string{
		"message":        message,
		"imageThumbnail": imageThumbnailURL,
		"imageFullsize":  imageFullsizeURL,
	}

	body, contentType, err := createFormData(formData)
	if err != nil {
		log.Fatal(err)
		return
	}

	req, err := addHeader(lineToken, body, contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
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

func generateMessage(openAPIKey string, prompt string) (message string, err error) {
	client := openai.NewClient(openAPIKey)
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
