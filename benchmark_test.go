package main

import (
	"testing"
)

func BenchmarkListAllFilesRecursive(b *testing.B) {
	// path := "src/db/maildir"
	path := "src/db/maildir"
	for i := 0; i < b.N; i++ {

		ListAllFilesRecursive(path)
	}
}
func BenchmarkListListAllFilesQuoteBasic(b *testing.B) {
	// path := "src/db/maildir"
	path := "src/db/maildir"
	for i := 0; i < b.N; i++ {

		ListAllFilesQuoteBasic(path)
	}
}

// func BenchmarkListListAllFilesSync(b *testing.B) {
// 	// path := "src/db/maildir"
// 	path := "src/db/maildir"
// 	for i := 0; i < b.N; i++ {

// 		ListAllFilesAsync(path, 50)
// 	}
// }
