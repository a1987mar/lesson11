package documentstore

import (
	"errors"
	"lesson4/pkg/err"
	"os"
	"testing"
)

func BenchmarkNewStore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewStore()
	}
}

func BenchmarkReadDamp(b *testing.B) {
	store := NewStore()
	store.CreateCollection("bench", "id")
	store.DumpToFile("bench")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewStoreFromFile("bench")
		if err != nil {
			b.Fatal(err)
		}
	}
	_ = os.Remove("bench.json")
}

func TestStore_CreateCollection(t *testing.T) {
	type fields struct {
		Collections map[string]*Collection
	}
	type args struct {
		name string
		id   string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     error
		wantCreated bool
	}{
		{
			name: "successfully create collection",
			fields: fields{
				Collections: map[string]*Collection{},
			},
			args: args{
				name: "users",
			},
			wantErr:     nil,
			wantCreated: true,
		},
		{
			name: "collection already exists",
			fields: fields{
				Collections: map[string]*Collection{
					"users": {},
				},
			},
			args: args{
				name: "users",
			},
			wantErr:     err.ErrCollectionAlreadyExists,
			wantCreated: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				collections: tt.fields.Collections,
			}
			err, coll := s.CreateCollection(tt.args.name, tt.args.id)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
			created := coll != nil
			if created != tt.wantCreated {
				t.Errorf("CreateCollection() collection created = %v, want %v", created, tt.wantCreated)
			}
		})
	}
}

func TestStore_GetCollection(t *testing.T) {
	type fields struct {
		Collections map[string]*Collection
	}
	type args struct {
		name string
		id   string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     error
		wantCreated bool
	}{
		{name: "collection exists",
			fields: fields{
				Collections: map[string]*Collection{
					"users": {documents: make(map[string]Document),
						config: CollectionConfig{
							PrimaryKey: "id",
						}},
				},
			}, args: args{
				name: "users",
			},

			wantErr:     nil,
			wantCreated: true,
		},
		{
			name: "collection does not exist",
			fields: fields{
				Collections: map[string]*Collection{},
			},
			args: args{
				name: "orders",
			},
			wantErr:     err.ErrCollectionNotFound,
			wantCreated: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				collections: tt.fields.Collections,
			}
			got, err := s.GetCollection(tt.args.name)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GetCollection() error = %v, wantErr = %v", err, tt.wantErr)
			}
			exists := got != nil
			if exists != tt.wantCreated {
				t.Errorf("GetCollection() got = %v, wantExists = %v", got, tt.wantCreated)
			}
		})
	}
}

func TestStore_DeleteCollection(t *testing.T) {
	type fields struct {
		Collections map[string]*Collection
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "collection exists",
			fields: fields{
				Collections: map[string]*Collection{
					"users": {documents: make(map[string]Document),
						config: CollectionConfig{
							PrimaryKey: "id",
						}},
				},
			}, args: args{
				name: "users",
			},
			want: true,
		},
		{
			name: "collection does not exist",
			fields: fields{
				Collections: map[string]*Collection{},
			},
			args: args{
				name: "orders",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Store{
				collections: tt.fields.Collections,
			}
			if got := s.DeleteCollection(tt.args.name); got != tt.want {
				t.Errorf("DeleteCollection() = %v, want %v", got, tt.want)
			}
		})
	}
}
