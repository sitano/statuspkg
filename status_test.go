/*
 *
 * Copyright 2017 gRPC authors.
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
 * Modifications:
 * - 2019 @john.koepi/@sitano deleted tests for details and that have
 *   external dependencies.
 *
 */

package statuspkg

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorsWithSameParameters(t *testing.T) {
	const description = "some description"
	e1 := Errorf(codes.AlreadyExists, description)
	e2 := Errorf(codes.AlreadyExists, description)
	if e1 == e2 || !reflect.DeepEqual(e1, e2) {
		t.Fatalf("Errors should be equivalent but unique - e1: %v  e2: %v", e1, e2)
	}
}

func TestError(t *testing.T) {
	err := Error(codes.Internal, "test description")
	if got, want := err.Error(), "rpc error: code = Internal desc = test description"; got != want {
		t.Fatalf("err.Error() = %q; want %q", got, want)
	}
	s, _ := FromError(err)
	if got, want := s.Code(), codes.Internal; got != want {
		t.Fatalf("err.Code() = %s; want %s", got, want)
	}
	if got, want := s.Message(), "test description"; got != want {
		t.Fatalf("err.Message() = %s; want %s", got, want)
	}
}

func TestErrorOK(t *testing.T) {
	err := Error(codes.OK, "foo")
	if err != nil {
		t.Fatalf("Error(codes.OK, _) = %p; want nil", err)
	}
}

func TestFromError(t *testing.T) {
	code, message := codes.Internal, "test description"
	err := Error(code, message)
	s, ok := FromError(err)
	if !ok || s.Code() != code || s.Message() != message || s.Err() == nil {
		t.Fatalf("FromError(%v) = %v, %v; want <Code()=%s, Message()=%q, Err()!=nil>, true", err, s, ok, code, message)
	}
}

func TestFromErrorOK(t *testing.T) {
	code, message := codes.OK, ""
	s, ok := FromError(nil)
	if !ok || s.Code() != code || s.Message() != message || s.Err() != nil {
		t.Fatalf("FromError(nil) = %v, %v; want <Code()=%s, Message()=%q, Err=nil>, true", s, ok, code, message)
	}
}

type customError struct {
	Code    codes.Code
	Message string
}

func (c customError) Error() string {
	return fmt.Sprintf("rpc error: code = %s desc = %s", c.Code, c.Message)
}

func (c customError) GRPCStatus() *status.Status {
	return status.New(c.Code, c.Message)
}

func TestFromErrorImplementsInterface(t *testing.T) {
	code, message := codes.Internal, "test description"
	err := customError{
		Code:    code,
		Message: message,
	}
	s, ok := FromError(err)
	if !ok || s.Code() != code || s.Message() != message || s.Err() == nil {
		t.Fatalf("FromError(%v) = %v, %v; want <Code()=%s, Message()=%q, Err()!=nil>, true", err, s, ok, code, message)
	}
}

func TestFromErrorUnknownError(t *testing.T) {
	code, message := codes.Unknown, "unknown error"
	err := errors.New("unknown error")
	s, ok := FromError(err)
	if ok || s.Code() != code || s.Message() != message {
		t.Fatalf("FromError(%v) = %v, %v; want <Code()=%s, Message()=%q>, false", err, s, ok, code, message)
	}
}

func TestConvertKnownError(t *testing.T) {
	code, message := codes.Internal, "test description"
	err := Error(code, message)
	s := Convert(err)
	if s.Code() != code || s.Message() != message {
		t.Fatalf("Convert(%v) = %v; want <Code()=%s, Message()=%q>", err, s, code, message)
	}
}

func TestConvertUnknownError(t *testing.T) {
	code, message := codes.Unknown, "unknown error"
	err := errors.New("unknown error")
	s := Convert(err)
	if s.Code() != code || s.Message() != message {
		t.Fatalf("Convert(%v) = %v; want <Code()=%s, Message()=%q>", err, s, code, message)
	}
}

func TestStatus_WithDetails_Fail(t *testing.T) {
	tests := []*status.Status{
		nil,
		status.FromProto(nil),
		New(codes.OK, ""),
	}
	for _, s := range tests {
		if s, err := s.WithDetails(); err == nil || s != nil {
			t.Fatalf("(%v).WithDetails(%+v) = %v; want nil, non-nil", str(s), s, err)
		}
	}
}

func str(s *status.Status) string {
	if s == nil {
		return "nil"
	}
	return fmt.Sprintf("<Code=%v, Message=%q, Details=%+v>", s.Code(), s.Message(), s.Details())
}

func TestFromContextError(t *testing.T) {
	testCases := []struct {
		in   error
		want *status.Status
	}{
		{in: nil, want: New(codes.OK, "")},
		{in: context.DeadlineExceeded, want: New(codes.DeadlineExceeded, context.DeadlineExceeded.Error())},
		{in: context.Canceled, want: New(codes.Canceled, context.Canceled.Error())},
		{in: errors.New("other"), want: New(codes.Unknown, "other")},
	}
	for _, tc := range testCases {
		got := FromContextError(tc.in)
		if got.Code() != tc.want.Code() || got.Message() != tc.want.Message() {
			t.Errorf("FromContextError(%v) = %v; want %v", tc.in, got, tc.want)
		}
	}
}
