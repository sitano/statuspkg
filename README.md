# statuspkg [![Travis-CI](https://travis-ci.org/sitano/statuspkg.svg)](https://travis-ci.org/sitano/statuspkg) [![AppVeyor](https://ci.appveyor.com/api/projects/status/b98mptawhudj53ep/branch/master?svg=true)](https://ci.appveyor.com/project/davecheney/errors/branch/master) [![GoDoc](https://godoc.org/github.com/sitano/statuspkg?status.svg)](http://godoc.org/github.com/sitano/statuspkg) [![Report card](https://goreportcard.com/badge/github.com/sitano/statuspkg)](https://goreportcard.com/report/github.com/sitano/statuspkg) [![Sourcegraph](https://sourcegraph.com/github.com/sitano/statuspkg/-/badge.svg)](https://sourcegraph.com/github.com/sitano/statuspkg?badge)

Package statuspkg provides primitives which enable
interoperability in between the gRPC status errors
and Dave`s Chaney github.com/pkg/errors errors.

Package allows wrapping errors with a status code,
extract status object from the middle of the cause
chain and have fully compatible API.

`go get github.com/sitano/statuspkg`

## Wrap gRPC errors

`Code`, `Convert`, `FromError` natively supports cause
chains and are able to extract gRPC status objects from the
inside of the cause chain:

```go
cause := status.Error(codes.IllegalArgument, "blah")
err := errors.Wrap(cause, "ouch")

if statuspkg.Code(err) == codes.IllegalArgument {}
```

`Code(err) == codes.IllegalArgument` is what this library
was created for.

## Wrap errors with gRPC status

Any error can be assigned a gRPC status:

```go
return WithStatus(err, codes.FailedPrecondition, "failed precondition")

// or

cause := something from spanner
wrap := errors.Wrap(cause, "operation xyz")
return WithStatus(err, spanner.Code(cause), cause.Error())
```

## Extract code from error

Implementation of all status related functions in this package
is able to extract gRPC status from any node of the cause chain:

```go
cause := something from spanner
wrap := errors.Wrap(cause, "operation xyz")
return WithStatus(wrap, spanner.Code(cause), cause.Error())

// or 

cause := something from spanner
wrap := WithStatus(err, spanner.Code(cause), cause.Error())
return errors.Wrap(wrap, "operation xyz")

// there is no difference for Code(), Convert(), FromError()
```

in both cases `Code, Convert, FromError` will return the right
status which comes from the WithStatus node.

## gRPC status override

`WithStatus` can override original gRPC status error:

```go
cause := status.Error(codes.IllegalArgument, "blah")
return WithStatus(cause, codes.Unknown, "whoa")
```

`status.Code` will return `codes.Unknown` for this error.

Thus it makes it possible to wrap gRPC statuses and then change
them:

```go
cause := status.Error(codes.IllegalArgument, "blah")
wrap := errors.Wrap(cause, "something went wrong here")
return WithStatus(wrap, codes.Unknown, "whoa")
```

`status.Code` will return `codes.Unknown` for this error,
and the `errors.Cause` will return original status.

## Multi layers architecture support

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

## How to develop custom wrappers

You can write different wrappers extending standard error context
and provide interface to it with functions that support searching over
a cause chain.

Implement `causer` interface as its made for WithStatus wrapper
and provide static helpers for checking state of wrapped errors.
`Cause` method just returns the next element in the cause chain.

```go
type causer interface {
    Cause() error
}
```

In example retryable errors or errors with hidden meta context
may be implement easily with the help of cause chains implemented
in this package.

## Search through the cause chain

```go
Search(err, func(t error) bool {
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
```

## Scan through every element in the cause chain

```go
Scan(err, func(t error) {
	count ++
})
```

## License

BSD-2-Clause
