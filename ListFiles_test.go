package main

import (
	"testing"

	Helpers "github.com/FranMT-S/Challenge-Go/src/helpers"
)

const pathTest string = "db/maildir"

func BenchmarkListFilesRecursive(b *testing.B) {
	// pathTest := "src/db/maildir"

	for i := 0; i < b.N; i++ {

		Helpers.ListAllFilesRecursive(pathTest)
	}
}
func BenchmarkListFilesQuoteBasic(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Helpers.ListAllFilesQuoteBasic(pathTest)
	}
}

func BenchmarkListFilesIterative(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Helpers.ListAllFilesIterative(pathTest)
	}
}
func BenchmarkListFilesWalker(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Helpers.ListAllFilesWalkDir(pathTest)
	}
}

// func BenchmarkListListAllFilesSync(b *testing.B) {
// 	// path := "src/db/maildir"
// 	path := "src/db/maildir"
// 	for i := 0; i < b.N; i++ {

// 		ListAllFilesAsync(path, 50)
// 	}
// }
