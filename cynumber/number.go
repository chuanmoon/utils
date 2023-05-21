package cynumber

import (
	"strconv"
	"strings"
)

// USD  	US$  	US$5.00     	. 点号  	, 逗号  	2
// GBP  	£    	GBP£5.00    	. 点号  	, 逗号  	2
// AUD  	AU$  	AU$5.00     	. 点号  	, 逗号  	2
// EUR  	€    	1.115,00€   	, 逗号  	. 点号  	2
// CAD  	CA$  	CA$5.00     	. 点号  	, 逗号  	2
// MXN  	$MXN 	$MXN1,205.00	. 点号  	, 逗号  	2
// JPY  	¥    	¥5,000      	. 点号  	, 逗号  	0
// BRL  	R$   	R$1.115,95  	, 逗号  	. 点号  	2
// ZAR  	R    	R500.00     	. 点号  	, 逗号  	2
var (
	USD = &Currency{Code: "USD", Symbol: "US$", Precision: 2, Thousand: ",", Decimal: "."}
	GBP = &Currency{Code: "GBP", Symbol: "£", Precision: 2, Thousand: ",", Decimal: "."}
	AUD = &Currency{Code: "AUD", Symbol: "AU$", Precision: 2, Thousand: ",", Decimal: "."}
	EUR = &Currency{Code: "EUR", Symbol: "€", Precision: 2, Thousand: ".", Decimal: ",", SymbolAfter: true}
	CAD = &Currency{Code: "CAD", Symbol: "CA$", Precision: 2, Thousand: ",", Decimal: "."}
	MXN = &Currency{Code: "MXN", Symbol: "$MXN", Precision: 2, Thousand: ",", Decimal: "."}
	JPY = &Currency{Code: "JPY", Symbol: "¥", Precision: 0, Thousand: ",", Decimal: "."}
	BRL = &Currency{Code: "BRL", Symbol: "R$", Precision: 2, Thousand: ".", Decimal: ","}
	ZAR = &Currency{Code: "ZAR", Symbol: "R", Precision: 2, Thousand: ",", Decimal: "."}

	Currencies = []*Currency{
		USD,
		GBP,
		AUD,
		EUR,
		CAD,
		MXN,
		JPY,
		BRL,
		ZAR,
	}
)

func GetFormatText(currencyCode, layout string) string {
	accounting := AccountingByCurrencyCode(currencyCode)
	if accounting.SymbolAfter {
		return layout + accounting.Symbol
	}
	return accounting.Symbol + layout
}

// FormatFloat 两位有效数字
func FormatFloat(price float64) string {
	return strconv.FormatFloat(price, 'f', 2, 64)
}

// FormatMoney usd转化为其他币种，并且格式化
func FormatMoneyFromUsd(usdAmount float64, currencyCode string, rate float64) string {
	return AccountingByCurrencyCode(currencyCode).formatFloat64Money(usdAmount * rate)
}

// Int64MoneyFromUsd usd转化为其他币种，并且变为int64
func Int64MoneyFromUsd(usdAmount float64, currencyCode string, rate float64) int64 {
	return AccountingByCurrencyCode(currencyCode).amountToInt64(usdAmount * rate)
}

// Int64MoneyFromAmount 币种变为int64，不涉及汇率
func Int64MoneyFromAmount(amount float64, currencyCode string) int64 {
	return AccountingByCurrencyCode(currencyCode).amountToInt64(amount)
}

// FormatInt64Money 格式化int64的金额
func FormatInt64Money(amount int64, currencyCode string) string {
	return AccountingByCurrencyCode(currencyCode).formatInt64Money(amount)
}

// Int64MoneyToUsd int64金额转化为美元
func Int64MoneyToUsd(amount int64, currencyCode string, rate float64) float64 {
	return AccountingByCurrencyCode(currencyCode).int64AmountToUsd(amount, rate)
}

// FormatInt64MoneyWithoutSymbol 格式化int64的金额
func FormatInt64MoneyWithoutSymbol(amount int64, currencyCode string) string {
	sign, amountStr := AccountingByCurrencyCode(currencyCode).formatInt64MoneyWithoutSymbol(amount)
	return sign + amountStr
}

// FormatInt64MoneyStandard 格式化int64的金额,符号在前,小数点为".",千分位为","
func FormatInt64MoneyStandard(amount int64, currencyCode string) string {
	sign, text1List, text2 := AccountingByCurrencyCode(currencyCode).splitInt64Money(amount)
	if text2 == "" {
		return sign + strings.Join(text1List, "")
	}
	return sign + strings.Join(text1List, "") + "." + text2
}

func AccountingByCurrencyCode(currencyCode string) *Currency {
	switch currencyCode {
	case "USD":
		return USD
	case "GBP":
		return GBP
	case "AUD":
		return AUD
	case "CAD":
		return CAD
	case "MXN":
		return MXN
	case "ZAR":
		return ZAR
	case "JPY":
		return JPY
	case "BRL":
		return BRL
	case "EUR":
		return EUR
	}

	return USD
}
