package gowsev

import "github.com/gorilla/websocket"

/* reader reads input from the web socket connection and sends it to the master.
 * reader runs in its own goroutine.
 * If the connection closes, the reader goroutine terminates.
 */
func reader(id uint64, conn *websocket.Conn, readerResultChan chan readerResult) {
	for {
		messageType, data, err := conn.ReadMessage()
		readerResultChan <- readerResult{id, messageType, data, err}
		if err != nil {
			return
		}
	}
}
