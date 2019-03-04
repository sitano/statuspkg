# statuspkg [![Travis-CI](https://travis-ci.org/sitano/statuspkg.svg)](https://travis-ci.org/sitano/statuspkg) [![AppVeyor](https://ci.appveyor.com/api/projects/status/b98mptawhudj53ep/branch/master?svg=true)](https://ci.appveyor.com/project/davecheney/errors/branch/master) [![GoDoc](https://godoc.org/github.com/sitano/statuspkg?status.svg)](http://godoc.org/github.com/sitano/statuspkg) [![Report card](https://goreportcard.com/badge/github.com/sitano/statuspkg)](https://goreportcard.com/report/github.com/sitano/statuspkg) [![Sourcegraph](https://sourcegraph.com/github.com/sitano/statuspkg/-/badge.svg)](https://sourcegraph.com/github.com/sitano/statuspkg?badge)

Package statuspkg provides compatibility of the gRPC status
errors with a Dave's Chaney `https://github.com/pkg/errors`.

It can wrap errors with a status code, extract status object
from the middle of the cause chain and have fully compatible API.

`go get github.com/sitano/statuspkg`

The package allows wrapping gRPC status responses, overrides
and developing of the smart wrappers.

    Business logic:
        err = errors.Wrap(err, "something bad happen")
        ...
        return statuspkg.WithStatus(err, codes.Internal, "not succeed")

    gRPC server handler:
        if statuspkg.Code(err) == codes.Internal {
            // do something about it
        }

    gRPC server interceptor:
        log.Info(err)
        return statuspkg.Convert(err)

It also provides implementation of all gRPC status package
methods which now supports cause chains of errors and context
errors by default.

## License

BSD-2-Clause
