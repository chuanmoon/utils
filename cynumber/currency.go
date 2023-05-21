package cynumber

import (
	"math"
	"strconv"
	"strings"
)

type Currency struct {
	Code        string
	Symbol      string
	Precision   int
	Thousand    string
	Decimal     string
	SymbolAfter bool
}

func (accounting *Currency) formatFloat64Money(amount float64) string {
	return accounting.formatInt64Money(accounting.amountToInt64(amount))
}

func (accounting *Currency) formatInt64Money(amount int64) string {
	sign, amountStr := accounting.formatInt64MoneyWithoutSymbol(amount)

	if accounting.SymbolAfter {
		return sign + amountStr + accounting.Symbol
	}

	return sign + accounting.Symbol + amountStr
}

func (accounting *Currency) formatInt64MoneyWithoutSymbol(amount int64) (string, string) {
	sign, text1List, text2 := accounting.splitInt64Money(amount)
	if text2 == "" {
		return sign, strings.Join(text1List, accounting.Thousand)
	}
	return sign, strings.Join(text1List, accounting.Thousand) + accounting.Decimal + text2
}

func (accounting *Currency) splitInt64Money(amount int64) (string, []string, string) {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	text := strconv.FormatInt(amount, 10)
	letText := len(text)

	text1 := "0"
	text2 := ""
	if accounting.Precision == 0 {
		text1 = text
	} else if letText > accounting.Precision {
		text1 = text[:letText-accounting.Precision]
		text2 = text[letText-accounting.Precision:]
	} else {
		for len(text) < accounting.Precision {
			text = "0" + text
		}
		text2 = text
	}

	lenText1 := len(text1)
	startIndex := lenText1 % 3
	count := lenText1/3 + 1
	if startIndex == 0 {
		startIndex = 3
		count--
	}

	text1List := make([]string, count)
	lastIndex := 0
	for i := 0; i < count; i++ {
		text1List[i] = text1[lastIndex : i*3+startIndex]
		lastIndex = i*3 + startIndex
	}

	return sign, text1List, text2
}

func (accounting *Currency) amountToInt64(amount float64) int64 {
	switch accounting.Precision {
	case 0:
		return int64(math.Round(amount))
	case 1:
		return int64(math.Round(amount * 10))
	case 2:
		return int64(math.Round(amount * 100))
	case 3:
		return int64(math.Round(amount * 1000))
	case 4:
		return int64(math.Round(amount * 10000))
	}
	return 0
}

func (accounting *Currency) int64AmountToUsd(amount int64, rate float64) float64 {
	switch accounting.Precision {
	case 0:
		return float64(amount) / rate
	case 1:
		return float64(amount) / rate / 10
	case 2:
		return float64(amount) / rate / 100
	case 3:
		return float64(amount) / rate / 1000
	case 4:
		return float64(amount) / rate / 10000
	}
	return 0
}
