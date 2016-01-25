package gowsev

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

func MakeContext(handler *Handler) Context {
	writerInitChan = make(chan writerInit)
	return Context{handler, 0, make(map[uint64]conn), make(chan readerResult), time.Minute}
}

func (context *Context) GetTimeout() time.Duration {
	return context.timeout
}

func (context *Context) SetTimeout(timeout time.Duration) {
	context.timeout = timeout
}

func (context *Context) AddConn(conn *websocket.Conn) {
	context.idCounter++
	go writer(context.idCounter, conn)
}

func (context *Context) ListenAndServe(port string) {
	var wsServer websocket.Server
	wsServer.Handler = websocket.Handler(wsHandler)

	var httpServer http.Server
	httpServer.Addr = ":" + port
	httpServer.Handler = wsServer

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			fmt.Printf("Listen error: %s", err)
		}
	}()
}

func (context *Context) EventLoopIteration() {

	select {
	case evConn := <-globalNewConnChan:
		acceptedConn := false
		if evConn.id == 0 {
			acceptedConn = true
			context.idCounter++
			evConn.id = context.idCounter
		}
		context.conns[evConn.id] = evConn
		go reader(evConn.id, evConn.conn, context.readerMessageChan, context.readerCloseChan)
		if acceptedConn {
			(*context.handler).ConnAccepted(context, evConn.id)
		}
	case evMessage := <-context.readerMessageChan:
		(*context.handler).MessageReceived(context, evMessage.id, evMessage.message)
	case id := <-context.readerCloseChan:
		evConn, ok := context.conns[id]
		if ok {
			evConn.writerCloseChan <- struct{}{}
			(*context.handler).ConnClosed(context, id)
			delete(context.conns, id)
		}
	case <-time.After(context.timeout):
		(*context.handler).EventLoopTimeout(context)
	}
}

func (context *Context) EventLoop() {
	for {
		context.EventLoopIteration()
	}
}

func (context *Context) Write(id uint64, message []byte) {
	evConn, ok := context.conns[id]
	if ok {
		evConn.writerMessageChan <- message
	}
}

func (context *Context) Close(id uint64) {
	evConn, ok := context.conns[id]
	if ok {
		evConn.conn.Close()
	}
}
