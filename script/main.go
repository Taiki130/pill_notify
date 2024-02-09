package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	lineToken := os.Getenv("LINE_TOKEN")
	message := os.Getenv("MESSAGE")
	imageThumbnailURL := os.Getenv("IMAGE_THUMBNAIL_URL")
	imageFullsizeURL := os.Getenv("IMAGE_FULLSIZE_URL")
	url := "https://notify-api.line.me/api/notify"
	method := "POST"

	formData := map[string]string{
		"message":        message,
		"imageThumbnail": imageThumbnailURL,
		"imageFullsize":  imageFullsizeURL,
	}

	body, contentType, err := createFormData(formData)
	if err != nil {
		fmt.Println(err)
		return
	}

	authHeader := fmt.Sprintf("Bearer %s", lineToken)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", authHeader)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
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
