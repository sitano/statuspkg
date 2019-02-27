package statuspkg

import (
	"errors"
	"testing"
)

type caused struct {
	error
}

func (e *caused) Error() string { return e.error.Error() }
func (e *caused) Cause() error  { return e.error }

func TestScan(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		Scan(nil, nil)
	})

	t.Run("simple", func(t *testing.T) {
		err := errors.New("e")
		i := 0
		Scan(err, func(e error) {
			if e != err {
				t.Error("unexpected error")
			} else {
				i++
			}
		})
		if i != 1 {
			t.Error("unexpected number of visits")
		}
	})

	t.Run("chained", func(t *testing.T) {
		cause := errors.New("e")
		var wrap error = &caused{cause}
		var err error = &caused{wrap}
		list := []error{err, wrap, cause}
		i := 0
		Scan(err, func(e error) {
			if e != list[i] {
				t.Error("unexpected error")
			} else {
				i++
			}
		})
		if i != len(list) {
			t.Error("unexpected number of visits")
		}
	})
}

func TestSearch(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		Search(nil, nil)
	})

	t.Run("simple", func(t *testing.T) {
		err := errors.New("e")
		if x := Search(err, func(e error) bool {
			return e == err
		}); x != err {
			t.Error("unexpected result")
		}
	})

	t.Run("chained", func(t *testing.T) {
		cause := errors.New("e")
		var wrap error = &caused{cause}
		var err error = &caused{wrap}
		list := []error{err, wrap, cause}
		for _, x := range list {
			if y := Search(err, func(e error) bool {
				return e == x
			}); y != x {
				t.Error("unexpected result")
			}
		}
	})
}
