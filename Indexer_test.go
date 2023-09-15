package main

import (
	"testing"

	"github.com/FranMT-S/Challenge-Go/src/core"
	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
	Helpers "github.com/FranMT-S/Challenge-Go/src/helpers"
)

const pathParserTest string = "db/maildir"

var listFilesParser []string = Helpers.ListAllFilesQuoteBasic(pathParserTest)[0:20000]

func BenchmarkIndexerParserNormal(b *testing.B) {
	// pathTest := "src/db/maildir"

	for i := 0; i < b.N; i++ {

		indexer := core.Indexer{
			FilePaths:  listFilesParser,
			Parser:     parser.ParserNormal{},
			Bulker:     bulker.CreateBulkerV1(),
			Pagination: 5000,
		}

		indexer.Start()
	}
}
func BenchmarkIndexerParserAsync(b *testing.B) {

	for i := 0; i < b.N; i++ {
		indexer := core.Indexer{
			FilePaths:  listFilesParser,
			Parser:     parser.NewParserAsyn(50),
			Bulker:     bulker.CreateBulkerV1(),
			Pagination: 5000,
		}

		indexer.Start()
	}
}

func BenchmarkIndexerParserAsyncSpliter(b *testing.B) {

	for i := 0; i < b.N; i++ {
		indexer := core.Indexer{
			FilePaths:  listFilesParser,
			Parser:     parser.NewParserAsyncSpliter(50),
			Bulker:     bulker.CreateBulkerV1(),
			Pagination: 5000,
		}

		indexer.Start()
	}
}
