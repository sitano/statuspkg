// Package statuspkg provides primitives which enable
// interoperability in between the gRPC status errors
// and Dave`s Chaney github.com/pkg/errors errors.
//
// This allows wrapping gRPC status responses, overrides
// and developing of the smart wrappers.
//
//     Business logic:
//         cause := status.Error(codes.Internal, "something bad happen")
//         ...
//         err := errors.Wrap(cause, "specific business level info")
//
//     gRPC server handler:
//         if statuspkg.Code(err) == codes.Internal {
//             // do something about it
//         }
//
//     gRPC server interceptor:
//         log.Info(err)
//         return statuspkg.Convert(err)
//
// It also provides implementation of all gRPC status package
// methods which now supports cause chains of errors.
//
package statuspkg
