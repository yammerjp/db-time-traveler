package system

func BuildStmtToUpdate(targetTable string, columns []string, primaryKeys []string, stmtInWhereIn string, interval Interval) (string, error) {
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

func BuildStmtToSelect(targetTable string, columns []string, primaryKeys []string, stmtInWhereIn string) (string, error) {
	return SelectSource{
		Table: Table{
			targetTable: targetTable,
		},
		columns:       columns,
		primaryKeys:   primaryKeys,
		stmtInWhereIn: stmtInWhereIn,
	}.buildStmtToSelect()
}

func BuildStmtToSelectBeforeAndAfter(targetTable string, columns []string, primaryKeys []string, stmtInWhereIn string, interval Interval) (string, error) {
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

func BuildStmtToSelectColumnNames(targetTable string) string {
	return Table{targetTable: targetTable}.buildStmtToSelectColumnNames()
}

func BuildStmtToSelectColumnNamesDateRelated(targetTable string, ignoreColumnNames []string) (string, error) {
	return Table{targetTable: targetTable}.buildStmtToSelectColumnNamesDateRelated(ignoreColumnNames)
}

func BuildStmtToSelectColumnNamesOfPrimaryKey(targetTable string) (string, error) {
	return Table{targetTable: targetTable}.buildStmtToSelectColumnNamesOfPrimaryKey()
}
