package system

import (
	"fmt"
)

type DatabaseColumn struct {
	name      string
	value     string
	isNull    bool
	isPrimary bool
}

type DatabaseRow struct {
	columns []DatabaseColumn
}

func (c *DatabaseColumn) ToString() string {
	if c == nil {
		return "nil"
	}
	if c.isNull {
		return fmt.Sprintf("%s: NULL", c.name)
	} else if c.isPrimary {
		return fmt.Sprintf("%s: %s (PRI)", c.name, c.value)
	} else {
		return fmt.Sprintf("%s: %s", c.name, c.value)
	}
}

func (r *DatabaseRow) ToString() string {
	if r == nil {
		return "nil"
	}
	ret := ""
	for _, v := range r.columns {
		ret += v.ToString() + "\n"
	}
	return ret
}
