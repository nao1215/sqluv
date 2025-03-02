package config

import (
	"testing"
)

func TestNewMemoryDB(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "generate new memory db",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, cleanup, err := NewMemoryDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInMemDB() error = %v, wantErr %v", err, tt.wantErr)
			}
			cleanup()
		})
	}
}
