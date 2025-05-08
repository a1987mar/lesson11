package documentstore

import "testing"

func BenchmarkDocument(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d1 := Document{Fields: make(map[string]DocumentField)}
		d1.Fields["id-1"] = DocumentField{
			Type:  DocumentFieldTypeString,
			Value: "setup.exe",
		}
	}
}

func BenchmarkMarshalDocument(b *testing.B) {
	s := &MyStruct{X: 15}
	for i := 0; i < b.N; i++ {
		MarshalDocument(s)
	}
}

func BenchmarkUnmarshalDocument(b *testing.B) {
	doc := &Document{Fields: map[string]DocumentField{}}
	doc.Fields["X"] = DocumentField{
		Type:  DocumentFieldTypeNumber,
		Value: int(32),
	}
	s := &MyStruct{}
	for i := 0; i < b.N; i++ {
		UnmarshalDocument(doc, s)
	}
}
