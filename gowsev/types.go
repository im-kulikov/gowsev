package gowsev

import (
	"github.com/gorilla/websocket"
	"time"
)
	
/* readerResult is sent from the reader to the master when the reader has read from the connection */
type readerResult struct {
	id          uint64
	messageType int
	data        []byte
	err         error
}

/* writerCommand is sent by the master to the writer to send it on the writer's connection. */
type writerCommand struct {
	close       bool
	messageType int
	data        []byte
}

/* writerInit is sent by the http handler in a new goroutine to the master */
type writerInit struct {
	conn              *websocket.Conn
	writerCommandChan chan writerCommand
}

/* The gobal channel to send the writerInit on. The master listens to this channel. */
var writerInitChan chan writerInit

/* A conn is a websocket connection and connection specific information needed by the context. */
type conn struct {
	id                uint64
	conn              *websocket.Conn
	writerCommandChan chan writerCommand
}

/* A context is the coordinator of the entire event loop. A context must be created by the user of this package.
 * A context controls the master goroutine and communicates with the readers and writers.
 * A context must be created with a handler of interface Handler. The handler is the created by the user of gowsev.
 */
type Context struct {
	handler          *Handler
	idCounter        uint64
	connMap          map[uint64]conn
	readerResultChan chan readerResult
	timeout          time.Duration // int64 nanosecond, timeout for the event loop.
}
