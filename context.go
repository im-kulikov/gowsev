package gowsev

import (
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
	timeout           int // milliseconds
}

func makeContext(handler *Handler) EvContext {
	return EvContext{handler, 0, make(map[uint64]evConn), make(chan evMessage), make(chan uint64), 1000}
}

