package gowsev

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)

type evMessage struct {
	id      uint64
	message []byte
}

type evConn struct {
	id                uint64
	conn              *websocket.Conn
	writerMessageChan chan []byte
	writerCloseChan   chan struct{}
}

var globalNewConnChan chan *evConn

type EvContext struct {
	handler           *Handler
	idCounter         uint64
	conns             map[uint64]*evConn
	readerMessageChan chan evMessage
	readerCloseChan   chan uint64
	timeout           time.Duration // int64 nanoseconds
}

func MakeContext(handler *Handler) EvContext {
	globalNewConnChan = make(chan *evConn)
	return EvContext{handler, 0, make(map[uint64]*evConn), make(chan evMessage), make(chan uint64), time.Minute}
}

func (context *EvContext) GetTimeout() time.Duration {
	return context.timeout
}

func (context *EvContext) SetTimeout(timeout time.Duration) {
	context.timeout = timeout
}

func (context *EvContext) AddConn(conn *websocket.Conn) {
	context.idCounter++
	go writer(context.idCounter, conn)
}

func (context *EvContext) ListenAndServe(port string) {
	var wsServer websocket.Server
	wsServer.Handler = websocket.Handler(wsHandler)

	var httpServer http.Server
	httpServer.Addr = ":" + port
	httpServer.Handler = wsServer

	go func () {
		err := httpServer.ListenAndServe()
		if err != nil {
			fmt.Printf("Listen error: %s", err)
		}
	}()
}

func (context *EvContext) EventLoopIteration() {

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

func (context *EvContext) EventLoop() {
	for {
		context.EventLoopIteration()
	}
}

func (context *EvContext) Write(id uint64, message []byte) {
	evConn, ok := context.conns[id]
	if ok {
		evConn.writerMessageChan <- message
	}
}

func (context *EvContext) Close(id uint64) {
	evConn, ok := context.conns[id]
	if ok {
		evConn.conn.Close()
	}
}
