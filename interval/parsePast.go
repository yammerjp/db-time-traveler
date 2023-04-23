package interval

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/yammerjp/db-time-traveler/system"
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

func ParseInterval(past string, future string) (*system.Interval, error) {
	pastNum, pastTerm, pastErr := devideNumAndTerm(past)
	futureNum, futureTerm, futureErr := devideNumAndTerm(future)
	if pastErr == nil && futureErr != nil {
		return &system.Interval{IsPast: true, Num: pastNum, Term: pastTerm}, nil
	}
	if pastErr != nil && futureErr == nil {
		return &system.Interval{IsPast: false, Num: futureNum, Term: futureTerm}, nil
	}
	return nil, errors.New("failed to parse interval (past or future)")
}

func devideNumAndTerm(str string) (int, string, error) {
	str = strings.ToUpper(strings.ReplaceAll(str, " ", ""))
	i := 0
	for i < len(str) && unicode.IsDigit(rune(str[i])) {
		i++
	}
	num, err := strconv.Atoi(str[0:i])
	if err != nil {
		return 0, "", err
	}
	term := str[i:]
	if term != "DAY" && term != "WEEK" && term != "MONTH" {
		return num, term, errors.New("past allow DAY or WEEK or MONTH")
	}
	return num, term, nil
}
