/*
 *
 * Copyright 2019 Ivan Prisyazhnyy <john.koepi@gmail.com>.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	var last error
	Scan(err, func(t error) {
		last = t
	})
	return last
}
