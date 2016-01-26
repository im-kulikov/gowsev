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
The simples possible use is

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

`MessageReceived(context *Context, id uint64, message []byte)` is called when the connection with id has delivered a message. Messages are of type []byte. Gowsev sends websocket messages of type binary.

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

The functions are self explanatory except for a few points. `EventLoopIteration()` performs a single wait for an event. `EventLoop` performs a loop of such iterations. Mostly, `EventLoop()` will be used. `Wrote` writes a message to the connection with id. The message is a binary websocket message. `AddConn` can be used to add external connections to  the event loop. The user app can dial a websocket connection to an external web service and put the socket into the event loop. The server can listen to multiple ports simultaneously in the same event loop; just call `ListenAndServe(port string)` or `ListenAndServeTLS(port string, certFile string, keyFile string)` several times.


## Single threaded event loop

#### Architecture of Gowsev

Gowsev uses multiple goroutines. There is a master goroutine that coordinates everything.
The handler callback functions are called in the master goroutine. Each listening socket has its own goroutine. 

When a new connection is accepted, the net/http go system creates a new goroutine for the accepted socket. The new goroutine sends a message to the master goroutine on a global channel made for that purpose. The new goroutine transforms itself into a writer. A writer waits for messages from the master and sends messages to to the socket. The master gives the new connection a unique uint64 id. The master also creates a reader goroutine. The reader signals to the master when a new websocket message has arrived. All of these goroutines block on either incoming network activity or on Go channels. 

From the point of view of the user code, the system looks like a single threaded event loop.

#### Rationale for a single threaded server

An idiomatic Go server would use one goroutine per connection instead of a single threaded event loop. Actually, the lower level Go code uses a callback scheme such as epoll or kqueue and then distributes the connections to new goroutines. Going back to a single threaded event loop seems backwards; the events fan out to many goroutines that then fan in again. And it is somewhat backwards compared to a situation where the lower level Go network code distributes new connections to a single goroutine and send channel messages every time there is a new event. However, there is also an advantage; the http parsing and websocket protocol handling takes place concurrently and in parallel.   

What is the rationale for the single goroutine callbacks. A minor advantage is that database access is automatically serialized. If the responses to websocket messages requires a database operation, the single goroutine can just operate directly on the data, whereas the multi-goroutine system would need to lock the data access or keep a dedicated database goroutine. Such a database goroutine would itself resemble the master goroutine of gowsev.

The major advantage is that websocket applications often need to coordinate the various open connections. The response to a message on a connection is often not just a reply to that connection, but instead a rely to multiple other connections. The subscribe-publish example included in Gowsev is an example. When a publish message arrives, all subscribed connections must receive a reply. In a single threaded event loop, this is easy. The application just keeps a map with information about subscribers and another map with information about connected websockets. 

In the architecture with one goroutine per connection, it is difficult for the connections to find each other. And a goroutine can not just send messages on connections owned by other goroutines. The goroutines must not block either.

The architecture with one connection per goroutine is really only useful for situations where the connections live independently of each other such as a static file web server. The single threaded event loop is more versatile.     

## Example

There is an example in the directory example/subscribe_publish. The subscribe_publish service
relays messages from users according to subscriptions. The details are described in the file main.go
