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
	writerMessageChan chan string
	writerCloseChan   chan struct{}
}

var globalNewConnChan chan *evConn

type evContext struct {
	idCounter         uint64
	conns             map[uint64]evConn
	readerMessageChan chan evMessage
	readerCloseChan   chan uint64
}

