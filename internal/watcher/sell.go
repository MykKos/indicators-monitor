package watcher

import "github.com/MykKos/stock-sdk/pkg/indicators"

func SellMacd(tokens, tf string, macdLine []*indicators.MACDLine) bool {
	lastMacd := macdLine[len(macdLine)-1]
	prelast := macdLine[len(macdLine)-2]

	prevpts := macdLine[len(macdLine)-PrevPoints-1:]

	if NotMACDSellPos(tf, lastMacd, macdLine) {
		return false
	}

	// если сигнальная линия ниже главной
	// или нет точки в выбраном периоде, в которой сигнальная линия ниже главной
	// т.е. нет пересечения
	// сигнала нет
	if lastMacd.MainLineValue.EMA > lastMacd.SignalLineValue.EMA || allBigger(prevpts) {
		return false
	}

	histAge := HistAge(lastMacd.HistogramValue, prelast.HistogramValue)
	// если гистограма растущая бычья, или стареющая медвежья — сигнала нет
	if histAge == GrowthGreenHist || histAge == OldRedHist {
		return false
	}
	return true
}

func NotMACDSellPos(tf string, lastMacd *indicators.MACDLine, macdLine []*indicators.MACDLine) bool {
	// получаем минимальный MACD за период
	mline := MinMacd(macdLine)
	var cond float64
	// получаем значение условия
	switch tf {
	case FiveMinuteTF:
		// FiveMinPercent процентов от минимального (на текущий момент 50%)
		// если мин = -1 -> cond = -0.5
		cond = mline.MainLineValue.EMA * FiveMinPercent
	default:
		// GeneralPercent процентов от минимального (на текущий момент 40%)
		// если мин = -1 -> cond = -0.4
		cond = mline.MainLineValue.EMA * GeneralPercent
	}
	// возвращаем текущее положение MACD
	// если оно ниже сравниваемого значения, то сигнала на продажу нет
	return lastMacd.MainLineValue.EMA < cond
}

// поиск минимального значения MACD
func MinMacd(macdLine []*indicators.MACDLine) *indicators.MACDLine {
	min := &indicators.MACDLine{
		MainLineValue: &indicators.EmaPoint{},
	}
	for _, ln := range macdLine {
		if ln.MainLineValue.EMA < min.MainLineValue.EMA {
			min = ln
		}
	}

	return min
}

func allBigger(comp []*indicators.MACDLine) bool {
	for _, vl := range comp {
		// если есть точка, в которой главная линия выше сигнальной
		// то есть пересечение, сигнал есть
		if vl.MainLineValue.EMA > vl.SignalLineValue.EMA {
			return false
		}
	}
	return true
}
