package infrastructure

import (
	"context"
	"database/sql"
	"strings"

	"github.com/nao1215/sqluv/domain/model"
)

// Query executes query.
func Query(ctx context.Context, tx *sql.Tx, query *model.SQL) (*model.Table, error) {
	rows, err := tx.QueryContext(ctx, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	header, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	if len(header) == 0 {
		return nil, ErrNoRows
	}

	scanDest := make([]any, len(header))
	rawResult := make([][]byte, len(header))
	for i := range header {
		scanDest[i] = &rawResult[i]
	}

	records := []model.Record{}
	for rows.Next() {
		result := make([]string, len(header))
		err := rows.Scan(scanDest...)
		if err != nil {
			return nil, err
		}

		for i, raw := range rawResult {
			result[i] = string(raw)
		}
		records = append(records, result)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return model.NewTable(ExtractTableName(query), header, records), nil
}

// ExtractTableName extract table name from query.
// The query must be "SELECT" or "EXPLAIN" statement.
func ExtractTableName(sql *model.SQL) string {
	if sql.IsWith() {
		return ""
	}

	query := strings.ReplaceAll(sql.String(), "`", "")
	words := strings.Fields(query)
	for i, v := range words {
		if strings.EqualFold(v, "FROM") || strings.EqualFold(v, "from") {
			return words[i+1]
		}
	}
	return ""
}
