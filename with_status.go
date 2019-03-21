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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type withStatus struct {
	error

	s *status.Status
}

func (w *withStatus) Error() string              { return w.error.Error() }
func (w *withStatus) Cause() error               { return w.error }
func (w *withStatus) GRPCStatus() *status.Status { return w.s }

func WithStatus(err error, code codes.Code, msg string) error {
	if err == nil {
		if code == codes.OK {
			return nil
		}

		return status.Error(code, msg)
	}

	return &withStatus{
		error: err,
		s:     status.New(code, msg),
	}
}
