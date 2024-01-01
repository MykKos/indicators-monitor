package tg

import (
	"fmt"
	"time"
)

type (
	TelegramResponse struct {
		Ok          bool                   `json:"ok"`
		ErrorCode   int                    `json:"error_code"`
		Description string                 `json:"description"`
		Parameters  map[string]interface{} `json:"parameters"`
		Result      TelegramMessageView    `json:"result"`
	}

	TelegramMessageView struct {
		MessageId int64            `json:"message_id"`
		Chat      TelegramChatView `json:"chat"`
		Date      int64            `json:"date"`
	}

	TelegramChatView struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Title     string `json:"title"`
		Type      string `json:"type"`
	}

	TelegramError interface{}

	TelegramRateLimit struct {
		Code        int
		Description string
		Delay       time.Duration
		ExtraDelay  time.Duration
	}

	TelegramBadRequest struct {
		Code        int
		Description string
	}
)

func (tgResponse TelegramResponse) String() string {
	dm := "unknown destination"
	switch tgResponse.Result.Chat.Type {
	case "channel":
		dm = fmt.Sprintf("channel %d(%s)", tgResponse.Result.Chat.ID, tgResponse.Result.Chat.Title)
	case "private":
		dm = fmt.Sprintf(
			"channel %d(%s %s)", tgResponse.Result.Chat.ID,
			tgResponse.Result.Chat.FirstName, tgResponse.Result.Chat.LastName,
		)
	}
	return fmt.Sprintf(
		"Message to %s was sent. [MessageID: %d]",
		dm, tgResponse.Result.MessageId,
	)
}

func (tgResponse TelegramResponse) GetError() TelegramError {
	switch tgResponse.ErrorCode {
	case 429:
		return NewRateLimitError(tgResponse)
	case 400:
		return NewBadRequestError(tgResponse)
	}
	return nil
}

func NewRateLimitError(tgResponse TelegramResponse) TelegramRateLimit {
	code, desc := tgResponse.ErrorCode, tgResponse.Description
	delay := 0.0
	if v, ok := tgResponse.Parameters["retry_after"]; ok {
		delay = v.(float64)
	}
	return TelegramRateLimit{
		Code:        code,
		Description: desc,
		Delay:       time.Duration(delay) * time.Second,
		ExtraDelay:  5 * time.Second,
	}
}

func (tgResponse TelegramRateLimit) String() string {
	delay := fmt.Sprintf("%s", tgResponse.Delay)
	if tgResponse.ExtraDelay != 0 {
		delay = fmt.Sprintf("%s (+ %s extra)", delay, tgResponse.ExtraDelay)
	}
	return fmt.Sprintf(
		"Rate limited. Your message will be sent after delay %s. [code: %d | description: %s]",
		delay, tgResponse.Code, tgResponse.Description,
	)
}

func NewBadRequestError(tgResponse TelegramResponse) TelegramBadRequest {
	code, desc := tgResponse.ErrorCode, tgResponse.Description
	return TelegramBadRequest{
		Code:        code,
		Description: desc,
	}
}

func (tgResponse TelegramBadRequest) String() string {
	return tgResponse.Description
}
