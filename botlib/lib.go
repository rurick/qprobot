package botlib

import (
	"strconv"
	"strings"
)

/*Cut * Обрезка строки
 * возвращает новую строку*/
func Cut(text string, limit int) string {
	runes := []rune(text)
	if len(runes) >= limit {
		return string(runes[:limit])
	}
	return text
}

//PriceFormatted -
func PriceFormatted(price int64) string {
	res := ""
	for o := price % 1000; price > 1000; {
		res = " " + strconv.Itoa(int(o)) + res
		price = price / 1000
	}
	return strings.Trim(strconv.Itoa(int(price))+res, " ")
}
