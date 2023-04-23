package query

func BuildStmtToUpdate(targetTable string, columns []string, primaryKeys []string, stmtInWhereIn string, interval Interval) (UpdateStatement, error) {
	return UpdateSource{
		SelectSource: SelectSource{
			Table: Table{
				targetTable: targetTable,
			},
			columns:       columns,
			primaryKeys:   primaryKeys,
			stmtInWhereIn: stmtInWhereIn,
		},
		Interval: interval,
	}.buildStmtToUpdate()
}

func BuildStmtToSelectBeforeAndAfter(targetTable string, columns []string, primaryKeys []string, stmtInWhereIn string, interval Interval) (SelectStatement, error) {
	return UpdateSource{
		SelectSource: SelectSource{
			Table: Table{
				targetTable: targetTable,
			},
			columns:       columns,
			primaryKeys:   primaryKeys,
			stmtInWhereIn: stmtInWhereIn,
		},
		Interval: interval,
	}.buildStmtToSelectBeforeAndAfter()
}

func BuildStmtToSelectColumnNames(targetTable string) SelectStatement {
	return Table{targetTable: targetTable}.buildStmtToSelectColumnNames()
}

func BuildStmtToSelectColumnNamesDateRelated(targetTable string, ignoreColumnNames []string) (SelectStatement, error) {
	return Table{targetTable: targetTable}.buildStmtToSelectColumnNamesDateRelated(ignoreColumnNames)
}

func BuildStmtToSelectColumnNamesOfPrimaryKey(targetTable string) (SelectStatement, error) {
	return Table{targetTable: targetTable}.buildStmtToSelectColumnNamesOfPrimaryKey()
}
