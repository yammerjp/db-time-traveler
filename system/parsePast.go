package system

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func ParsePast(str string) (string, error) {
	str = strings.ToUpper(strings.ReplaceAll(str, " ", ""))
	i := 0
	for i < len(str) && unicode.IsDigit(rune(str[i])) {
		i++
	}
	num, err := strconv.Atoi(str[0:i])
	if err != nil {
		return "", err
	}
	term := str[i:]
	if term != "DAY" && term != "WEEK" && term != "MONTH" {
		return "", errors.New("past allow DAY or WEEK or MONTH")
	}
	return strconv.Itoa(num) + " " + term, nil
}
