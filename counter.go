package specfile

import (
	"reflect"

	"github.com/marguerite/go-stdlib/slice"
)

// Counter if a line syntactically ends
type Counter struct {
	NextLineConcats int64 // means "\\"
	CurlyBrackets   int64 // means "{"
	SquareBrackets  int64 // means "["
	Parentheses     int64 // means "("
}

// Reset reset the counter
func (c *Counter) Reset() {
	c.NextLineConcats = 0
	c.CurlyBrackets = 0
	c.SquareBrackets = 0
	c.Parentheses = 0
}

// Valid if this is a syntactically valid line
func (c Counter) Valid(options ...string) bool {
	v := reflect.ValueOf(c)
	for i := 0; i < v.NumField(); i++ {
		// can't use "!= 0" check here, because such shell syntax may occur in rpm specfile:
		// case `value` in
		// *)
		//     `command`
		//		 ;;
		// esac
		// there's no opening "(" here, so the parenteses counter will be negative, but it's still valid
		if ok, _ := slice.Contains(options, v.Type().Field(i).Name); !ok && v.Field(i).Int() > 0 {
			return false
		}
	}
	return true
}

// Count count the line break stuff in line and write to counter
func (c *Counter) Count(b []byte) {
	for i := 0; i < len(b); i++ {
		switch b[i] {
		// '\' means to concat next line
		case '\\':
			//len(b)-2 because len(b)-1 is "\n", we take line-break into consideration
			if i == len(b)-1 || i == len(b)-2 {
				c.NextLineConcats = 1
				break
			}
		case '{':
			c.CurlyBrackets++
		case '}':
			c.CurlyBrackets--
		case '(':
			c.Parentheses++
		case ')':
			c.Parentheses--
		case '[':
			c.SquareBrackets++
		case ']':
			c.SquareBrackets--
		}
	}
}
