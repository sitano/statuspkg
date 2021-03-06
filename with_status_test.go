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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestWithStatus(t *testing.T) {
	t.Run("no wrap for nil good", func(t *testing.T) {
		if WithStatus(nil, codes.OK, "") != nil {
			t.Error("nil with ok must be nil")
		}
	})

	t.Run("nil with no good must return original status", func(t *testing.T) {
		err := WithStatus(nil, 1000, "oh")
		if err == nil {
			t.Error("nil with unknown wrap must be not nil")
		}
		if Code(err) != 1000 {
			t.Error("unexpected code")
		}
	})

	t.Run("status wrap works", func(t *testing.T) {
		cause := errors.New("cause")
		err := WithStatus(cause, 1000, "oh")
		if err == nil {
			t.Error("nil with unknown wrap must be not nil")
		}
		s, ok := FromError(err)
		if !ok {
			t.Error("unexpected result")
		}
		if s.Code() != 1000 {
			t.Error("unexpected code")
		}
		if s.Message() != "oh" {
			t.Error("unexpected desc")
		}
		if Cause(err) != cause {
			t.Error("invalid cause")
		}
	})

	t.Run("wrapped status wrap works", func(t *testing.T) {
		cause := errors.New("cause")
		status := WithStatus(cause, 1000, "oh")
		var err error = &withMessage{cause: status, msg: "wrap"}
		s, ok := FromError(err)
		if !ok {
			t.Error("unexpected result")
		}
		if s.Code() != 1000 {
			t.Error("unexpected code")
		}
		if s.Message() != "oh" {
			t.Error("unexpected desc")
		}
		if Cause(err) != cause {
			t.Error("invalid cause")
		}
	})

	t.Run("status override works", func(t *testing.T) {
		cause := status.Error(codes.InvalidArgument, "x")
		err := WithStatus(cause, codes.FailedPrecondition, "y")
		s, ok := FromError(err)
		if !ok {
			t.Error("unexpected result")
		}
		if s.Code() != codes.FailedPrecondition {
			t.Error("unexpected code")
		}
		if s.Message() != "y" {
			t.Error("unexpected desc")
		}
		if Cause(err) != cause {
			t.Error("invalid cause")
		}
	})
}

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }
func (w *withMessage) Cause() error  { return w.cause }
