package errors

import (
	"strings"
)

type MultiErr struct {
	errs []error
}

func (me *MultiErr) IsNil() bool {
	return len(me.errs) == 0
}

func (me *MultiErr) Err() error {
	if len(me.errs) > 0 {
		return me
	}
	return nil
}

func (me *MultiErr) Append(err error) {
	if err != nil {
		me.errs = append(me.errs, err)
	}
}

func (me *MultiErr) Error() string {
	if len(me.errs) > 0 {
		sb := &strings.Builder{}
		sb.WriteString("multierr: [")
		for i, err := range me.errs {
			if i > 0 {
				sb.WriteString("  ")
			}
			sb.WriteString("``")
			sb.WriteString(err.Error())
			sb.WriteString("``")
		}
		sb.WriteString("]")
		return sb.String()
	}
	return "<nil>"
}

func (me *MultiErr) String() string {
	return me.Error()
}
