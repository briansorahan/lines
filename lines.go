// Package lines helps you write a little less boilerplate for programs that process newline delimited data.
package lines

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// DefaultErrCode is the default exit code when an error occurs.
var DefaultErrCode = 1

// Processor is a thing that can process lines.
type Processor interface {
	Process(line string, count int64) error
}

// Func processes a line of input.
type Func func(string, int64) error

// Process processes a line with the LineFunc.
func (f Func) Process(line string, count int64) error {
	return f(line, count)
}

// Error is an error type that contains an exit code.
type Error struct {
	Code int
	Msg  string
}

// Error returns an error message.
func (e Error) Error() string {
	return e.Msg
}

// From processes lines from the provided reader.
func From(r io.Reader, p Processor) int {
	var (
		br    = bufio.NewReader(r)
		count = int64(1)
	)
	for {
		line, err := br.ReadString(0x0A)
		if err == io.EOF {
			return 0
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return DefaultErrCode
		}
		if err := p.Process(strings.TrimRightFunc(line, isNewline), count); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())

			if e, ok := err.(Error); ok {
				return e.Code
			}
			return DefaultErrCode
		}
		count++
	}
	return 0
}

func isNewline(c rune) bool {
	return c == 0x0A
}
