package utils

import (
	"strconv"
	"strings"

	"github.com/filswan/go-swan-lib/utils"
)

func ConvertPrice2AttoFil(price string) string {
	fields := strings.Fields(price)
	if len(fields) < 1 {
		return ""
	}
	if len(fields) < 2 {
		return fields[0]
	}
	priceAttoFil := utils.GetInt64FromStr(fields[0])
	unit := strings.ToUpper(fields[1])
	switch unit {
	case "FIL":
		priceAttoFil = priceAttoFil * 1e18
	case "MILLIFIL":
		priceAttoFil = priceAttoFil * 1e15
	case "MICROFIL":
		priceAttoFil = priceAttoFil * 1e12
	case "NANOFIL":
		priceAttoFil = priceAttoFil * 1e9
	case "PICOFIL":
		priceAttoFil = priceAttoFil * 1e6
	case "FEMTOFIL":
		priceAttoFil = priceAttoFil * 1e3
	case "ATTOFIL":
		priceAttoFil = priceAttoFil
	default:
		priceAttoFil = priceAttoFil
	}

	return strconv.FormatInt(priceAttoFil, 10)
}

func GetPriceFormat(price string) string {
	fields := strings.Fields(price)
	if len(fields) < 1 {
		return ""
	}
	if len(fields) < 2 {
		return fields[0]
	}
	priceAttoFil := int64(1)
	unit := strings.ToUpper(fields[1])
	switch unit {
	case "FIL":
		priceAttoFil = priceAttoFil * 1e18
	case "MILLIFIL":
		priceAttoFil = priceAttoFil * 1e15
	case "MICROFIL":
		priceAttoFil = priceAttoFil * 1e12
	case "NANOFIL":
		priceAttoFil = priceAttoFil * 1e9
	case "PICOFIL":
		priceAttoFil = priceAttoFil * 1e6
	case "FEMTOFIL":
		priceAttoFil = priceAttoFil * 1e3
	case "ATTOFIL":
		priceAttoFil = priceAttoFil
	default:
		priceAttoFil = priceAttoFil
	}

	result := strconv.FormatInt(priceAttoFil, 10)
	result = strings.TrimPrefix(result, "1")

	return result
}
