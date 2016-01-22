package gowsev

import (
	"golang.org/x/net/websocket"
)

func reader(id uint64, conn *websocket.Conn, readerMessageChan chan evMessage, readerCloseChan chan uint64) {

	buffer := make([]byte, 2)

	for {
		bytesRead := 0
		allOfFrameRead := false
		for !allOfFrameRead {
			n, err := conn.Read(buffer[bytesRead:])
			if err != nil {
				readerCloseChan <- id
				return
			}
			bytesRead += n
			if bytesRead < len(buffer) {
				allOfFrameRead = true
			} else {
				newBuffer := make([]byte, 2*len(buffer))
				copy(newBuffer, buffer)
				buffer = newBuffer
			}
		}
		readerMessageChan <- evMessage{id, buffer[:bytesRead]}
	}
}
