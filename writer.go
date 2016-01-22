package gowsev

import (
	"golang.org/x/net/websocket"
)

func writer(conn *websocket.Conn) {

	writerMessageChan := make(chan []byte)
	writerCloseChan := make(chan struct{})
	globalNewConnChan <- &evConn{0, conn, writerMessageChan, writerCloseChan}

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
