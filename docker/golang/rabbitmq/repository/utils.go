package repository

import (
	"bytes"
	"io"
	"log"
)

// streamToString は、io.ReaderをStringに変換
func streamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(stream)
	log.Print("[ERROR] streamToString()", err)
	return buf.String()
}
