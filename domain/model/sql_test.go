package model

import (
	"testing"
)

func TestContains(t *testing.T) {
	t.Parallel()

	type args struct {
		list []string
		v    string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success to find value",
			args: args{
				list: []string{"a", "b", "c"},
				v:    "b",
			},
			want: true,
		},
		{
			name: "failed to find value",
			args: args{
				list: []string{"a", "b", "c"},
				v:    "d",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := contains(tt.args.list, tt.args.v); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimWordGaps(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success to trim word gaps. delete head and tail spaces",
			args: args{
				s: "  a  b  c  ",
			},
			want: "a b c",
		},
		{
			name: "success to trim word gaps. delete no spaces",
			args: args{
				s: "a b c",
			},
			want: "a b c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := trimWordGaps(tt.args.s); got != tt.want {
				t.Errorf("trimWordGaps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLIsDDL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		sql  *SQL
		want bool
	}{
		{
			name: "is DDL - CREATE",
			sql: func() *SQL {
				sql, _ := NewSQL("CREATE TABLE test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is DDL - DROP",
			sql: func() *SQL {
				sql, _ := NewSQL("DROP TABLE test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is DDL - ALTER",
			sql: func() *SQL {
				sql, _ := NewSQL("ALTER TABLE test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is not DDL - SELECT",
			sql: func() *SQL {
				sql, _ := NewSQL("SELECT * FROM test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.sql.IsDDL(); got != tt.want {
				t.Errorf("SQL.IsDDL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLIsDML(t *testing.T) {
	tests := []struct {
		name string
		sql  *SQL
		want bool
	}{
		{
			name: "is DML - INSERT",
			sql: func() *SQL {
				sql, _ := NewSQL("INSERT INTO test VALUES (1, 'test')") //nolint:errcheck
				return sql
			}(),
			want: true,
		},
		{
			name: "is DML - UPDATE",
			sql: func() *SQL {
				sql, _ := NewSQL("UPDATE test SET name = 'updated' WHERE id = 1") //nolint:errcheck
				return sql
			}(),
			want: true,
		},
		{
			name: "is DML - DELETE",
			sql: func() *SQL {
				sql, _ := NewSQL("DELETE FROM test WHERE id = 1") //nolint:errcheck
				return sql
			}(),
			want: true,
		},
		{
			name: "is not DML - CREATE TABLE",
			sql: func() *SQL {
				sql, _ := NewSQL("CREATE TABLE test (id INT)") //nolint:errcheck
				return sql
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.sql.IsDML(); got != tt.want {
				t.Errorf("SQL.IsDML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLString(t *testing.T) {
	tests := []struct {
		name   string
		sqlStr string
		want   string
	}{
		{
			name:   "SELECT statement",
			sqlStr: "SELECT * FROM table",
			want:   "SELECT * FROM table",
		},
		{
			name:   "INSERT statement",
			sqlStr: "INSERT INTO table VALUES (1, 'test')",
			want:   "INSERT INTO table VALUES (1, 'test')",
		},
		{
			name:   "CREATE TABLE statement",
			sqlStr: "CREATE TABLE test (id INT)",
			want:   "CREATE TABLE test (id INT)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sql, err := NewSQL(tt.sqlStr)
			if err != nil {
				t.Fatalf("Failed to create SQL: %v", err)
			}

			if got := sql.String(); got != tt.want {
				t.Errorf("SQL.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLIsTCL(t *testing.T) {
	tests := []struct {
		name string
		sql  *SQL
		want bool
	}{
		{
			name: "is TCL - BEGIN",
			sql: func() *SQL {
				sql, _ := NewSQL("BEGIN") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is TCL - COMMIT",
			sql: func() *SQL {
				sql, _ := NewSQL("COMMIT") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is TCL - ROLLBACK",
			sql: func() *SQL {
				sql, _ := NewSQL("ROLLBACK") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is not TCL - SELECT",
			sql: func() *SQL {
				sql, _ := NewSQL("SELECT * FROM test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.sql.IsTCL(); got != tt.want {
				t.Errorf("SQL.IsTCL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLIsDCL(t *testing.T) {
	tests := []struct {
		name string
		sql  *SQL
		want bool
	}{
		{
			name: "is DCL - GRANT",
			sql: func() *SQL {
				sql, _ := NewSQL("GRANT SELECT ON table TO user") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is DCL - REVOKE",
			sql: func() *SQL {
				sql, _ := NewSQL("REVOKE SELECT ON table FROM user") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is not DCL - SELECT",
			sql: func() *SQL {
				sql, _ := NewSQL("SELECT * FROM test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.sql.IsDCL(); got != tt.want {
				t.Errorf("SQL.IsDCL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLIsSelect(t *testing.T) {
	tests := []struct {
		name string
		sql  *SQL
		want bool
	}{
		{
			name: "is SELECT - simple",
			sql: func() *SQL {
				sql, _ := NewSQL("SELECT * FROM test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is SELECT - with WHERE",
			sql: func() *SQL {
				sql, _ := NewSQL("SELECT id, name FROM test WHERE id = 1") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is not SELECT - INSERT",
			sql: func() *SQL {
				sql, _ := NewSQL("INSERT INTO test VALUES (1, 'test')") //nolint:errcheck // ignore error
				return sql
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.sql.IsSelect(); got != tt.want {
				t.Errorf("SQL.IsSelect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLIsInsert(t *testing.T) {
	tests := []struct {
		name string
		sql  *SQL
		want bool
	}{
		{
			name: "is INSERT - values",
			sql: func() *SQL {
				sql, _ := NewSQL("INSERT INTO test VALUES (1, 'test')") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is INSERT - with columns",
			sql: func() *SQL {
				sql, _ := NewSQL("INSERT INTO test (id, name) VALUES (1, 'test')") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is not INSERT - SELECT",
			sql: func() *SQL {
				sql, _ := NewSQL("SELECT * FROM test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.sql.IsInsert(); got != tt.want {
				t.Errorf("SQL.IsInsert() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLIsUpdate(t *testing.T) {
	tests := []struct {
		name string
		sql  *SQL
		want bool
	}{
		{
			name: "is UPDATE - simple",
			sql: func() *SQL {
				sql, _ := NewSQL("UPDATE test SET name = 'updated'") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is UPDATE - with WHERE",
			sql: func() *SQL {
				sql, _ := NewSQL("UPDATE test SET name = 'updated' WHERE id = 1") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is not UPDATE - SELECT",
			sql: func() *SQL {
				sql, _ := NewSQL("SELECT * FROM test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.sql.IsUpdate(); got != tt.want {
				t.Errorf("SQL.IsUpdate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLIsDelete(t *testing.T) {
	tests := []struct {
		name string
		sql  *SQL
		want bool
	}{
		{
			name: "is DELETE - simple",
			sql: func() *SQL {
				sql, _ := NewSQL("DELETE FROM test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is DELETE - with WHERE",
			sql: func() *SQL {
				sql, _ := NewSQL("DELETE FROM test WHERE id = 1") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is not DELETE - SELECT",
			sql: func() *SQL {
				sql, _ := NewSQL("SELECT * FROM test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.sql.IsDelete(); got != tt.want {
				t.Errorf("SQL.IsDelete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLIsExplain(t *testing.T) {
	tests := []struct {
		name string
		sql  *SQL
		want bool
	}{
		{
			name: "is EXPLAIN - SELECT",
			sql: func() *SQL {
				sql, _ := NewSQL("EXPLAIN SELECT * FROM test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is EXPLAIN - with ANALYZE",
			sql: func() *SQL {
				sql, _ := NewSQL("EXPLAIN ANALYZE SELECT * FROM test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is not EXPLAIN - SELECT",
			sql: func() *SQL {
				sql, _ := NewSQL("SELECT * FROM test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.sql.IsExplain(); got != tt.want {
				t.Errorf("SQL.IsExplain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSQLIsWith(t *testing.T) {
	tests := []struct {
		name string
		sql  *SQL
		want bool
	}{
		{
			name: "is WITH - simple CTE",
			sql: func() *SQL {
				sql, _ := NewSQL("WITH cte AS (SELECT id FROM test) SELECT * FROM cte") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is WITH - multiple CTEs",
			sql: func() *SQL {
				sql, _ := NewSQL("WITH cte1 AS (SELECT id FROM test), cte2 AS (SELECT name FROM test) SELECT * FROM cte1 JOIN cte2") //nolint:errcheck // ignore error
				return sql
			}(),
			want: true,
		},
		{
			name: "is not WITH - SELECT",
			sql: func() *SQL {
				sql, _ := NewSQL("SELECT * FROM test") //nolint:errcheck // ignore error
				return sql
			}(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.sql.IsWith(); got != tt.want {
				t.Errorf("SQL.IsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}
