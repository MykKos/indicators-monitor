package tg

import (
	"bytes"
	"encoding/json"
	"io"
)

type (
	TelegramMessage struct {
		Channel   string `json:"chat_id"`
		Text      string `json:"text"`
		ParseMode string `json:"parse_mode"`
	}
)

func (message TelegramMessage) Bytes() []byte {
	message.ParseMode = "MarkdownV2"
	b, _ := json.Marshal(message)
	// fmt.Println(string(b))
	return b
}

func (message TelegramMessage) BytesReader() io.Reader {
	return bytes.NewReader(message.Bytes())
}
