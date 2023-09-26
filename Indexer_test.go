package main

import (
	"testing"

	_core "github.com/FranMT-S/Challenge-Go/src/core"
	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
	Helpers "github.com/FranMT-S/Challenge-Go/src/helpers"
)

const pathParserTest string = "db/maildir"

var listFilesParser []string

func listFiles() {
	if listFilesParser == nil {
		listFilesParser, _ = Helpers.ListAllFilesQuoteBasic(pathParserTest)
		listFilesParser = listFilesParser[:20000]
	}
}

func BenchmarkIndexerParserNormal(b *testing.B) {
	listFiles()

	for i := 0; i < b.N; i++ {

		indexer := _core.Indexer{

			Parser:     parser.NewParserBasic(),
			Bulker:     bulker.CreateBulkerV1(),
			Pagination: 5000,
		}

		indexer.StartFromArray(listFilesParser)
	}
}
func BenchmarkIndexerParserAsync(b *testing.B) {
	listFiles()
	for i := 0; i < b.N; i++ {
		indexer := _core.Indexer{
			Parser:     parser.NewParserAsync(50),
			Bulker:     bulker.CreateBulkerV1(),
			Pagination: 5000,
		}

		indexer.StartFromArray(listFilesParser)
	}
}

func BenchmarkIndexerParserAsyncRegex(b *testing.B) {
	listFiles()
	for i := 0; i < b.N; i++ {
		indexer := _core.Indexer{
			Parser:     parser.NewParserAsyncRegex(50),
			Bulker:     bulker.CreateBulkerV2(),
			Pagination: 5000,
		}

		indexer.StartFromArray(listFilesParser)
	}
}
