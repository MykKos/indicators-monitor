package watcher

import (
	"fmt"

	"github.com/MykKos/stock-sdk/pkg/indicators"
)

func BuyMacd(token, tf string, macdLine []*indicators.MACDLine) bool {
	lastMacd := macdLine[len(macdLine)-1]
	prelast := macdLine[len(macdLine)-2]

	prevpts := macdLine[len(macdLine)-PrevPoints-1:]

	fmt.Printf(
		"[%s] TF: %s. Last: %.6f | %.6f | %.6f, Prelast: %.6f | %.6f | %.6f. Age: %s\n",
		token, tf,
		lastMacd.MainLineValue.EMA, lastMacd.SignalLineValue.EMA, lastMacd.HistogramValue.EMA,
		prelast.MainLineValue.EMA, prelast.SignalLineValue.EMA, prelast.HistogramValue.EMA,
		HistAge(lastMacd.HistogramValue, prelast.HistogramValue),
	)

	if NotMACDBuyPos(tf, lastMacd, macdLine) {
		return false
	}

	// если главная линия ниже сигнальной (находимся в зоне покупателей)
	// или нет точки в выбраном периоде, в которой сигнальная линия выше главной (находимся в зоне продавцов)
	// т.е. нет пересечения
	// сигнала нет
	if lastMacd.MainLineValue.EMA < lastMacd.SignalLineValue.EMA || allSmaller(prevpts) {
		return false
	}

	histAge := HistAge(lastMacd.HistogramValue, prelast.HistogramValue)
	if histAge == GrowthRedHist || histAge == OldGreenHist {
		return false
	}

	// если мы находимся в зоне покупателей
	// и 10 точек назад главная линия была выше, чем сейчас
	// то тренд нисходящий, сигнала на покупку нет
	// trendDirPts := macdLine[len(macdLine)-10]
	// if lastMacd.MainLineValue.EMA > 0 && lastMacd.MainLineValue.EMA < trendDirPts.MainLineValue.EMA {
	// 	return false
	// }
	return true
}

func NotMACDBuyPos(tf string, lastMacd *indicators.MACDLine, macdLine []*indicators.MACDLine) bool {
	// получаем минимальный MACD за период
	mline := MaxMacd(macdLine)
	var cond float64
	// получаем значение условия
	switch tf {
	case FiveMinuteTF:
		// FiveMinPercent процентов от максимального (на текущий момент 50%)
		// если макс = 1 -> cond = -0.5
		cond = mline.MainLineValue.EMA * FiveMinPercent
	default:
		// GeneralPercent процентов от максимального (на текущий момент 40%)
		// если макс = -1 -> cond = -0.4
		cond = mline.MainLineValue.EMA * GeneralPercent
	}
	// возвращаем текущее положение MACD
	// если оно выше сравниваемого значения, то сигнала на покупку нет
	return lastMacd.MainLineValue.EMA > cond
}

func MaxMacd(macdLine []*indicators.MACDLine) *indicators.MACDLine {
	max := &indicators.MACDLine{
		MainLineValue: &indicators.EmaPoint{},
	}
	for _, ln := range macdLine {
		if ln.MainLineValue.EMA > max.MainLineValue.EMA {
			max = ln
		}
	}

	return max
}

func allSmaller(comp []*indicators.MACDLine) bool {
	for _, vl := range comp {
		// если есть точка, в которой главная линия ниже сигнальной
		// то есть пересечение, сигнал есть
		if vl.MainLineValue.EMA < vl.SignalLineValue.EMA {
			return false
		}
	}
	return true
}
