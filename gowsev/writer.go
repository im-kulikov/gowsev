package gowsev

import "github.com/gorilla/websocket"

/* The writer keeps a websocket connection and waits for messages from the master.
 * The writer writes the message from the master to its connection.
 */
func writer(conn *websocket.Conn, writerCommandChan chan writerCommand) {
	for {
		select {
		case writerCommand := <-writerCommandChan:
			if writerCommand.close {
				conn.Close()
				return
			} else {
				conn.WriteMessage(writerCommand.messageType, writerCommand.data)
			}
		}
	}
}
