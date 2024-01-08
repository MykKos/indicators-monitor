package watcher

import (
	"fmt"
	"indicators-monitor/internal/tg"
	"strings"

	"github.com/MykKos/stock-sdk/pkg/indicators"
)

const (
	MRAZBPm       = "154800244"
	CEXSigChannel = "-1002141332914"
)

type (
	MACDSignal struct {
		S5    string
		S60   string
		S1440 string

		MACD5    *indicators.MACD
		MACD60   *indicators.MACD
		MACD1440 *indicators.MACD
	}
)

func (w *Watcher) SendSignal(token string, ms *MACDSignal) {
	header := ""
	switch ms.S5 {
	case BuySignal:
		header = "📈 BUY"
	case SellSignal:
		header = "📉 SELL"
	}
	headerLine := fmt.Sprintf("%s *%s* signal", header, token)

	macdLines := []string{
		macdline(ms.S5, "5 minutes", ms.MACD5),
		macdline(ms.S60, "1 hour", ms.MACD60),
		macdline(ms.S1440, "1 day", ms.MACD1440),
	}
	macdLine := strings.Join(macdLines, "\n\n")

	msgText := fmt.Sprintf("%s\n\n%s", headerLine, macdLine)

	m, err := w.SignalChannel.SendMessage(tg.TelegramMessage{
		Channel: CEXSigChannel,
		Text:    msgText,
	})

	fmt.Println(m, err)
}

func macdline(signal, macdtf string, macd *indicators.MACD) string {
	macdLine := macd.MACDLine()
	lastMacd := macdLine[len(macdLine)-1]
	prelast := macdLine[len(macdLine)-2]

	line1 := fmt.Sprintf("*MACD for %s* \\(%s\\)", macdtf, signal)
	line2 := fmt.Sprintf(
		"MACD line position: %s, Signal line position: %s",
		escaped(fmt.Sprintf("%.4f", lastMacd.MainLineValue.EMA)),
		escaped(fmt.Sprintf("%.4f", lastMacd.SignalLineValue.EMA)),
	)
	line3 := fmt.Sprintf(
		"Histogram direction: *%s*\\(%s \\-\\> %s\\)",
		HistAge(lastMacd.HistogramValue, prelast.HistogramValue),
		escaped(fmt.Sprintf("%.4f", prelast.HistogramValue.EMA)),
		escaped(fmt.Sprintf("%.4f", lastMacd.HistogramValue.EMA)),
	)

	return fmt.Sprintf("%s\n\n%s\n%s", line1, line2, line3)
}

var (
	escChars = []string{
		"_", "*", "[",
		"]", "(", ")",
		"~", "`", ">",
		"#", "+", "-",
		"=", "|", "{",
		"}", ".", "!",
	}
)

func escaped(source string) string {
	for _, c := range escChars {
		source = strings.ReplaceAll(source, c, fmt.Sprintf("\\%s", c))
	}
	return source
}
