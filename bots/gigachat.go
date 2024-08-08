package bots

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/emicklei/go-restful/v3/log"
	"gopkg.in/telebot.v3"
)

const (
	GigaChatAPIURL = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"
	OAuthURL       = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
)

type GigaChatRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream            bool `json:"stream"`
	RepetitionPenalty int  `json:"repetition_penalty"`
}

type GigaChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   int64  `json:"expires_at"`
}

func handleGenCommand(c telebot.Context) error {
	prompt := c.Text()
	if len(prompt) <= len("/gen ") {
		return c.Send("Пожалуйста, введите текст для генерации.")
	}
	prompt = prompt[len("/gen "):]

	token, err := getOAuthToken()
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка получения токена: %v", err))
	}

	request := GigaChatRequest{
		Model: "GigaChat",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "user", Content: prompt},
		},
		Stream:            false,
		RepetitionPenalty: 1,
	}

	response, err := sendGigaChatRequest(request, token)
	if err != nil {
		return c.Send(fmt.Sprintf("Ошибка генерации текста: %v", err))
	}

	if len(response.Choices) == 0 {
		return c.Send("Ошибка генерации текста: пустой ответ от сервера.")
	}

	return c.Send(response.Choices[0].Message.Content)
}

func sendGigaChatRequest(request GigaChatRequest, token string) (*GigaChatResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("ошибка формирования запроса: %v", err)
	}

	req, err := http.NewRequest("POST", GigaChatAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	var gigaChatResponse GigaChatResponse
	err = json.Unmarshal(body, &gigaChatResponse)
	if err != nil {
		return nil, fmt.Errorf("ошибка обработки ответа: %v", err)
	}

	return &gigaChatResponse, nil
}

func getOAuthToken() (string, error) {
	data := "grant_type=client_credentials&scope=GIGACHAT_API_PERS"
	req, err := http.NewRequest("POST", OAuthURL, bytes.NewBufferString(data))
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("RqUID", os.Getenv("RQUID"))
	req.Header.Set("Authorization", os.Getenv("AUTHORIZATION"))
	log.Printf("URL: %s", req.URL.String())
	log.Printf("Header: %v", req.Header)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	log.Printf("Response Body: %s", string(body))

	var oauthResponse OAuthResponse
	err = json.Unmarshal(body, &oauthResponse)
	if err != nil {
		return "", fmt.Errorf("ошибка обработки ответа: %v", err)
	}
	log.Printf("Token: %s", oauthResponse.AccessToken)
	return oauthResponse.AccessToken, nil
}
