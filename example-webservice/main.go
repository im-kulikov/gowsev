package main

import (
	"fmt"
	"github.com/morten-krogh/gowsev/gowsev"
)

type Handler struct{
	id uint64
}

func (handler *Handler) ConnAccepted(context *gowsev.EvContext, id uint64) {
	handler.id = id
	fmt.Printf("Connection accepted %d\n", id)
}

func (handler *Handler) ConnClosed(context *gowsev.EvContext, id uint64) {
	fmt.Printf("Connection closed %d\n", id)
}

func (handler *Handler) EventLoopTimeout(context *gowsev.EvContext) {
	fmt.Printf("Timeout\n")
}

func (handler *Handler) MessageReceived(context *gowsev.EvContext, id uint64, message []byte) {
	fmt.Printf("Connection %d sent message %s\n", id, string(message))

	context.Write(handler.id, []byte("You sent me: " + string(message)))
}

func main() {

	var handler gowsev.Handler

	handler = &Handler{0}

	context := gowsev.MakeContext(&handler)

	context.ListenAndServe("9000")
	context.EventLoop()

}
