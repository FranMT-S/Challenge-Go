package main

import (
	"fmt"
	"time"

	"github.com/FranMT-S/Challenge-Go/src/constants"
	"github.com/FranMT-S/Challenge-Go/src/core"
	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
	myDatabase "github.com/FranMT-S/Challenge-Go/src/db"
)

var path string = "db/maildir"

func main() {

	constants.InitializeVarEnviroment()
	// Registra el tiempo de inicio
	startTime := time.Now()
	myDatabase.ZincDatabase().CreateIndex()

	// // path := "src/db/maildir/allen-p"

	// listFiles := Helpers.ListAllFilesQuoteBasic(path)
	// listFiles := Helpers.ListAllFilesQuoteBasic(path)[0:20000]

	// // listFiles := []string{"db/maildir/arora-h/sent_items/26"}
	// // listFiles := []string{"db/maildir/allen-p/straw/7"}
	// listFiles := []string{"db/maildir/buy-r/inbox/15"}
	// listFiles := []string{"db/maildir/buy-r/inbox/99"}
	// listFiles := []string{"db/maildir/allen-p/_sent_mail/100"}

	// // for _, v := range listFiles {
	// // 	fmt.Println(v)
	// // }

	// var opt string
	// Taking input from user
	// fmt.Println("ingrese un comando: \nclient \nserver \nquite \ncommand: ")
	// fmt.Scanln(&opt)
	// switch strings.ToLower(opt) {
	// case "client":
	// 	fmt.Println("start client")
	// 	mysocket.Client()
	// case "server":
	// 	fmt.Println("start server")
	// 	mysocket.Server()
	// case "quite":
	// 	fmt.Println("saliendo")
	// default:
	// 	fmt.Println("ingrese un comando: client, server, quite")
	// 	fmt.Scanln(&opt)
	// }

	// listFiles := []string{"db/maildir/lokey-t/calendar/33"}
	indexer := core.Indexer{
		// Parser:     parser.NewParserBasic(),
		// Parser: parser.NewParserAsync(20),
		Parser:     parser.NewParserAsyncRegex(15),
		Bulker:     bulker.CreateBulkerV2(),
		Pagination: 100,
	}

	// // indexer.StartFromArray(listFiles)

	// indexer.Start(`db/TestFormat`)
	// indexer.Start(`db/TestFormat`)
	// indexer.StartAsync(`db/maildir`, 30)
	indexer.StartAsync(`db/Test`, 3)
	// indexer.StartAsync(`db/Test`, 10)

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	seconds := duration.Seconds()

	fmt.Printf("El código se ejecutó en %.2f segundos\n", seconds)
}
