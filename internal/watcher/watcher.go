package watcher

import (
	"encoding/json"
	"fmt"
	"indicators-monitor/internal/metrics"
	"indicators-monitor/internal/metrics/influx"
	"indicators-monitor/internal/tg"
	"os"
	"sync"
	"time"

	"github.com/MykKos/stock-sdk/pkg/indicators"
	"github.com/MykKos/stock-sdk/pkg/models"
	"github.com/MykKos/stock-sdk/pkg/stock-api/kcs"
	"github.com/google/uuid"
)

const (
	OldRedHist    = "old bearish"
	GrowthRedHist = "growth bearish"
	ZeroBearish   = "zero bearish"

	OldGreenHist    = "old bullish"
	GrowthGreenHist = "growth bullish"
	ZeroBullish     = "zero bullish"

	Lower  = "lower"
	Bigger = "bigger"
	Equals = "equals"

	BuySignal  = "buy"
	SellSignal = "sell"
	NoSignal   = "no signal"

	SellStrategy = "sell"
	BuyStrategy  = "buy"

	MACDStrategy = "macd"

	FiveMinuteTF = "5m"

	Precision  = 0.001
	PrecDigits = 4

	GeneralPercent = 0.3
	FiveMinPercent = 0.2

	PrevPoints = 3
)

type (
	Watcher struct {
		KCS *kcs.KcsApiCaller

		Tokens []string

		SignalChannel *tg.TgMessageSender

		Limiters map[string]SignalLimiter

		sync.Mutex

		IC *influx.Client
	}

	SignalLimiter struct {
		Token   string
		Signals map[string]struct{}
	}
)

func NewWatcher(tokens []string) *Watcher {
	ic, _ := influx.NewFromConfig(influx.InfluxConfig{
		Url:       "http://89.104.66.163:8086",
		Database:  "crypto-sig",
		Precision: "ns",
	})

	bots := tgbots()

	w := &Watcher{
		KCS: kcs.NewKcsApiCaller("", "", ""),

		SignalChannel: tg.NewSender(bots.SignalBot),

		Tokens: tokens,

		IC: ic,
	}

	w.Limiters = map[string]SignalLimiter{}

	return w
}

type (
	TgBots struct {
		SignalBot string `json:"signal_bot"`
	}
)

func tgbots() TgBots {
	f, _ := os.ReadFile("configs/tg-config.json")
	var bots TgBots
	json.Unmarshal(f, &bots)

	return bots
}

func (w *Watcher) Start() {
	for _, token := range w.Tokens {
		if token == "" {
			continue
		}
		go w.Signals(token)
	}
}

func (w *Watcher) Signals(token string) {
	for ; ; time.Sleep(1 * time.Minute) {
		fmt.Printf("[DEBUG] Trying to find signals for %s\n", token)
		hasSignals := w.WatchToken(token)
		if hasSignals {
			fmt.Printf("[INFO] Have found signals for %s\n", token)
		}
	}
}

func (w *Watcher) WatchToken(token string) bool {
	ms := &MACDSignal{}

	var last5 float64

	ms.MACD5, last5 = w.SearchTF(token, models.FiveMinuteTF())
	ms.MACD60, _ = w.SearchTF(token, models.OneHourTF())
	ms.MACD1440, _ = w.SearchTF(token, models.OneDayTF())

	macd5 := ms.MACD5.MACDLine()
	macd60 := ms.MACD60.MACDLine()
	macd1440 := ms.MACD1440.MACDLine()

	ms.S5 = MacdType(token, "5m", macd5)
	ms.S60 = MacdType(token, "1h", macd60)
	ms.S1440 = MacdType(token, "1d", macd1440)
	defer func() {
		macd5Last := macd5[len(macd5)-1]
		macd60Last := macd60[len(macd60)-1]
		macd1440Last := macd1440[len(macd1440)-1]

		w.MacdMetrics(token, "5m", ms.S5, macd5Last)
		w.MacdMetrics(token, "1h", ms.S60, macd60Last)
		w.MacdMetrics(token, "1d", ms.S1440, macd1440Last)

		fmt.Printf("[%s] 5m: %s, 1h: %s, 1d: %s\n", token, ms.S5, ms.S60, ms.S1440)
	}()
	if ms.S5 == NoSignal || ms.S60 == NoSignal {
		return false
	}
	if ms.S5 != ms.S60 {
		return false
	}

	w.Lock()

	if _, ok := w.Limiters[token]; !ok {
		w.Limiters[token] = SignalLimiter{
			Token:   token,
			Signals: map[string]struct{}{},
		}
	}

	if _, ok := w.Limiters[token].Signals[fmt.Sprint(last5)]; ok {
		w.Unlock()
		return false
	}

	w.Limiters[token].Signals[fmt.Sprint(last5)] = struct{}{}

	w.Unlock()

	w.SendSignal(token, ms)
	return true
}

