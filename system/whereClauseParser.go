package system

import (
	"errors"
	"strings"
)

type WhereClause struct {
	LeftHand  string
	RightHand string
	Operator  string
}

func ParseWhereClause(str string) (*WhereClause, error) {
	hands := strings.Split(str, "=")
	if len(hands) != 2 {
		return nil, errors.New("Where Clause must Includes only 1 equal operator")
	}
	return &WhereClause{
		LeftHand:  strings.TrimSpace(hands[0]),
		RightHand: strings.TrimSpace(hands[1]),
		Operator:  "=",
	}, nil
}
