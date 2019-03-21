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
		var i = 0
		var err = errors.New("e")

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
		var cause = errors.New("e")
		var wrap error = &caused{cause}
		var err error = &caused{wrap}

		var i = 0
		var list = []error{err, wrap, cause}

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
		if Search(nil, nil) != nil {
			t.Error("search over nil must do nothing")
		}
	})

	t.Run("simple", func(t *testing.T) {
		var err = errors.New("e")
		if x := Search(err, func(e error) bool {
			return e == err
		}); x != err {
			t.Error("unexpected result")
		}
	})

	t.Run("chained", func(t *testing.T) {
		var cause = errors.New("e")
		var wrap error = &caused{cause}
		var err error = &caused{wrap}
		var list = []error{err, wrap, cause}
		for _, x := range list {
			if y := Search(err, func(e error) bool {
				return e == x
			}); y != x {
				t.Error("unexpected result")
			}
		}
	})
}
