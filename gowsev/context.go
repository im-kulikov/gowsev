package gowsev

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

func MakeContext(handler *Handler) Context {
	writerInitChan = make(chan writerInit)
	return Context{handler, 0, make(map[uint64]econn), make(chan readerResult), time.Minute}
}

func (context *Context) GetTimeout() time.Duration {
	return context.timeout
}

func (context *Context) SetTimeout(timeout time.Duration) {
	context.timeout = timeout
}

func (context *Context) ListenAndServe(port string) {
	var httpServer http.Server
	httpServer.Addr = ":" + port
	httpServer.Handler = http.HandlerFunc(serverHandler)

	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Printf("Listen error: %s", err)
		}
	}()
}

func (context *Context) EventLoopIteration() {
	handler := *context.handler

	select {
	case writerInit := <-writerInitChan:
		context.idCounter++
		econn := econn{context.idCounter, writerInit.conn, writerInit.writerCommandChan}
		context.econnMap[context.idCounter] = econn
		go reader(econn.id, econn.conn, context.readerResultChan)
		handler.ConnAccepted(context, econn.id)
	case readerResult := <-context.readerResultChan:
		econn := context.econnMap[readerResult.id]
		if readerResult.err != nil {
			econn.writerCommandChan <- writerCommand{true, 0, nil}
			handler.ConnClosed(context, econn.id)
		} else {
			handler.MessageReceived(context, econn.id, readerResult.data)
		}
	case <-time.After(context.timeout):
		handler.EventLoopTimeout(context)
	}
}

func (context *Context) EventLoop() {
	for {
		context.EventLoopIteration()
	}
}

func (context *Context) Write(id uint64, message []byte) error {
	econn, ok := context.econnMap[id]
	if ok {
		econn.writerCommandChan <- writerCommand{false, 2, message}
		return nil
	} else {
		return errors.New("id not found")
	}
}

func (context *Context) Close(id uint64) {
	econn, ok := context.econnMap[id]
	if ok {
		econn.conn.Close() // The reader will notice and send a message and then the writer will be closed by a message from the master.
	}
}

func (context *Context) AddConn(conn *websocket.Conn) uint64 {
	context.idCounter++
	id := context.idCounter
	writerCommandChan := make(chan writerCommand)
	econn := econn{id, conn, writerCommandChan}
	context.econnMap[id] = econn
	go writer(conn, writerCommandChan)
	go reader(id, conn, context.readerResultChan)
	return id
}
