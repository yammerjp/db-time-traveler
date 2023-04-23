package query

func BuildStmtToUpdate(targetTable string, columns []string, primaryKeys []string, stmtInWhereIn string, interval Interval) (UpdateStatement, error) {
	return UpdateSource{
		SelectSource: SelectSource{
			Table: Table{
				TargetTable: targetTable,
			},
			Columns:       columns,
			PrimaryKeys:   primaryKeys,
			StmtInWhereIn: stmtInWhereIn,
		},
		Interval: interval,
	}.buildStmtToUpdate()
}

func BuildStmtToSelectBeforeAndAfter(targetTable string, columns []string, primaryKeys []string, stmtInWhereIn string, interval Interval) (SelectStatement, error) {
	return UpdateSource{
		SelectSource: SelectSource{
			Table: Table{
				TargetTable: targetTable,
			},
			Columns:       columns,
			PrimaryKeys:   primaryKeys,
			StmtInWhereIn: stmtInWhereIn,
		},
		Interval: interval,
	}.buildStmtToSelectBeforeAndAfter()
}

func BuildStmtToSelectColumnNames(targetTable string) SelectStatement {
	return Table{TargetTable: targetTable}.buildStmtToSelectColumnNames()
}

func BuildStmtToSelectColumnNamesDateRelated(targetTable string, ignoreColumnNames []string) (SelectStatement, error) {
	return Table{TargetTable: targetTable}.buildStmtToSelectColumnNamesDateRelated(ignoreColumnNames)
}

func BuildStmtToSelectColumnNamesOfPrimaryKey(targetTable string) (SelectStatement, error) {
	return Table{TargetTable: targetTable}.buildStmtToSelectColumnNamesOfPrimaryKey()
}
