package gormbulkupdate

import (
	"gorm.io/gorm/clause"
)

func AssignmentColumns(values []string) clause.Set {
	assignments := make([]clause.Assignment, len(values))
	for idx, value := range values {
		assignments[idx] = clause.Assignment{Column: clause.Column{Name: value}, Value: clause.Column{Table: "tmp", Name: value}}
	}
	return assignments
}

// FromValues is a clause which represents SQLs like "FROM VALUES ((a, b, c), (d, e, f))"
type FromValues struct {
	Values [][]interface{}
}

// Name from clause name
func (FromValues) Name() string {
	return "FROM VALUES"
}

// Build build from clause
func (fv FromValues) Build(builder clause.Builder) {
	builder.WriteString("FROM VALUES ")
	builder.WriteByte('(')
	for idx, value := range fv.Values {
		if idx > 0 {
			builder.WriteByte(',')
		}

		builder.WriteByte('(')
		builder.AddVar(builder, value...)
		builder.WriteByte(')')
	}
	builder.WriteByte(')')
}

// MergeClause merge values clauses
func (fv FromValues) MergeClause(clause *clause.Clause) {
	clause.Name = ""
	clause.Expression = fv
}

// As is a clause which represents SQLs like "AS xxx_table(col1, col2)"
type As struct {
	Table   clause.Table
	Columns []string
}

// Name from clause name
func (As) Name() string {
	return "AS"
}

// Build build from clause
func (as As) Build(builder clause.Builder) {
	builder.WriteString("AS ")
	builder.WriteQuoted(as.Table)
	builder.WriteByte('(')
	for idx, col := range as.Columns {
		if idx > 0 {
			builder.WriteByte(',')
		}
		builder.WriteQuoted(col)
	}
	builder.WriteByte(')')
}

// MergeClause merge values clauses
func (as As) MergeClause(clause *clause.Clause) {
	clause.Name = ""
	clause.Expression = as
}

// EqTableColumn is a clause which represents SQLs like "table1.col1 = table2.col2"
type EqTableColumn struct {
	SourceTable  clause.Table
	TargetTable  clause.Table
	SourceColumn clause.Column
	TargetColumn clause.Column
}

func (eq EqTableColumn) Build(builder clause.Builder) {
	builder.WriteQuoted(eq.SourceTable)
	builder.WriteByte('.')
	builder.WriteQuoted(eq.SourceColumn)
	builder.WriteByte('=')
	builder.WriteQuoted(eq.TargetTable)
	builder.WriteByte('.')
	builder.WriteQuoted(eq.TargetColumn)
}
