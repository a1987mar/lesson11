package main

import (
	"fmt"
	"log/slog"
	"sync"

	"lesson4/pkg/documentstore"
	"lesson4/pkg/users"
)

func main() {
	slog.Info("App started")
	lesson11()
}

func marshalExample() {
	s := &documentstore.MyStruct{X: 15}
	doc, err := documentstore.MarshalDocument(s)
	if err != nil {
		fmt.Printf("failed to marshal document: %+v\n", err)
		return
	}
	fmt.Printf("marshaled document: %+v\n", doc)
}

func unmarshalExample() {
	doc := &documentstore.Document{Fields: map[string]documentstore.DocumentField{}}
	doc.Fields["X"] = documentstore.DocumentField{
		Type:  documentstore.DocumentFieldTypeNumber,
		Value: int(32),
	}

	s := &documentstore.MyStruct{}
	err := documentstore.UnmarshalDocument(doc, s)
	if err != nil {
		fmt.Printf("failed to unmarshal document: %+v\n", err)
		return
	}
	fmt.Printf("unmarshaled document: %+v\n", s)
}

func lesson11() {

	slog.Info("start app")
	doc1 := []documentstore.Document{
		{
			Fields: map[string]documentstore.DocumentField{
				"id":   {Type: documentstore.DocumentFieldTypeString, Value: "u1"},
				"name": {Type: documentstore.DocumentFieldTypeString, Value: "Andrii"},
			},
		},
		{
			Fields: map[string]documentstore.DocumentField{
				"id":   {Type: documentstore.DocumentFieldTypeString, Value: "u2"},
				"name": {Type: documentstore.DocumentFieldTypeString, Value: "Lubov"},
			},
		},
		{
			Fields: map[string]documentstore.DocumentField{
				"id":   {Type: documentstore.DocumentFieldTypeString, Value: "u4"},
				"name": {Type: documentstore.DocumentFieldTypeString, Value: "Taras"},
			},
		},
		{
			Fields: map[string]documentstore.DocumentField{
				"id":   {Type: documentstore.DocumentFieldTypeString, Value: "u3"},
				"name": {Type: documentstore.DocumentFieldTypeString, Value: "Roman"},
			},
		},
	}

	usersCreated := make([]users.User, 0)
	slog.Info("add new store")
	st := documentstore.NewStore()
	ser := users.NewService(st)

	/// Додавання користувача
	wg := sync.WaitGroup{}
	var mu sync.Mutex
	for l := 0; l < 1000; l++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i, doc := range doc1 {
				u1, err := ser.CreateUser(doc.Fields["id"].Value.(string), doc.Fields["name"].Value.(string), &doc)
				if err != nil {
					fmt.Println("", i+1, err)
				} else {
					mu.Lock()
					usersCreated = append(usersCreated, *u1)
					mu.Unlock()
				}
				getU, err := ser.GetUser(doc.Fields["id"].Value.(string))
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("%+v знайдено \n", getU)
				delUser := ser.DeleteUser(doc.Fields["id"].Value.(string))
				fmt.Println(delUser)

			}
		}()
	}
	wg.Wait()

	slog.Info("App done")
}
