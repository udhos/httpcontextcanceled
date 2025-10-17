# httpcontextcanceled

An ongoing http request context might get canceled by the http server because the connection with the client was broken.

If the context is propagated by the http handler into a forwarded http request, this cancellation will result in `context canceled` errors in the forwarded http request.

Unfortunately, the http server code in the standard library that cancels the request does not provide the errored connection as cause of the cancellation.

[context-canceled-http](./cmd/context-canceled-http) demonstrates how an http client can cancel a request with a cause.

[context-canceled-chi](./cmd/context-canceled-chi) demonstrates how an http client (curl) can cancel a request just by breaking the connection. However in this case the http server will cancel the request without providing a cause.

The broken connection cancels the request without a cause: `cr.conn.cancelCtx()`

```go
func (cr *connReader) handleReadErrorLocked(_ error) {
	if cr.conn == nil {
		return
	}
	cr.conn.cancelCtx()
	if res := cr.conn.curReq.Load(); res != nil {
		res.closeNotify()
	}
}
```

Source: https://cs.opensource.google/go/go/+/refs/tags/go1.25.3:src/net/http/server.go;l=769

Context cacellation without a cause:

```go
ctx, cancelCtx := context.WithCancel(ctx)
```

Source: https://cs.opensource.google/go/go/+/refs/tags/go1.25.3:src/net/http/server.go;l=2012

# References

https://github.com/golang/go/issues/75939
