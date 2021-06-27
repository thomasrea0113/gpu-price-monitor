package monitor

import (
	"fmt"
	"io"
	"log"
)

func writeAndLog(w io.Writer, format string, a ...interface{}) {
	log.Fatalf(format, a...)
	fmt.Fprintf(w, format, a...)
}

func NewFalse() *bool {
	b := false
	return &b
}
