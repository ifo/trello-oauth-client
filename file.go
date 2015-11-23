package trelgo

import (
	"bufio"
	"os"
	"strings"
)

func overwriteFile(dest string, info []byte) error {
	// ensure dest's folders exists
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		index := strings.LastIndex(dest, "/")
		if index != -1 {
			errMkdirAll := os.MkdirAll(dest[:index], 0755)
			if errMkdirAll != nil {
				return errMkdirAll
			}
		}
	}

	outFile, errCreate := os.Create(dest)
	if errCreate != nil {
		return errCreate
	}
	defer outFile.Close()

	w := bufio.NewWriter(outFile)
	_, errWrite := w.Write(info)
	w.Flush()
	return errWrite
}