func (w *Watcher) MacdMetrics(token, tf, signal string, macd *indicators.MACDLine) {
	pt := metricPoint(token, tf, signal, macd)
	w.IC.HitSave(pt)
}

func metricPoint(token, tf, signal string, macd *indicators.MACDLine) metrics.Point {
	pt := metrics.Point{
		Table: "macd-signals",
		Tags: map[string]string{
			"token":     token,
			"timeframe": tf,
			"signal":    signal,
		},
		Fields: map[string]interface{}{
			"uuid":           uuid.NewString(),
			"macd-main":      macd.MainLineValue.EMA,
			"macd-signal":    macd.SignalLineValue.EMA,
			"macd-histogram": macd.HistogramValue.EMA,
		},
	}

	return pt
}

func (w *Watcher) SearchTF(token string, tf models.TimeFrame) (*indicators.MACD, float64) {
	prices := w.KCS.GetTokenPrices(token, tf)
	lastPrice := prices.Prices[len(prices.Prices)-1]

	indprices := w.IndicatorsPrices(prices.Prices)

	macd := indicators.CalculateMACD(indprices, indicators.MACDSetup{
		SlowLine:   indicators.EMASetup{Period: 26, Smooth: 9},
		FastLine:   indicators.EMASetup{Period: 12, Smooth: 9},
		SignalLine: indicators.EMASetup{Period: 9, Smooth: 9},
	})

	return macd, lastPrice.Timestamp
}

func (w *Watcher) IndicatorsPrices(prices []models.PriceValues) []indicators.Price {
	indprices := []indicators.Price{}
	for _, price := range prices {
		indprices = append(indprices, indicators.Price{
			OpenPrice:  price.OpenPrice,
			ClosePrice: price.ClosePrice,
			MaxPrice:   price.MaxPrice,
			MinPrice:   price.MinPrice,
		})
	}
	return indprices
}

func MacdType(token, tf string, macdLine []*indicators.MACDLine) string {
	// macdLine := source.MACDLine()
	fmt.Printf("MACD line len for token %s, %s: %d\n", token, tf, len(macdLine))
	switch {
	case SellMacd(token, tf, macdLine):
		return SellSignal
	case BuyMacd(token, tf, macdLine):
		return BuySignal
	}
	return NoSignal
}

func HistAge(source *indicators.EmaPoint, prev *indicators.EmaPoint) string {
	switch {
	case source.EMA > 0:
		return GreenAge(source, prev)
	case source.EMA < 0:
		return RedAge(source, prev)
	}
	return ""
}

func GreenAge(source *indicators.EmaPoint, prev *indicators.EmaPoint) string {
	switch PrecCompare(source.EMA, prev.EMA) {
	case Bigger:
		return GrowthGreenHist
	case Lower:
		return OldGreenHist
	case Equals:
		return ZeroBullish
	}
	return ""
}

func RedAge(source *indicators.EmaPoint, prev *indicators.EmaPoint) string {
	switch PrecCompare(source.EMA, prev.EMA) {
	case Bigger:
		return OldRedHist
	case Lower:
		return GrowthRedHist
	case Equals:
		return ZeroBearish
	}
	return ""
}

func PrecCompare(source, dest float64) string {
	diff := source - dest
	switch {
	case diff < -Precision:
		return Lower
	case diff > Precision:
		return Bigger
	}
	return Equals
}
