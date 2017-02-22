package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv" // for converting integers to strings
)

func main() {
	var buf bytes.Buffer
	writer := NewRESPWriter(&buf)
	writer.WriteCommand("GET", "foo")
	fmt.Printf("%", string(buf.Bytes())) // *2\r\n$3\r\nGET\r\n$3\r\nfoo\r\n
}

var (
	arrayPrefixSlice      = []byte{'*'}
	bulkStringPrefixSlice = []byte{'$'}
	lineEndingSlice       = []byte{'\r', '\n'}
)

type RESPWriter struct {
	*bufio.Writer
}

func NewRESPWriter(writer io.Writer) *RESPWriter {
	return &RESPWriter{
		Writer: bufio.NewWriter(writer),
	}
}

func (w *RESPWriter) WriteCommand(args ...string) (err error) {
	// Write the array prefix and the number of arguments in the array.
	w.Write(arrayPrefixSlice)
	w.WriteString(strconv.Itoa(len(args)))
	w.Write(lineEndingSlice)

	// Write a bulk string for each argument.
	for _, arg := range args {
		w.Write(bulkStringPrefixSlice)
		w.WriteString(strconv.Itoa(len(arg)))
		w.Write(lineEndingSlice)
		w.WriteString(arg)
		w.Write(lineEndingSlice)
	}

	return w.Flush()
}
