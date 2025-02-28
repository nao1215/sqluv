package model

import (
	"github.com/nao1215/sqluv/domain"
)

// Table represents database record.
type Table struct {
	// Name is table name.
	name string
	// Header is table header.
	header Header
	// Records is table records.
	records []Record
}

// NewTable create new Table.
func NewTable(
	name string,
	header Header,
	records []Record,
) *Table {
	return &Table{
		name:    name,
		header:  header,
		records: records,
	}
}

// Name return table name.
func (t *Table) Name() string {
	return t.name
}

// Header return table header.
func (t *Table) Header() Header {
	return t.header
}

// Records return table records.
func (t *Table) Records() []Record {
	return t.records
}

// Equal compare Table.
func (t *Table) Equal(t2 *Table) bool {
	if t.Name() != t2.Name() {
		return false
	}
	if !t.header.Equal(t2.header) {
		return false
	}
	if len(t.Records()) != len(t2.Records()) {
		return false
	}
	for i, record := range t.Records() {
		if !record.Equal(t2.Records()[i]) {
			return false
		}
	}
	return true
}

// Valid check the contents of a Table.
func (t *Table) Valid() error {
	if t.IsEmptyName() {
		return domain.ErrEmptyTableName
	}

	if t.IsEmptyHeader() {
		return domain.ErrEmptyHeader
	}

	if t.IsEmptyRecords() {
		return domain.ErrEmptyRecords
	}

	if t.IsSameHeaderColumnName() {
		return domain.ErrSameHeaderColumns
	}

	return nil
}

// IsEmptyName return wherther table name is empty or not
func (t *Table) IsEmptyName() bool {
	return t.name == ""
}

// IsEmptyHeader return wherther table header is empty or not
func (t *Table) IsEmptyHeader() bool {
	return len(t.header) == 0
}

// IsEmptyRecords return wherther table records is empty or not
func (t *Table) IsEmptyRecords() bool {
	return len(t.records) == 0
}

// IsSameHeaderColumnName return whether the table has a header column with the same name
func (t *Table) IsSameHeaderColumnName() bool {
	encountered := map[string]bool{}
	for i := range t.header {
		if !encountered[t.Header()[i]] {
			encountered[t.Header()[i]] = true
			continue
		}
		return true
	}
	return false
}
