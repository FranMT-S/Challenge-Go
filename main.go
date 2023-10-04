package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"

	"github.com/FranMT-S/Challenge-Go/src/constants"
	constants_log "github.com/FranMT-S/Challenge-Go/src/constants/logs"
	_core "github.com/FranMT-S/Challenge-Go/src/core"
	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
	myDatabase "github.com/FranMT-S/Challenge-Go/src/db"
	Helpers "github.com/FranMT-S/Challenge-Go/src/helpers"
	_logs "github.com/FranMT-S/chi-zinc-server/src/logs"
)

var _time = time.Now().Format("010206_030405")

func main() {

	constants.InitializeVarEnviroment()
	Helpers.CreateDirectoryLogIfNotExist("profiling")

	// indexer := _core.Indexer{
	// 	// Parser:     parser.NewParserBasic(),
	// 	// Parser: parser.NewParserAsync(20),
	// 	Parser:     parser.NewParserAsyncRegex(15),
	// 	Bulker:     bulker.CreateBulkerV2(),
	// 	Pagination: 2000,
	// }
	// indexer.StartFromArray([]string{"db/Test/bailey-s/sent_items/15"})

	startLocal()
	// socketServerStart()

}

func startLocal() {

	// Configuration
	path := selectPathToIndex()
	pagination := selectPagination()
	filesConcurrency := selectNumFiles()
	linesConcurrency := selectLinesReadingAtSameTime()
	activateProf := activateProfiling()

	// register time
	startTime := time.Now()
	if err := myDatabase.ZincDatabase().CreateIndex(); err != nil {
		log.Panic(err)
	}

	// cpu profiling
	if activateProf {
		cpuProf, err := os.Create(fmt.Sprintf("profiling/cpu_%v.prof", _time))
		if err != nil {
			_logs.LogSVG(
				constants_log.FILE_NAME_ERROR_PROFILING,
				constants_log.OPERATION_PROFILING,
				constants_log.ERROR_CPU_PROFILING_CREATE,
				err,
			)
		}

		defer cpuProf.Close() // error handling omitted for example

		if err := pprof.StartCPUProfile(cpuProf); err != nil {
			_logs.LogSVG(
				constants_log.FILE_NAME_ERROR_PROFILING,
				constants_log.OPERATION_PROFILING,
				constants_log.ERROR_CPU_PROFILING_START,
				err,
			)
		}

		defer pprof.StopCPUProfile()
	}

	indexer := _core.Indexer{
		// Parser:     parser.NewParserBasic(),
		// Parser: parser.NewParserAsync(20),
		Parser:     parser.NewParserAsyncRegex(linesConcurrency),
		Bulker:     bulker.CreateBulkerV2(),
		Pagination: pagination,
	}

	indexer.StartAsync(path, filesConcurrency)

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	seconds := duration.Seconds()
	fmt.Printf("The code ran in %.2f seconds\n", seconds)

	// Mem profiling
	if activateProf {
		memProf, err := os.Create(fmt.Sprintf("profiling/mem_%v.prof", _time))
		if err != nil {
			_logs.LogSVG(
				constants_log.FILE_NAME_ERROR_PROFILING,
				constants_log.OPERATION_PROFILING,
				constants_log.ERROR_MEM_PROFILING_CREATE,
				err,
			)
		}

		defer memProf.Close()

		runtime.GC()
		if err := pprof.WriteHeapProfile(memProf); err != nil {
			_logs.LogSVG(
				constants_log.FILE_NAME_ERROR_PROFILING,
				constants_log.OPERATION_PROFILING,
				constants_log.ERROR_MEM_PROFILING_WRITE,
				err,
			)
		}
	}

}

func selectPagination() int {
	var opt string

	for {
		fmt.Printf("Select the pagination that establishes how many files will be uploaded to the database max %v: ", _core.MAXPAGINATION)
		fmt.Scanln(&opt)
		num, err := strconv.Atoi(opt)

		if err != nil || num < 1 || num > _core.MAXPAGINATION {
			fmt.Println("invalid value entered")
			continue
		}
		return num
	}
}

func selectNumFiles() int {
	var opt string

	for {
		fmt.Printf("enter the number of files that will be read at the same time, min:1 , max:%v: ", _core.MAXCONCURRENTALLOWED)
		fmt.Scanln(&opt)
		num, err := strconv.Atoi(opt)

		if err != nil || num < 1 || num > _core.MAXCONCURRENTALLOWED {
			fmt.Println("invalid value entered")
			continue
		}
		return num

	}
}

func selectLinesReadingAtSameTime() int {
	var opt string

	for {
		fmt.Printf("Select the number of lines that will be analyzed at the same time, min:1 , max:%v: ", parser.MAX_CONCURRENT_LINES)
		fmt.Scanln(&opt)
		num, err := strconv.Atoi(opt)

		if err != nil || num < 1 || num > parser.MAX_CONCURRENT_LINES {
			fmt.Println("invalid value entered")
			continue
		}
		return num

	}

}

func selectPathToIndex() string {
	var opt string

	for {
		fmt.Printf("write the path to index:")
		fmt.Scanln(&opt)

		if _, err := os.Stat(opt); errors.Is(err, os.ErrNotExist) {
			fmt.Println("the path was not found, enter an existing path")
			continue
		}

		return opt
	}
}

func activateProfiling() bool {
	var opt string

	for {
		fmt.Printf("want to enable profiling, write Y or N:")
		fmt.Scanln(&opt)

		if opt == "Y" || opt == "N" {
			return opt == "Y"
		}

		fmt.Println("invalid option entered")
	}
}

// func selectMode() {
// 	var opt string

// 	for {
// 		fmt.Println("input command: \nlocal \nsocket \nexit \ncommand: ")
// 		fmt.Scanln(&opt)
// 		switch strings.ToLower(opt) {
// 		case "local":
// 			startLocal()
// 			return
// 		case "socket":
// 			socketServerStart()
// 			return
// 		case "exit":
// 			fmt.Println("exit")
// 			return
// 		default:
// 			fmt.Println("select a option valid")
// 		}
// 	}

// }

// func socketServerStart() {
// 	var opt string
// 	for {
// 		fmt.Println("input command: \n1.server \n2.client \n3.exit \ncommand: ")
// 		fmt.Scanln(&opt)
// 		switch strings.ToLower(opt) {
// 		case "server":
// 			fmt.Println("start server")
// 			mysocket.Server()
// 			return
// 		case "client":
// 			fmt.Println("start client")
// 			mysocket.Client()
// 			return
// 		case "exit":
// 			fmt.Println("Exit")
// 			return
// 		default:
// 			fmt.Println("select a option valid")
// 			fmt.Scanln(&opt)
// 		}
// 	}

// }
