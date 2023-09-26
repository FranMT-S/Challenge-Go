package Helpers

import (
	"errors"
	"log"
	"os"
)

func CreateDirectoryLogIfNotExist(name string) {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(name, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}
