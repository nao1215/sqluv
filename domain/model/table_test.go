package model

import (
	"testing"
)

func TestTableIsSameHeaderColumnName(t *testing.T) {
	t.Parallel()

	type fields struct {
		Name    string
		Header  Header
		Records []Record
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "table has same header column",
			fields: fields{
				Name:    "table_name",
				Header:  Header{"aaa", "bbb", "ccc", "aa", "aaa"},
				Records: []Record{},
			},
			want: true,
		},
		{
			name: "table does not have same header column",
			fields: fields{
				Name:    "table_name",
				Header:  Header{"aaa", "bbb", "ccc"},
				Records: []Record{},
			},
			want: false,
		},
		{
			name: "table does not have header column",
			fields: fields{
				Name:    "table_name",
				Header:  Header{},
				Records: []Record{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tr := NewTable(
				tt.fields.Name,
				tt.fields.Header,
				tt.fields.Records,
			)
			if got := tr.IsSameHeaderColumnName(); got != tt.want {
				t.Errorf("Table.IsSameHeaderColumnName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableValid(t *testing.T) {
	type fields struct {
		Name    string
		Header  Header
		Records []Record
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Name:   "valid_table",
				Header: Header{"aaa", "bbb", "ccc"},
				Records: []Record{
					{"111", "222", "333"},
					{"444", "555", "666"},
					{"777", "888", "999"},
				},
			},
			wantErr: false,
		},
		{
			name: "table name is empty",
			fields: fields{
				Name:   "",
				Header: Header{"aaa", "bbb", "ccc"},
				Records: []Record{
					{"111", "222", "333"},
					{"444", "555", "666"},
					{"777", "888", "999"},
				},
			},
			wantErr: true,
		},
		{
			name: "header is empty",
			fields: fields{
				Name:   "invalid_table",
				Header: Header{},
				Records: []Record{
					{"111", "222", "333"},
					{"444", "555", "666"},
					{"777", "888", "999"},
				},
			},
			wantErr: true,
		},
		{
			name: "record is empty",
			fields: fields{
				Name:    "invalid_table",
				Header:  Header{"aaa", "bbb", "ccc"},
				Records: []Record{},
			},
			wantErr: true,
		},
		{
			name: "header has same name colomn",
			fields: fields{
				Name:   "valid_table",
				Header: Header{"aaa", "bbb", "aaa"},
				Records: []Record{
					{"111", "222", "333"},
					{"444", "555", "666"},
					{"777", "888", "999"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := NewTable(
				tt.fields.Name,
				tt.fields.Header,
				tt.fields.Records,
			)
			if err := tr.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("Table.Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTableEqual(t *testing.T) {
	t.Parallel()

	type fields struct {
		name    string
		Header  Header
		Records []Record
	}
	type args struct {
		t2 *Table
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "table is equal",
			fields: fields{
				name:   "table_name",
				Header: Header{"aaa", "bbb", "ccc"},
				Records: []Record{
					{"111", "222", "333"},
					{"444", "555", "666"},
					{"777", "888", "999"},
				},
			},
			args: args{
				t2: NewTable(
					"table_name",
					Header{"aaa", "bbb", "ccc"},
					[]Record{
						{"111", "222", "333"},
						{"444", "555", "666"},
						{"777", "888", "999"},
					},
				),
			},
			want: true,
		},
		{
			name: "table is not equal (name)",
			fields: fields{
				name:   "table_name",
				Header: Header{"aaa", "bbb", "ccc"},
				Records: []Record{
					{"111", "222", "333"},
					{"444", "555", "666"},
					{"777", "888", "999"},
				},
			},
			args: args{
				t2: NewTable(
					"table_name2",
					Header{"aaa", "bbb", "ccc"},
					[]Record{
						{"111", "222", "333"},
						{"444", "555", "666"},
						{"777", "888", "999"},
					},
				),
			},
			want: false,
		},
		{
			name: "table is not equal (header)",
			fields: fields{
				name:   "table_name",
				Header: Header{"aaa", "bbb", "ccc"},
				Records: []Record{
					{"111", "222", "333"},
					{"444", "555", "666"},
					{"777", "888", "999"},
				},
			},
			args: args{
				t2: NewTable(
					"table_name",
					Header{"aaa", "bbb", "ccc", "ddd"},
					[]Record{
						{"111", "222", "333"},
						{"444", "555", "666"},
						{"777", "888", "999"},
					},
				),
			},
			want: false,
		},
		{
			name: "table is not equal (record)",
			fields: fields{
				name:   "table_name",
				Header: Header{"aaa", "bbb", "ccc"},
				Records: []Record{
					{"111", "222", "333"},
					{"444", "555", "666"},
					{"777", "888", "999"},
				},
			},
			args: args{
				t2: NewTable(
					"table_name",
					Header{"aaa", "bbb", "ccc"},
					[]Record{
						{"111", "222", "333"},
						{"444", "555", "666"},
						{"777", "888", "999"},
						{"aaa", "bbb", "ccc"},
					},
				),
			},
			want: false,
		},
		{
			name: "table is not equal (record value)",
			fields: fields{
				name:   "table_name",
				Header: Header{"aaa", "bbb", "ccc"},
				Records: []Record{
					{"111", "222", "333"},
					{"444", "555", "666"},
					{"777", "888", "999"},
				},
			},
			args: args{
				t2: NewTable(
					"table_name",
					Header{"aaa", "bbb", "ccc"},
					[]Record{
						{"111", "222", "333"},
						{"444", "555", "666"},
						{"777", "888", "99"},
					},
				),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tr := NewTable(
				tt.fields.name,
				tt.fields.Header,
				tt.fields.Records,
			)
			if got := tr.Equal(tt.args.t2); got != tt.want {
				t.Errorf("Table.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
