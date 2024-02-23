package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/sashabaranov/go-openai"
)

type config struct {
	FirstRunDate string   `env:"FIRST_RUN_DATE,required,notEmpty"`
	LineToken    string   `env:"LINE_TOKEN,required,notEmpty"`
	OpenAIAPIKey string   `env:"OPENAI_API_KEY,required,notEmpty"`
	Prompt       string   `env:"OPENAI_PROMPT,required,notEmpty"`
	ImageURLs    []string `env:"IMAGE_URL,required,notEmpty" envSeparator:","`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(fmt.Errorf("環境変数の設定に失敗しました。: %w", err))
	}

	firstRunDate, err := time.Parse("2006-01-02", cfg.FirstRunDate)
	if err != nil {
		log.Fatal(fmt.Errorf("日付の設定に失敗しました。: %w", err))
	}

	isHoliday, err := calculateDay(firstRunDate, time.Now())
	if err != nil {
		log.Fatal(fmt.Errorf("日付の計算に失敗しました。: %w", err))
	}
	if isHoliday {
		log.Printf("休みです。isHoliday: %t", isHoliday)
		return
	}

	message, err := generateMessage(cfg.OpenAIAPIKey, cfg.Prompt)
	if err != nil {
		log.Fatal(fmt.Errorf("メッセージの生成に失敗しました。: %w", err))
	}

	imageURL, err := getRandomImage(cfg.ImageURLs)
	if err != nil {
		log.Fatal(fmt.Errorf("画像の取得に失敗しました。: %w", err))
	}

	formData := map[string]string{
		"message":        message,
		"imageThumbnail": imageURL,
		"imageFullsize":  imageURL,
	}

	body, contentType, err := createFormData(formData)
	if err != nil {
		log.Fatal(fmt.Errorf("フォームデータの生成に失敗しました。: %w", err))
	}

	req, err := addHeader(cfg.LineToken, body, contentType)
	if err != nil {
		log.Fatal(fmt.Errorf("ヘッダーの追加に失敗しました。: %w", err))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(fmt.Errorf("LINEへのメッセージ送信に失敗しました。: %w", err))
	}
	if resp.StatusCode >= 400 {
		log.Fatal(fmt.Errorf("LINEへのメッセージ送信に失敗しました。ステータスコード: %d", resp.StatusCode))
	}
	defer resp.Body.Close()
}

func calculateDay(firstRunDate time.Time, currentDate time.Time) (bool, error) {
	weekDiff := currentDate.Sub(firstRunDate) / (7 * 24 * time.Hour)
	if weekDiff%(4*7) == 3 {
		return true, nil
	}
	return false, nil
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

func getRandomImage(imageURLs []string) (imageURL string, err error) {
	sliceLen := len(imageURLs)
	if sliceLen == 0 {
		return "", errors.New("画像URLが指定されていません")
	}

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	randInt := r.Intn(sliceLen)
	imageURL = imageURLs[randInt]
	return
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
