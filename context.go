package gowsev

import (
	"time"
	"golang.org/x/net/websocket"
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

type evNewConn struct {
	conn              *websocket.Conn
	writerMessageChan chan []byte
	writerCloseChan   chan struct{}
}

var globalNewConnChan chan *evNewConn

type EvContext struct {
	handler           *Handler
	idCounter         uint64
	conns             map[uint64]evConn
	readerMessageChan chan evMessage
	readerCloseChan   chan uint64
	timeout           time.Duration // int64 nanoseconds
}

func makeContext(handler *Handler) EvContext {
	return EvContext{handler, 0, make(map[uint64]evConn), make(chan evMessage), make(chan uint64), 1000}
}

func (context *EvContext) GetTimeout() time.Duration {
	return context.timeout
}

func (context *EvContext) SetTimeout(timeout time.Duration) {
	context.timeout = timeout
}

func (context *EvContext) AddConn(conn *websocket.Conn) {
	go writer(conn)
}

func (context *EvContext) ListenAndServe(port string) {

}

