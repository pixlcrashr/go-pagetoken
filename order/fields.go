package order

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Field struct {
	Path  string
	Order Order
}

type Fields []Field

func (f *Fields) UnmarshalString(s string) error {
	// thanks to https://github.com/einride/aip-go/blob/master/ordering/orderby.go
	// for einride's order parsing implementation

	// reset slice
	*f = (*f)[:0]

	if s == "" { // fast path for no ordering
		return nil
	}

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '_' && r != ' ' && r != ',' && r != '.' {
			return fmt.Errorf("unmarshal order by '%s': invalid character %s", s, strconv.QuoteRune(r))
		}
	}

	fields := strings.Split(s, ",")
	fs := make([]Field, 0, len(fields))

	for _, field := range fields {
		parts := strings.Fields(field)
		switch len(parts) {
		case 1: // default ordering (ascending)
			fs = append(fs, Field{Path: parts[0]})
		case 2: // specific ordering
			var o Order
			if err := o.UnmarshalString(parts[1]); err != nil {
				return fmt.Errorf("unmarshal order by '%s': %w", s, err)
			}

			fs = append(fs, Field{Path: parts[0], Order: o})
		case 0:
			fallthrough
		default:
			return fmt.Errorf("unmarshal order by '%s': invalid format", s)
		}
	}

	*f = fs
	return nil
}
