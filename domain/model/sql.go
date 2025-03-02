package model

import (
	"errors"
	"strings"
)

// ddl is Data Definition Language List
type ddl []string

// dml is Data Manipulation Language List
type dml []string

// tcl is Transaction Control Language List
type tcl []string

// dcl is Data Control Language List
type dcl []string

// SQL is a struct that contains SQL statements.
type SQL struct {
	query string
	ddl   ddl
	dml   dml
	tcl   tcl
	dcl   dcl
}

// NewSQL return *SQL
func NewSQL(q string) (*SQL, error) {
	if q == "" {
		return nil, errors.New("query is empty")
	}

	return &SQL{
		query: q,
		ddl:   []string{"CREATE", "DROP", "ALTER", "REINDEX"},
		dml:   []string{"SELECT", "INSERT", "UPDATE", "DELETE", "EXPLAIN", "WITH"},
		tcl:   []string{"BEGIN", "COMMIT", "ROLLBACK", "SAVEPOINT", "RELEASE"},
		dcl:   []string{"GRANT", "REVOKE"},
	}, nil
}

// IsDDL return wherther string is ddl or not.
func (sql *SQL) IsDDL() bool {
	return contains(sql.ddl, strings.ToUpper(sql.firstWord()))
}

// IsDML return wherther string is dml or not.
func (sql *SQL) IsDML() bool {
	return contains(sql.dml, strings.ToUpper(sql.firstWord()))
}

// IsTCL return wherther string is tcl or not.
func (sql *SQL) IsTCL() bool {
	return contains(sql.tcl, strings.ToUpper(sql.firstWord()))
}

// IsDCL returns true if the given string represents a Data Control Language (DCL) statement.
func (sql *SQL) IsDCL() bool {
	return contains(sql.dcl, strings.ToUpper(sql.firstWord()))
}

// IsSelect returns true if the given string represents a SELECT statement.
func (sql *SQL) IsSelect() bool {
	return strings.ToUpper(sql.firstWord()) == "SELECT"
}

// IsInsert returns true if the given string represents an INSERT statement.
func (sql *SQL) IsInsert() bool {
	return strings.ToUpper(sql.firstWord()) == "INSERT"
}

// IsUpdate returns true if the given string represents an UPDATE statement.
func (sql *SQL) IsUpdate() bool {
	return strings.ToUpper(sql.firstWord()) == "UPDATE"
}

// IsDelete returns true if the given string represents a DELETE statement.
func (sql *SQL) IsDelete() bool {
	return strings.ToUpper(sql.firstWord()) == "DELETE"
}

// IsExplain returns true if the given string represents an EXPLAIN statement.
func (sql *SQL) IsExplain() bool {
	return strings.ToUpper(sql.firstWord()) == "EXPLAIN"
}

// IsWith returns true if the given string represents a WITH statement.
func (sql *SQL) IsWith() bool {
	return strings.ToUpper(sql.firstWord()) == "WITH"
}

// firstWord returns the first word of the given string.
func (sql *SQL) firstWord() string {
	return strings.Split(trimWordGaps(sql.query), " ")[0]
}

// String returns the query string.
func (sql *SQL) String() string {
	return sql.query
}

// contains checks if a string exists in a slice of strings.
func contains(list []string, v string) bool {
	for _, s := range list {
		if v == s {
			return true
		}
	}
	return false
}

// trimWordGaps trims extra spaces between words in a string.
func trimWordGaps(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
