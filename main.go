package main

import (
	"fmt"
	"time"

	"github.com/FranMT-S/Challenge-Go/src/constants"
	"github.com/FranMT-S/Challenge-Go/src/core"
	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
)

var path string = "db/maildir"

func main() {
	constants.InitializeVarEnviroment()

	// Registra el tiempo de inicio
	startTime := time.Now()

	// // path := "src/db/maildir/allen-p"

	// listFiles := Helpers.ListAllFilesQuoteBasic(path)
	// listFiles := Helpers.ListAllFilesQuoteBasic(path)[0:20000]

	// // listFiles := []string{"db/maildir/arora-h/sent_items/26"}
	// // listFiles := []string{"db/maildir/allen-p/straw/7"}
	// listFiles := []string{"db/maildir/buy-r/inbox/15"}
	// listFiles := []string{"db/maildir/buy-r/inbox/99"}
	// listFiles := []string{"db/maildir/allen-p/_sent_mail/100"}
	// listFiles := []string{"db/maildir/allen-p/_sent_mail/1002"}

	// // for _, v := range listFiles {
	// // 	fmt.Println(v)
	// // }

	// myDatabase.ZincDatabase().CreateIndex()

	// indexer := core.Indexer{

	// 	Parser:     parser.ParserNormal{},
	// 	Bulker:     bulker.CreateBulkerV1(),
	// 	Pagination: 5000,
	// }

	// indexer := core.Indexer{

	// 	Parser:     parser.NewParserAsyn(50),
	// 	Bulker:     bulker.CreateBulkerV1(),
	// 	Pagination: 5000,
	// }

	// indexer := core.Indexer{

	// 	Parser:     parser.NewParserAsyncSpliter(3),
	// 	Bulker:     bulker.CreateBulkerV1(),
	// 	Pagination: 5000,
	// }

	indexer := core.Indexer{

		Parser:     parser.NewParserAsyncRegex(20),
		Bulker:     bulker.CreateBulkerV2(),
		Pagination: 10,
	}

	// indexer.Start(`db/Testformat`)
	// indexer.StartFromArray(listFiles)
	// indexer.StartAsync(`db/maildir/arora-h`, 5)
	indexer.StartAsync(`db/Testformat`, 5)

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	seconds := duration.Seconds()

	fmt.Printf("El código se ejecutó en %.2f segundos\n", seconds)
}
