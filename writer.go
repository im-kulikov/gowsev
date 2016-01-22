package gowsev

import (
	"golang.org/x/net/websocket"
)

/* wsHandler is called automatically by "golang.org/x/net/websocket". 
 * The id is set to zero in writer to signal that the connection is new. The master
 * will handle id assignment.
 */
func wsHandler(conn *websocket.Conn) {
	writer(0, conn)
}

func writer(id uint64, conn *websocket.Conn) {

	writerMessageChan := make(chan []byte)
	writerCloseChan := make(chan struct{})
	globalNewConnChan <- &evConn{id, conn, writerMessageChan, writerCloseChan}

	for {
		select {
		case message := <- writerMessageChan:
			bytesWritten := 0
			for bytesWritten < len(message) {
				n, err := conn.Write(message[bytesWritten:])
				if err != nil {
					continue
				}
				bytesWritten += n
			}
		case <- writerCloseChan:
			return
		}
	}
}
