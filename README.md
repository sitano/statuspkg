# statuspkg

Package statuspkg provides compatibility of the gRPC status
errors with a Dave's Chaney `https://github.com/pkg/errors`.

It can wrap errors with a status code, extract status object
from the cause chain and have fully compatible API.

`go get github.com/sitano/statuspkg`

The package allows wrapping gRPC status responses, overrides
and developing of the smart wrappers.

    Business logic:
        cause := status.Error(codes.Internal, "something bad happen")
        ...
        err := errors.Wrap(cause, "specific business level info")

    gRPC server handler:
        if statuspkg.Code(err) == codes.Internal {
            // do something about it
        }

    gRPC server interceptor:
        log.Info(err)
        return statuspkg.Convert(err)

It also provides implementation of all gRPC status package
methods which now supports cause chains of errors.

## License

BSD-2-Clause
