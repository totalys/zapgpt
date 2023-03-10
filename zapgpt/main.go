package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	methodPost            = "POST"
	headerContentType     = "Content-Type"
	headerContentTypeJson = "application/json"
	headerAuth            = "Authorization"
	openApiURL            = "https://api.openai.com/v1/chat/completions"
	openApiModel          = "gpt-3.5-turbo"
	openApiRole           = "user"
	requestBody           = "Body"
)

var errorBodyNotFound = errors.New("Body not found")

var gptApiKey = os.Getenv("GPT_API_KEY")

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index   int `json:"index"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
}

func GenerateGPTText(query string) (string, error) {
	req := Request{
		Model: openApiModel,
		Messages: []Message{
			{
				Role:    openApiRole,
				Content: query,
			},
		},
		MaxTokens: 150,
	}

	reqJson, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest(methodPost, openApiURL, bytes.NewBuffer(reqJson))
	if err != nil {
		return "", err
	}

	request.Header.Set(headerContentType, headerContentTypeJson)
	request.Header.Set(headerAuth, fmt.Sprintf("Bearer %s", gptApiKey))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var resp Response
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil

}

func parseBase64RequestData(r string) (string, error) {
	dataBytes, err := base64.StdEncoding.DecodeString(r)
	if err != nil {
		return "", err
	}

	data, err := url.ParseQuery(string(dataBytes))
	if data.Has(requestBody) {
		return data.Get(requestBody), nil
	}

	return "", errorBodyNotFound
}

func process(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	result, err := parseBase64RequestData(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	text, err := GenerateGPTText(result)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       text,
	}, nil
}

func main() {
	lambda.Start(process)
}
