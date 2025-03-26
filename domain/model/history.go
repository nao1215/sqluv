package model

import (
	"strconv"
)

// Histories is sqluv history all record.
type Histories []History

// History is sqluv history record.
type History struct {
	// ID is history id. 1 is oldest
	ID int
	// Request is sqluv history record that is user input from sqluv prompt
	Request string
}

// NewHistory create new History.
func NewHistory(id int, request string) History {
	return History{
		ID:      id,
		Request: request,
	}
}

// ToTable convert History to Table.
func (h Histories) ToTable() *Table {
	records := make([]Record, 0, len(h))
	for _, v := range h {
		records = append(records, Record{
			strconv.Itoa(v.ID), v.Request,
		})
	}
	return NewTable("history", []string{"id", "request"}, records)
}

// ToStringList convert history to string list.
func (h Histories) ToStringList() []string {
	histories := make([]string, 0, len(h))
	for _, v := range h {
		histories = append(histories, v.Request)
	}
	return histories
}
