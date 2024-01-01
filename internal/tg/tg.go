package tg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	ApiUrl = "https://api.telegram.org"
)

type (
	TgMessageSender struct {
		Token string
		Url   string
	}
)

func NewSender(token string) *TgMessageSender {
	return &TgMessageSender{
		Token: token,
		Url:   fmt.Sprintf("%s/bot%s", ApiUrl, token),
	}
}

func (tg *TgMessageSender) SendMessage(data TelegramMessage) (TelegramResponse, error) {
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/sendMessage", tg.Url), data.BytesReader())
	req.Header.Add("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)

	if err != nil {
		return TelegramResponse{}, err
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Println(string(body))

	var tgresp TelegramResponse

	json.Unmarshal(body, &tgresp)

	if strings.Contains(tgresp.String(), "Message to unknown destination was sent") {
		fmt.Printf("[DEBUG] Error message: %s\n", data.Text)
	}

	return tgresp, nil
}
