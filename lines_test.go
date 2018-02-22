package lines_test

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/briansorahan/lines"
)

func TestFrom(t *testing.T) {
	t.Run("count test", func(t *testing.T) {
		var myCount int64

		if expected, got := 0, lines.From(testReader(), lines.Func(func(line string, count int64) error {
			myCount = count
			return nil
		})); expected != got {
			t.Fatalf("expected %d, got %d", expected, got)
		}
		if expected, got := int64(3), myCount; expected != got {
			t.Fatalf("expected %d, got %d", expected, got)
		}
	})
	t.Run("lines test", func(t *testing.T) {
		var lns []string

		if expected, got := 0, lines.From(testReader(), lines.Func(func(line string, count int64) error {
			lns = append(lns, line)
			return nil
		})); expected != got {
			t.Fatalf("expected %d, got %d")
		}
		if expected, got := []string{"foo", "bar", "baz"}, lns; !reflect.DeepEqual(expected, got) {
			t.Fatalf("expected %#v, got %#v", expected, got)
		}
	})
	t.Run("default err code test", func(t *testing.T) {
		code := lines.From(testReader(), lines.Func(func(line string, count int64) error {
			return errors.New("should cause From to return DefaultErrCode")
		}))
		if expected, got := lines.DefaultErrCode, code; expected != got {
			t.Fatalf("expected %d, got %d", expected, got)
		}
	})
	t.Run("custom err code test", func(t *testing.T) {
		code := lines.From(testReader(), lines.Func(func(line string, count int64) error {
			return lines.Error{
				Code: 2,
				Msg:  "custom error code",
			}
		}))
		if expected, got := 2, code; expected != got {
			t.Fatalf("expected %d, got %d", expected, got)
		}
	})
}

func testReader() io.Reader {
	return strings.NewReader(`foo
bar
baz
`)
}
