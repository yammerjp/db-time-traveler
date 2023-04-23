package query

type Interval struct {
	IsPast bool
	Num    int
	Term   string
}

type Table struct {
	TargetTable string
}

type SelectSource struct {
	Table
	Columns       []string
	PrimaryKeys   []string
	StmtInWhereIn string
}

type UpdateSource struct {
	SelectSource
	Interval
}
