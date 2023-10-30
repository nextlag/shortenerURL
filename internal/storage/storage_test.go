package storage

import (
	"sync"
	"testing"
)

func TestInMemoryStorage_Load(t *testing.T) {
	type fields struct {
		Data  map[string]string
		Mutex sync.Mutex
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &InMemoryStorage{
				Data:  tt.fields.Data,
				Mutex: tt.fields.Mutex,
			}
			if err := s.Load(tt.args.filename); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
