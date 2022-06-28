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
