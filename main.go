package main

import (
	"flag"
	"log"
	"runtime/pprof"

	"os"
	"runtime"

	"github.com/FranMT-S/Challenge-Go/src/constants"
	"github.com/FranMT-S/Challenge-Go/src/core"
	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
	// "github.com/FranMT-S/Challenge-Go/src/core"
	// "github.com/FranMT-S/Challenge-Go/src/core/bulker"
	// "github.com/FranMT-S/Challenge-Go/src/core/parser"
)

// var path string = "db/test"

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
var path string = "db/maildir"

func main() {
	constants.InitializeVarEnviroment()

	flag.Parse()

	// ... rest of the program ...

	fe, err := os.Create(`mem.prof`)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer fe.Close() // error handling omitted for example
	runtime.GC()     // get up-to-date statistics
	if err := pprof.WriteHeapProfile(fe); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}

	// Registra el tiempo de inicio
	// startTime := time.Now()

	// // path := "src/db/maildir/allen-p"

	// listFiles := Helpers.ListAllFilesQuoteBasic(path)
	// listFiles := Helpers.ListAllFilesQuoteBasic(path)[0:20000]

	// listFiles := []string{"db/maildir/buy-r/inbox/15"}
	// // listFiles := []string{"db/maildir/arora-h/sent_items/26"}
	// // listFiles := []string{"db/maildir/allen-p/straw/7"}
	listFiles := []string{"db/maildir/allen-p/_sent_mail/100"}

	// // for _, v := range listFiles {
	// // 	fmt.Println(v)
	// // }

	// myDatabase.ZincDatabase().CreateIndex()

	// indexer := core.Indexer{
	// 	FilePaths:  listFiles,
	// 	Parser:     parser.ParserNormal{},
	// 	Bulker:     bulker.CreateBulkerV1(),
	// 	Pagination: 5000,
	// }

	// indexer := core.Indexer{
	// 	FilePaths:  listFiles,
	// 	Parser:     parser.NewParserAsyn(50),
	// 	Bulker:     bulker.CreateBulkerV1(),
	// 	Pagination: 5000,
	// }

	indexer := core.Indexer{
		FilePaths:  listFiles,
		Parser:     parser.NewParserAsyncSpliter(3),
		Bulker:     bulker.CreateBulkerV1(),
		Pagination: 5000,
	}

	indexer.Start()

	// endTime := time.Now()
	// duration := endTime.Sub(startTime)
	// seconds := duration.Seconds()

	// fmt.Printf("El código se ejecutó en %.2f segundos\n", seconds)
}
