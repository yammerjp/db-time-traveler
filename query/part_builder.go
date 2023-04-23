package query

import (
	"fmt"
	"strings"
)

func (i *Interval) buildInterval() string {
	var symbol string
	if i.IsPast {
		symbol = "-"
	} else {
		symbol = "+"
	}
	return fmt.Sprintf(" %s INTERVAL %d %s", symbol, i.Num, i.Term)
}

func (t Table) buildFrom() string {
	return " FROM " + t.TargetTable
}

func (s SelectSource) buildWhereIn() string {
	pks := strings.Join(s.PrimaryKeys, ", ")

	if !strings.HasPrefix(strings.ToUpper(strings.TrimSpace(s.StmtInWhereIn)), "SELECT") {
		return fmt.Sprintf(" WHERE %s IN ( %s )", pks, s.StmtInWhereIn)
	}

	// MySQL does not allow specifying a table with the same name in WHERE IN SELECT during UPDATE, so an alias is applied to avoid this problem.
	return fmt.Sprintf(" WHERE %s IN ( SELECT %s FROM ( %s ) as any )", pks, pks, s.StmtInWhereIn)
}
