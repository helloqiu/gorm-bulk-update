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
