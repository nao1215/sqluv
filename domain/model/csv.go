package model

import "strings"

// CSV is csv data with header.
type CSV struct {
	name    string
	header  Header
	records []Record
}

// NewCSV create new CSV.
func NewCSV(
	name string,
	header Header,
	records []Record,
) *CSV {
	return &CSV{
		name:    name,
		header:  header,
		records: records,
	}
}

// ToTable convert CSV to Table.
func (c *CSV) ToTable() *Table {
	return NewTable(
		strings.TrimSuffix(c.name, ".csv"),
		c.header,
		c.records,
	)
}

// Equal compare CSV.
func (c *CSV) Equal(c2 *CSV) bool {
	if c.name != c2.name {
		return false
	}
	if !c.header.Equal(c2.header) {
		return false
	}
	if len(c.records) != len(c2.records) {
		return false
	}
	for i, record := range c.records {
		if !record.Equal(c2.records[i]) {
			return false
		}
	}
	return true
}
