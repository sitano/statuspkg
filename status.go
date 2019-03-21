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

// Package statuspkg provides primitives which enable
// interoperability in between the gRPC status errors
// and Dave`s Chaney github.com/pkg/errors errors.
//
// Package allows wrapping errors with a status code,
// extract status object from the middle of the cause
// chain and have fully compatible API.
//
//    Business logic:
//        err = errors.Wrap(err, "something bad happen")
//        ...
//        return statuspkg.WithStatus(err, codes.Internal, "not succeed")
//
//    gRPC server handler:
//        if statuspkg.Code(err) == codes.Internal {
//            // do something about it
//        }
//
//    gRPC server interceptor:
//        log.Info(err)
//        return statuspkg.Convert(err)
//
// It also provides implementation of all gRPC status package
// methods which now supports cause chains of errors.
//
package statuspkg

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type statuser interface {
	GRPCStatus() *status.Status
}

// New returns a Status representing c and msg.
func New(c codes.Code, msg string) *status.Status {
	return status.New(c, msg)
}

// Newf returns New(c, fmt.Sprintf(format, a...)).
func Newf(c codes.Code, format string, a ...interface{}) *status.Status {
	return New(c, fmt.Sprintf(format, a...))
}

// Error returns an error representing c and msg.  If c is OK, returns nil.
func Error(c codes.Code, msg string) error {
	return New(c, msg).Err()
}

// Errorf returns Error(c, fmt.Sprintf(format, a...)).
func Errorf(c codes.Code, format string, a ...interface{}) error {
	return Error(c, fmt.Sprintf(format, a...))
}

// FromError returns a Status representing err if it was produced from this
// package or has a method `GRPCStatus() *Status`. Otherwise, ok is false and a
// Status is returned with codes.Unknown and the original error message.
func FromError(err error) (s *status.Status, ok bool) {
	if err == nil {
		return status.FromError(nil)
	}
	_ = Search(err, func(t error) bool {
		if t == context.DeadlineExceeded {
			s = New(codes.DeadlineExceeded, t.Error())
			ok = true
		} else if t == context.Canceled {
			s = New(codes.Canceled, t.Error())
			ok = true
		} else if r, k := t.(statuser); k {
			s = r.GRPCStatus()
			ok = true
		}
		return ok
	})
	if ok {
		return s, ok
	}

	return New(codes.Unknown, err.Error()), false
}

// Convert is a convenience function which removes the need to handle the
// boolean return value from FromError.
func Convert(err error) *status.Status {
	s, _ := FromError(err)
	return s
}

// Code returns the Code of the error if it is a Status error, codes.OK if err
// is nil, or codes.Unknown otherwise.
func Code(err error) codes.Code {
	// Don't use FromError to avoid allocation of OK status.
	if err == nil {
		return codes.OK
	}
	var c = codes.Unknown
	_ = Search(err, func(t error) bool {
		if t == context.DeadlineExceeded {
			c = codes.DeadlineExceeded
			return true
		} else if t == context.Canceled {
			c = codes.Canceled
			return true
		} else if r, k := t.(statuser); k {
			c = r.GRPCStatus().Code()
			return true
		}
		return false
	})
	return c
}

// FromContextError converts a context error into a Status.  It returns a
// Status with codes.OK if err is nil, or a Status with codes.Unknown if err is
// non-nil and not a context error.
func FromContextError(err error) *status.Status {
	switch err {
	case nil:
		return New(codes.OK, "")
	case context.DeadlineExceeded:
		return New(codes.DeadlineExceeded, err.Error())
	case context.Canceled:
		return New(codes.Canceled, err.Error())
	default:
		return New(codes.Unknown, err.Error())
	}
}
