package gowsev

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

/* serverHandler is registered with the http server.
 * A new accetd connection is sent to serverHandler in a new goroutine.
 * ServerHandler upgrades to the web socket protocol and tells the master.
 * The new goroutine stays alive and becomes a writer.
 */

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("%s\n", err)
		return
	}

	writerCommandChan := make(chan writerCommand)
	writerInitChan <- writerInit{conn, writerCommandChan}
	writer(conn, writerCommandChan)
}
