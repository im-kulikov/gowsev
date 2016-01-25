# Gowsev

A Go event loop library for a websocket server. The gowsev event loop runs inside a single goroutine. See the discussion below for the rationale for a single threaded(single goroutine) architecture. 

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

## Usage

To use the server, an event handler must be defined, and a context must be created.
The simples possibe use is

```
var handler gowsev.Handler
// define handler

context := gowsev.MakeContext(&handler)
context.ListenAndServe("9000")
context.EventLoop()
```


### gowsev.Handler

A Gowsev server is an event loop where the user code receives callbacks at certain events.
A user of gowsev must create a handler of interface type gowsev.Handler.

```
type Handler interface {
	ConnAccepted(context *Context, id uint64)
	ConnClosed(context *Context, id uint64)
	EventLoopTimeout(context *Context)
	MessageReceived(context *Context, id uint64, message []byte)
}
```

Websocket connections have ids of type uint64 assigned internally by gowsev.

`ConnAccepted(context *Context, id uint64)` is called every time a new connection is accepted by the server. The user can choose to store the id.

`ConnClosed(context *Context, id uint64)` is called when the connection with the id is closed.

`EventLoopTimeout(context *Context)` is called when the event loops times out.

`MessageReceived(context *Context, id uint64, message []byte)` is called when the conection with id has delivered a message. Messages are of type []byte. Gowsev sends websocket messages of type binary.

### gowsev.Context

The context of type gowsev.Context is an opaque data structure that keeps track of the event loop. The api is

```
func MakeContext(handler *Handler) Context
func (context *Context) GetTimeout() time.Duration
func (context *Context) SetTimeout(timeout time.Duration)
func (context *Context) ListenAndServe(port string)
func (context *Context) ListenAndServeTLS(port string, certFile string, keyFile string)
func (context *Context) EventLoopIteration()
func (context *Context) EventLoop()
func (context *Context) Write(id uint64, message []byte) error
func (context *Context) Close(id uint64)
func (context *Context) AddConn(conn *websocket.Conn) uint64
```

The functions are self explanatory except for a few points. `EventLoopIteration()` performs a single wait for an event. `EventLoop` performs a loop of such iterations. Mostly, `EventLoop()` will be used. `Wrote` writes a message to the connection with id. The message is a binary websocket message. `AddConn` can be used to add external connections to  the event loop. The user app can dial a websocket connection to an external webservice and put the socket into the event loop. The server can listen to multiple ports simultaneously in the same event loop; just call `ListenAndServe(port string)` or `ListenAndServeTLS(port string, certFile string, keyFile string)` several times.


## Rationale for a single threaded event loop





## Architecture of gowsev






## Example

There is an exampke in the directory example/subscribe_publish. The subscribe_publish service
relays messages from users according to subscriptions. The details are described in the file main.go
