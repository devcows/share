package lib

import (
	"os"

	"github.com/jhoonb/archivex"
)

func CompressFile(inputFilePath string, outputFilePath string) error {
	zip := new(archivex.ZipFile)
	zip.Create(outputFilePath)
	if info, err := os.Stat(inputFilePath); err == nil && info.IsDir() {
		zip.AddFile(inputFilePath)
	} else {
		zip.AddAll(inputFilePath, true)
	}
	zip.Close()

	return nil
}
