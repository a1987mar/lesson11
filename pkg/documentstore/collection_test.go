package documentstore

import (
	"reflect"
	"testing"
)

func TestCollection_Put(t *testing.T) {
	docs := map[string]Document{}
	docs["id1"] = GetTestDocuments(GetTestFields("123", DocumentFieldTypeNumber))
	docs["id2"] = GetTestDocuments(GetTestFields("id2", DocumentFieldTypeString))
	docs["id3"] = GetTestDocuments(GetTestFields("true", DocumentFieldTypeBool))
	type fields struct {
		documents map[string]Document
		config    CollectionConfig
	}
	type args struct {
		doc Document
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid document with correct primary key",
			fields: fields{
				documents: docs,
				config:    CollectionConfig{PrimaryKey: "id"},
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{
						"id": {
							Type:  DocumentFieldTypeString,
							Value: "123",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing primary key field",
			fields: fields{
				documents: map[string]Document{},
				config:    CollectionConfig{PrimaryKey: "id"},
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{
						"name": {
							Type:  DocumentFieldTypeString,
							Value: "Test",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "primary key not string type",
			fields: fields{
				documents: map[string]Document{},
				config:    CollectionConfig{PrimaryKey: "id"},
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{
						"id": {
							Type:  DocumentFieldTypeNumber,
							Value: 123,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Collection{
				documents: tt.fields.documents,
				config:    tt.fields.config,
			}
			if err := s.Put(tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollection_Get(t *testing.T) {
	docs := map[string]Document{}
	docs["id1"] = GetTestDocuments(GetTestFields("123", DocumentFieldTypeNumber))
	docs["id2"] = GetTestDocuments(GetTestFields("id2", DocumentFieldTypeString))
	docs["id3"] = GetTestDocuments(GetTestFields("true", DocumentFieldTypeBool))

	type fields struct {
		Documents map[string]Document
		Config    CollectionConfig
	}
	type args struct {
		doc Document
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "valid document with correct primary key",
			fields: fields{
				Documents: docs,
				Config:    CollectionConfig{PrimaryKey: "id"},
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{
						"id": {
							Type:  DocumentFieldTypeString,
							Value: "123",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing primary key field",
			fields: fields{
				Documents: map[string]Document{},
				Config:    CollectionConfig{PrimaryKey: "id"},
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{
						"name": {
							Type:  DocumentFieldTypeString,
							Value: "Test User",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "primary key field is not string type",
			fields: fields{
				Documents: map[string]Document{},
				Config:    CollectionConfig{PrimaryKey: "id"},
			},
			args: args{
				doc: Document{
					Fields: map[string]DocumentField{
						"id": {
							Type:  DocumentFieldTypeNumber, // not string!
							Value: 123,
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Collection{
				documents: tt.fields.Documents,
				config:    tt.fields.Config,
			}
			if err := s.Put(tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCollection_Delete(t *testing.T) {
	docs := map[string]Document{}
	docs["id"] = GetTestDocuments(GetTestFields("id", DocumentFieldTypeNumber))
	docs["id2"] = GetTestDocuments(GetTestFields("id2", DocumentFieldTypeString))
	docs["id3"] = GetTestDocuments(GetTestFields("id3", DocumentFieldTypeString))

	type fields struct {
		Documents map[string]Document
		Config    CollectionConfig
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "DELETE document with correct primary key",
			fields: fields{
				Documents: docs,
				Config:    CollectionConfig{PrimaryKey: "id"},
			},
			args: args{
				key: "id",
			},
			want: true,
		},
		{
			name: "missing primary key field",
			fields: fields{
				Documents: docs,
				Config:    CollectionConfig{PrimaryKey: "id4"},
			},
			args: args{
				"id4",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Collection{
				documents: tt.fields.Documents,
				config:    tt.fields.Config,
			}
			if got := s.Delete(tt.args.key); got != tt.want {
				t.Errorf("Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCollection_List(t *testing.T) {
	docs := map[string]Document{}
	docs["id1"] = GetTestDocuments(GetTestFields("id1", DocumentFieldTypeNumber))
	docs["id2"] = GetTestDocuments(GetTestFields("id2", DocumentFieldTypeString))
	docs["id3"] = GetTestDocuments(GetTestFields("id3", DocumentFieldTypeString))
	tests := []struct {
		name string
		s    *Collection
		want []Document
	}{
		{name: "valid List with correct documnet",
			s: &Collection{
				documents: map[string]Document{},
				config: CollectionConfig{
					PrimaryKey: "id-1",
				},
			},
			want: []Document{},
		},

		{
			name: "collection with one document",
			s: &Collection{
				documents: map[string]Document{
					"doc1": {
						Fields: map[string]DocumentField{
							"id": {Type: DocumentFieldTypeString, Value: "123"},
						},
					},
				},
				config: CollectionConfig{
					PrimaryKey: "id",
				},
			},
			want: []Document{
				{
					Fields: map[string]DocumentField{
						"id": {Type: DocumentFieldTypeString, Value: "123"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.List(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Collection.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func GetTestFields(v string, t DocumentFieldType) map[string]DocumentField {
	docs := make(map[string]DocumentField)
	docs[v] = DocumentField{
		Type:  DocumentFieldTypeNumber,
		Value: v,
	}
	return docs
}

func GetTestDocuments(fields ...map[string]DocumentField) Document {
	document := Document{
		Fields: make(map[string]DocumentField),
	}
	for _, field := range fields {
		for k, v := range field {
			document.Fields[k] = v
		}
	}
	return document
}
