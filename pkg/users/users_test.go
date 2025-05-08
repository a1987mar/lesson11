package users

import (
	"lesson4/pkg/documentstore"
	"reflect"
	"testing"
)

func BenchmarkCreateUser(b *testing.B) {
	st := documentstore.NewStore()
	s := NewService(st)
	_, _ = st.CreateCollection("name", "id")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		d1 := documentstore.Document{
			Fields: GetTestFields("u1", "Andrii", documentstore.DocumentFieldTypeString),
		}
		s.CreateUser("id1", "user1", &d1)
	}
}

func TestService_GetUser(t *testing.T) {
	doc1 := []documentstore.Document{
		{Fields: GetTestFields("u1", "Andrii", documentstore.DocumentFieldTypeString)},
		{Fields: GetTestFields("u2", "Lubov", documentstore.DocumentFieldTypeString)},
		{Fields: GetTestFields("u4", "Taras", documentstore.DocumentFieldTypeString)},
		{Fields: GetTestFields("u3", "Roman", documentstore.DocumentFieldTypeString)},
	}
	store := documentstore.NewStore()
	_, colection := store.CreateCollection("name", "id") // створює колекцію з первинним ключем "id"
	for _, doc := range doc1 {
		colection.Put(doc)
	}
	collect, _ := store.GetCollection("name")
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "GET user",
			s: &Service{
				coll: collect,
			},
			args: args{
				userID: "u3",
			},
			want: &User{
				ID:   "u3",
				Name: "Roman",
			},
			wantErr: false,
		},
		{
			name: "GET user",
			s: &Service{
				coll: collect,
			},
			args: args{
				userID: "u5",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.GetUser(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func GetTestFields(v, name string, t documentstore.DocumentFieldType) map[string]documentstore.DocumentField {
	docs := make(map[string]documentstore.DocumentField)
	docs["id"] = documentstore.DocumentField{
		Type:  documentstore.DocumentFieldTypeString,
		Value: v,
	}
	docs["name"] = documentstore.DocumentField{
		Type:  documentstore.DocumentFieldTypeString,
		Value: name,
	}
	return docs
}

func GetTestDocuments(fields ...map[string]documentstore.DocumentField) documentstore.Document {
	document := documentstore.Document{
		Fields: make(map[string]documentstore.DocumentField),
	}
	for _, field := range fields {
		for k, v := range field {
			document.Fields[k] = v
		}
	}
	return document
}
func TestService_DeleteUser(t *testing.T) {

	doc1 := []documentstore.Document{
		{Fields: GetTestFields("u1", "Andrii", documentstore.DocumentFieldTypeString)},
		{Fields: GetTestFields("u4", "Taras", documentstore.DocumentFieldTypeString)},
		{Fields: GetTestFields("u3", "Roman", documentstore.DocumentFieldTypeString)},
	}
	store := documentstore.NewStore()
	_, colection := store.CreateCollection("name", "id") // створює колекцію з первинним ключем "id"
	for _, doc := range doc1 {
		colection.Put(doc)
	}
	collect, _ := store.GetCollection("name")

	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		s       *Service
		args    args
		wantErr bool
	}{

		{
			name: "delete user for key",
			s:    &Service{coll: collect},
			args: args{
				userID: "u1",
			},
			wantErr: false,
		},
		{
			name: "not delete user for key",
			s:    &Service{coll: collect},
			args: args{
				userID: "u2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.DeleteUser(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("Service.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_ListUsers(t *testing.T) {
	doc1 := []documentstore.Document{
		{Fields: GetTestFields("u1", "Andrii", documentstore.DocumentFieldTypeString)},
		{Fields: GetTestFields("u4", "Taras", documentstore.DocumentFieldTypeString)},
		{Fields: GetTestFields("u3", "Roman", documentstore.DocumentFieldTypeString)},
	}
	store := documentstore.NewStore()
	_, colection := store.CreateCollection("name", "id")
	for _, doc := range doc1 {
		colection.Put(doc)
	}
	collect, _ := store.GetCollection("name")

	storeEmpty := documentstore.NewStore()
	_, colectionEmpty := storeEmpty.CreateCollection("name", "id")
	doc2 := []documentstore.Document{}
	for _, docEmpty := range doc2 {
		colectionEmpty.Put(docEmpty)
	}
	collectEmpty, _ := storeEmpty.GetCollection("name")

	tests := []struct {
		name    string
		s       *Service
		want    []User
		wantErr bool
	}{
		{
			name: "List for users",
			s: &Service{
				coll: collect,
			},
			want: []User{
				{ID: "u1", Name: "Andrii"},
				{ID: "u4", Name: "Taras"},
				{ID: "u3", Name: "Roman"},
			},
			wantErr: false,
		},
		{
			name: "List for empty",
			s: &Service{
				coll: collectEmpty,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.ListUsers()
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.ListUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_CreateUser(t *testing.T) {
	doc1 := documentstore.Document{
		Fields: GetTestFields("u1", "Andrii", documentstore.DocumentFieldTypeString),
	}
	docBad := documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id": {
				Type:  documentstore.DocumentFieldTypeNumber,
				Value: 123,
			},
			"name": {
				Type:  documentstore.DocumentFieldTypeString,
				Value: "badUser",
			},
		},
	}

	type args struct {
		id   string
		name string
		doc  *documentstore.Document
	}
	tests := []struct {
		name    string
		s       func() *Service
		args    args
		want    *User
		wantErr bool
	}{
		{name: "Create new user successfully",
			s: func() *Service {
				store := documentstore.NewStore()
				_, coll := store.CreateCollection("name", "id")
				return &Service{coll: coll}
			},
			args: args{
				id:   "u1",
				name: "Andrii",
				doc:  &doc1},
			want: &User{
				ID:   "u1",
				Name: "Andrii",
			},
			wantErr: false,
		},
		{name: "Create user that already exists",
			s: func() *Service {
				store := documentstore.NewStore()
				_, coll := store.CreateCollection("name", "id")
				coll.Put(doc1)
				return &Service{coll: coll}
			},
			args: args{
				id:   "u1",
				name: "Andrii",
				doc:  &doc1},
			want:    nil,
			wantErr: true,
		},
		{name: "Create user with invalid document field type",
			s: func() *Service {
				store := documentstore.NewStore()
				_, coll := store.CreateCollection("name", "id")
				return &Service{coll: coll}
			},
			args: args{
				id:   "bad",
				name: "badUser",
				doc:  &docBad,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup := tt.s()
			got, err := setup.CreateUser(tt.args.id, tt.args.name, tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.CreateUser() = %v, want %v", got, tt.want)
			}
		})
	}

}
