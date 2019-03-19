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
