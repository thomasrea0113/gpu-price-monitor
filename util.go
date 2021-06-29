package monitor

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/url"
)

func writeAndLog(w io.Writer, format string, a ...interface{}) {
	log.Printf(format, a...)
	fmt.Fprintf(w, format, a...)
}

func NewFalse() *bool {
	b := false
	return &b
}

// ensures all arguments to sprinf are properly escaped
func Uprintf(format string, vv ...string) string {
	vCopy := make([]interface{}, len(vv))
	for i, v := range vv {
		vCopy[i] = url.PathEscape(v)
	}
	return fmt.Sprintf(format, vCopy...)
}

// ensures all arguments to sprinf are properly escaped
func Hprintf(format string, vv ...string) string {
	vCopy := make([]interface{}, len(vv))
	for i, v := range vv {
		vCopy[i] = html.EscapeString(v)
	}
	return fmt.Sprintf(format, vCopy...)
}
