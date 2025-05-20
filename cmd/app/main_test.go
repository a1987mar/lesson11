package main

import "testing"

func BenchmarkNewStore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		lesson11()
	}
}
