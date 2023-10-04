package _core

import (
	"fmt"
	"log"
	"os"
	"sync"

	constants_err "github.com/FranMT-S/Challenge-Go/src/constants/errors"
	constants_log "github.com/FranMT-S/Challenge-Go/src/constants/logs"
	"github.com/FranMT-S/Challenge-Go/src/core/bulker"
	"github.com/FranMT-S/Challenge-Go/src/core/parser"
	Helpers "github.com/FranMT-S/Challenge-Go/src/helpers"
	_logs "github.com/FranMT-S/Challenge-Go/src/logs"
	"github.com/FranMT-S/Challenge-Go/src/model"
)

const (
	MAXCONCURRENTALLOWED = 30
	MAXPAGINATION        = 5000
)

/*
Indexer - Contains methods for files that are registered in the database

Parameters:

  - Parser: object that transforms data to emails, must implement IParserMail.

  - Bulker: makes the request to upload the content to the database

  - Pagination: will help divide the number of requests to reduce the load.
    If no paging is assigned (paging = 0), it defaults to 1000, max = 5000.
*/
type Indexer struct {
	Parser     parser.IParserMail
	Bulker     bulker.IBulker
	Pagination int // max:5000, min:1.
}

// Asynchronous methods

// index all files in the directory of the specified path
//
// parameters:
//
//   - path: path of directory to index.
//
//   - maxConcurrent: number of files to be indexed at the same time.
//     max = 30, min = 1.
//
//     If maxConcurrent is greater or less, it will be automatically assigned in min or max.
func (indexer Indexer) StartAsync(path string, maxConcurrent int) {

	pathCh := make(chan string)
	mutex := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	if maxConcurrent < 0 {
		maxConcurrent = 1
	} else if maxConcurrent > MAXCONCURRENTALLOWED {
		maxConcurrent = MAXCONCURRENTALLOWED
	}

	// Initialize Worker fro Thead Pool
	for i := 0; i < maxConcurrent; i++ {
		wg.Add(1)
		go indexer.workerAsync(pathCh, mutex, wg, i+1)
	}

	if err := Helpers.ListAllFilesQuoteChannel(path, pathCh); err != nil {
		log.Println(err)
	}

	wg.Wait()
}

func (indexer Indexer) workerAsync(pathCh chan string, mutex *sync.Mutex, wg *sync.WaitGroup, id int) {
	defer wg.Done()

	var mails []*model.Mail
	NumRequest := 0

	for path := range pathCh {
		file, err := os.Open(path)
		if err != nil {
			_logs.LogSVG(
				constants_log.FILE_NAME_ERROR_INDEXER,
				constants_log.OPERATION_PARSER,
				constants_log.ERROR_OPEN_FILE+": "+path,
				err,
			)

			file.Close()
			continue
		}

		fmt.Printf("\rWorker %v parsing: %v", id, path)

		parsedMail, err := indexer.Parser.Parse(file)
		if err != nil {
			_logs.LogSVG(
				constants_log.FILE_NAME_ERROR_INDEXER,
				constants_log.OPERATION_PARSER,
				constants_log.ERROR_PARSER_FAILED+": "+path,
				err,
			)
			file.Close()
			continue
		}

		mails = append(mails, parsedMail)

		file.Close()

		if len(mails) == indexer.Pagination {
			NumRequest++
			indexer.safeRequest(mails, mutex, id, NumRequest)
			mails = nil
		}
	}

	// upload remaining emails
	if len(mails) > 0 {
		indexer.safeRequest(mails, mutex, id, NumRequest)
		mails = nil

	}
}

func (indexer Indexer) safeRequest(mails []*model.Mail, mutex *sync.Mutex, id int, NumRequest int) {
	mutex.Lock()
	indexer.Bulker.Bulk(mails)
	fmt.Println("---------------------------")
	fmt.Printf("--Worker %v, Request %v Completed--------\n", id, NumRequest)
	fmt.Println("---------------------------")
	mutex.Unlock()
}

// synchronous methods

// StartFromArray index all files from a array of string with path of files
func (indexer Indexer) StartFromArray(FilePaths []string) {

	if indexer.Pagination <= 0 {
		indexer.Pagination = 1000
	}

	if indexer.Pagination > MAXPAGINATION {
		indexer.Pagination = MAXPAGINATION
	}

	if indexer.Parser == nil {
		panic(constants_err.ERROR_PARSER_UNINITIALIZED)
	}

	if indexer.Bulker == nil {
		panic(constants_err.ERROR_BULKER_UNINITIALIZED)
	}

	if len(FilePaths) == 0 {
		panic(constants_err.ERROR_ARRAY_EMPTY)
	}

	indexer.work(FilePaths)
}

// Start index all files in the directory of the specified path
func (indexer Indexer) Start(path string) {

	if indexer.Pagination <= 0 {
		indexer.Pagination = 1000
	}

	if indexer.Pagination > MAXPAGINATION {
		indexer.Pagination = MAXPAGINATION
	}

	if indexer.Parser == nil {
		panic(constants_err.ERROR_PARSER_UNINITIALIZED)
	}

	if indexer.Bulker == nil {
		panic(constants_err.ERROR_BULKER_UNINITIALIZED)
	}

	FilePaths, err := Helpers.ListAllFilesQuoteBasic(path)
	if err == nil {
		indexer.work(FilePaths)
	}
}

func (indexer Indexer) work(FilePaths []string) {

	paths := make([]string, len(FilePaths))
	count := (len(FilePaths) / indexer.Pagination)

	// fix for pagination
	if (len(FilePaths) % indexer.Pagination) != 0 {
		count++
	}

	for i := 0; i < count; i++ {
		start := i * indexer.Pagination
		end := (i + 1) * indexer.Pagination

		// end must be smaller size of the array
		// start must be less than the length of the array
		// the remainder when dividing by the pagination should not be 0
		if end > len(FilePaths) && len(FilePaths)%indexer.Pagination != 0 {
			paths = FilePaths[start:]
		} else if start < len(FilePaths) {
			paths = FilePaths[start:end]
		}

		if len(paths) > 0 {
			var mails []*model.Mail
			for j := 0; j < len(paths); j++ {

				file, err := os.Open(paths[j])
				if err != nil {
					_logs.LogSVG(
						constants_log.FILE_NAME_ERROR_INDEXER,
						constants_log.OPERATION_PARSER,
						constants_log.ERROR_OPEN_FILE+": "+paths[j],
						err,
					)
					file.Close()
					continue
				}

				fmt.Print("\rParsing: " + paths[j])
				parsedMail, err := indexer.Parser.Parse(file)
				if err != nil {
					_logs.LogSVG(
						constants_log.FILE_NAME_ERROR_INDEXER,
						constants_log.OPERATION_PARSER,
						constants_log.ERROR_PARSER_FAILED+": "+paths[j],
						err,
					)
					file.Close()
					continue
				}

				mails = append(mails, parsedMail)

				file.Close()
			}

			indexer.Bulker.Bulk(mails)

			mails = nil
			fmt.Println("---------------------------")
			fmt.Printf("---------Request %v Completed--------\n", i+1)
			fmt.Println("---------------------------")
		}
	}
}
