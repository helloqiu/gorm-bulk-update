package gormbulkupdate

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils/tests"
)

var db, _ = gorm.Open(tests.DummyDialector{}, nil)

func checkBuildClauses(t *testing.T, clauses []clause.Interface, result string, vars []interface{}) {
	var (
		buildNames    []string
		buildNamesMap = map[string]bool{}
		user, _       = schema.Parse(&tests.User{}, &sync.Map{}, db.NamingStrategy)
		stmt          = gorm.Statement{DB: db, Table: user.Table, Schema: user, Clauses: map[string]clause.Clause{}}
	)

	for _, c := range clauses {
		if _, ok := buildNamesMap[c.Name()]; !ok {
			buildNames = append(buildNames, c.Name())
			buildNamesMap[c.Name()] = true
		}

		stmt.AddClause(c)
	}

	stmt.Build(buildNames...)

	if strings.TrimSpace(stmt.SQL.String()) != result {
		t.Errorf("SQL expects %v got %v", result, stmt.SQL.String())
	}

	if !reflect.DeepEqual(stmt.Vars, vars) {
		t.Errorf("Vars expects %+v got %v", stmt.Vars, vars)
	}
}

func TestAssignmentColumns(t *testing.T) {
	results := []struct {
		Clauses []clause.Interface
		Result  string
		Vars    []interface{}
	}{
		{[]clause.Interface{
			clause.Update{},
			AssignmentColumns([]string{"gorm", "helloqiu"}),
		},
			"UPDATE `users` SET `gorm`=`tmp`.`gorm`,`helloqiu`=`tmp`.`helloqiu`",
			nil,
		},
	}
	for idx, result := range results {
		t.Run(fmt.Sprintf("case #%v", idx), func(t *testing.T) {
			checkBuildClauses(t, result.Clauses, result.Result, result.Vars)
		})
	}
}

func TestFromValues(t *testing.T) {
	results := []struct {
		Clauses []clause.Interface
		Result  string
		Vars    []interface{}
	}{
		{
			[]clause.Interface{
				clause.Update{},
				AssignmentColumns([]string{"gorm", "helloqiu"}),
				FromValues{
					Values: [][]interface{}{
						{
							"gorm1", "helloqiu1",
						},
						{
							"gorm2", "helloqiu2",
						},
					},
				},
			},
			"UPDATE `users` SET `gorm`=`tmp`.`gorm`,`helloqiu`=`tmp`.`helloqiu` FROM VALUES ((?,?),(?,?))",
			[]interface{}{"gorm1", "helloqiu1", "gorm2", "helloqiu2"},
		},
	}
	for idx, result := range results {
		t.Run(fmt.Sprintf("case #%v", idx), func(t *testing.T) {
			checkBuildClauses(t, result.Clauses, result.Result, result.Vars)
		})
	}
}

func TestAs(t *testing.T) {
	results := []struct {
		Clauses []clause.Interface
		Result  string
		Vars    []interface{}
	}{
		{
			[]clause.Interface{
				clause.Update{},
				AssignmentColumns([]string{"gorm", "helloqiu"}),
				FromValues{
					Values: [][]interface{}{
						{
							"gorm1", "helloqiu1",
						},
						{
							"gorm2", "helloqiu2",
						},
					},
				},
				As{
					Table:   clause.Table{Name: "tmp"},
					Columns: []string{"gorm", "helloqiu"},
				},
			},
			"UPDATE `users` SET `gorm`=`tmp`.`gorm`,`helloqiu`=`tmp`.`helloqiu` FROM VALUES ((?,?),(?,?)) AS `tmp`(`gorm`,`helloqiu`)",
			[]interface{}{"gorm1", "helloqiu1", "gorm2", "helloqiu2"},
		},
	}
	for idx, result := range results {
		t.Run(fmt.Sprintf("case #%v", idx), func(t *testing.T) {
			checkBuildClauses(t, result.Clauses, result.Result, result.Vars)
		})
	}
}

func TestEqTableColumn(t *testing.T) {
	results := []struct {
		Clauses []clause.Interface
		Result  string
		Vars    []interface{}
	}{
		{
			[]clause.Interface{
				clause.Update{},
				AssignmentColumns([]string{"gorm", "helloqiu"}),
				FromValues{
					Values: [][]interface{}{
						{
							"gorm1", "helloqiu1",
						},
						{
							"gorm2", "helloqiu2",
						},
					},
				},
				As{
					Table:   clause.Table{Name: "tmp"},
					Columns: []string{"gorm", "helloqiu"},
				},
				clause.Where{
					Exprs: []clause.Expression{
						EqTableColumn{
							SourceTable:  clause.Table{Name: "users"},
							SourceColumn: clause.Column{Name: "gorm"},
							TargetTable:  clause.Table{Name: "tmp"},
							TargetColumn: clause.Column{Name: "gorm"},
						},
					},
				},
			},
			"UPDATE `users` SET `gorm`=`tmp`.`gorm`,`helloqiu`=`tmp`.`helloqiu` FROM VALUES ((?,?),(?,?)) AS `tmp`(`gorm`,`helloqiu`) WHERE `users`.`gorm`=`tmp`.`gorm`",
			[]interface{}{"gorm1", "helloqiu1", "gorm2", "helloqiu2"},
		},
	}
	for idx, result := range results {
		t.Run(fmt.Sprintf("case #%v", idx), func(t *testing.T) {
			checkBuildClauses(t, result.Clauses, result.Result, result.Vars)
		})
	}
}
