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
	_core "github.com/FranMT-S/Challenge-Go/src/core"
	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
	myDatabase "github.com/FranMT-S/Challenge-Go/src/db"
	Helpers "github.com/FranMT-S/Challenge-Go/src/helpers"
)

func main() {

	constants.InitializeVarEnviroment()
	Helpers.CreateDirectoryLogIfNotExist("profiling")

	// register time
	startTime := time.Now()
	myDatabase.ZincDatabase().CreateIndex()

	path := selectPathToIndex()
	pagination := selectPagination()
	filesConcurrency := selectNumFiles()
	linesConcurrency := selectLinesReadingAtSameTime()

	indexer := _core.Indexer{
		// Parser:     parser.NewParserBasic(),
		// Parser: parser.NewParserAsync(20),
		Parser:     parser.NewParserAsyncRegex(linesConcurrency),
		Bulker:     bulker.CreateBulkerV2(),
		Pagination: pagination,
	}

	// cpu profiling
	f, err := os.Create(fmt.Sprintf("cpu_%v.prof", time.Now().Format("020106_030405")))
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}

	defer f.Close() // error handling omitted for example
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	indexer.StartAsync(path, filesConcurrency)

	endTime := time.Now()
	duration := endTime.Sub(startTime)
	seconds := duration.Seconds()

	fmt.Printf("The code ran in %.2f seconds\n", seconds)

	// Mem profiling
	f2, err := os.Create(fmt.Sprintf("mem_%v.prof", time.Now().Format("020106_030405")))
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer f2.Close() // error handling omitted for example
	runtime.GC()     // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f2); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}

func selectPagination() int {
	var opt string

	for {
		fmt.Printf("Select the pagination that establishes how many files will be uploaded to the database max %v: ", _core.GetMaxmaxPaginationAllowed())
		fmt.Scanln(&opt)
		num, err := strconv.Atoi(opt)

		if err != nil || num < 1 || num > _core.GetMaxmaxPaginationAllowed() {
			fmt.Println("invalid value entered")
			continue
		}
		return num
	}
}

func selectNumFiles() int {
	var opt string

	for {
		fmt.Printf("enter the number of files that will be read at the same time, min:1 , max:%v: ", _core.GetMaxConcurrentAllow())
		fmt.Scanln(&opt)
		num, err := strconv.Atoi(opt)

		if err != nil || num < 1 || num > _core.GetMaxConcurrentAllow() {
			fmt.Println("invalid value entered")
			continue
		}
		return num

	}
}

func selectLinesReadingAtSameTime() int {
	var opt string

	for {
		fmt.Printf("Select the number of lines that will be analyzed at the same time, min:1 , max:%v: ", parser.GetMaxConcurrentLines())
		fmt.Scanln(&opt)
		num, err := strconv.Atoi(opt)

		if err != nil || num < 1 || num > parser.GetMaxConcurrentLines() {
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
