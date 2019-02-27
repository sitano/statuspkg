package statuspkg

// Scan visits every item in a cause chain with a predicate f.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// A cause chain is a linked list of errors each of which
// references next one implementing Cause().
//
// e1 -> e2 -> e3 -> e4 -> ...
// f(e1), f(e2), ...
func Scan(err error, f func(error)) {
	type causer interface {
		Cause() error
	}

	for err != nil {
		f(err)

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
}

// Search looks for the first occurrence of an item in a cause
// chain of errors for which the predicate f returns true.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// A cause chain is a linked list of errors each of which
// references next one implementing Cause().
//
// e1 -> e2 -> e3 -> f(e4) == true -> ...
//
// Search returns first error in a cause chain for which
// the predicate returned true or nil otherwise.
func Search(err error, f func(error) bool) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		if f(err) {
			return err
		}

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}

	return nil
}
