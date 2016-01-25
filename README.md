# Gowsev

A Go event loop library for a websocket server. The gowsev event loop runs inside a single goroutine. See the discussion below for the rationale for a single threaded(single goroutine) architecture. 

## Api

A Gowsev server is an eventloop where the user code receives callbacks at certain events.
A user of gowsev must create a handler of interface type gowsev.Handler.

```
type Handler interface {
	ConnAccepted(context *Context, id uint64)
	ConnClosed(context *Context, id uint64)
	EventLoopTimeout(context *Context)
	MessageReceived(context *Context, id uint64, message []byte)
}
```

Websocket connections have ids of type uint64 assigned internally by gowsev. The context is described below.


`ConnAccepted(context *Context, id uint64)` is called every time a new connection is accepted by the server. The user might want to store the id.













## Rationale for a single threaded event loop







## Installation

Gowsev is a standard Go program and can be installed with

```
go get github.com/morten-krogh/gowsev
go install github.com/morten-krogh/gowsev
```

#### Dependencies

The only dependency besides the Go standard library is the Gorilla websocket library:

```
go get github.com/gorilla/websocket
```


## Example

