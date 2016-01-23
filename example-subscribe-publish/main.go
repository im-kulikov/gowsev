package main

import (
	"fmt"
	"github.com/morten-krogh/gowsev/gowsev"
)

/* This program is an example program using the gowsev libary.
 * A websocket connection can send four types of messages.
 *
 * 1. "subscribe:(channel)" where (channel) is any string representing a channel, e.g. "subscribe:123"
 * 2. "unsubscribe:(channel)"
 * 3. "unsusbscribe-all"
 * 4. "publish:(channel):(body)"
 *
 * (channel) and (message) are placeholders for any string. (channel) can not contain colon ":".
 *
 * The operation is straight forward: A publish message leads to body being sent to all web sockets subscribed to that channel.
 */

type Service struct {
	idChMap map[uint64](map[string]struct{}) // wsChMap[id][ch] exists if the websocket id is subscribed to the channel ch.
	chIdMap map[string](map[uint64]struct{}) // chWsMap[ch][id] exists if the websocket id is subscribed to the channel ch.
}

func (service *Service) ConnAccepted(context *gowsev.EvContext, id uint64) {
	fmt.Printf("Connection accepted %d\n", id)
}

func (service *Service) ConnClosed(context *gowsev.EvContext, id uint64) {
	fmt.Printf("Connection closed %d\n", id)
	service.unsubscribe_all(id)
}

func (service *Service) EventLoopTimeout(context *gowsev.EvContext) {
	fmt.Printf("Timeout\n")
}

func (service *Service) MessageReceived(context *gowsev.EvContext, id uint64, message []byte) {
	messageStr := string(message)
	service.handle_message(context, id, messageStr)
	fmt.Printf("Connection %d sent message %s\n", id, messageStr)
}

func (service *Service) handle_message(context *gowsev.EvContext, id uint64, message string) {
	var channel, body string

	_, err := fmt.Sscanf(message, "subscribe:%s", &channel)
	if err == nil {
		service.subscribe(id, channel)
		return
	}

	_, err = fmt.Sscanf(message, "unsubscribe:%s", &channel)
	if err == nil {
		service.unsubscribe(id, channel)
		return
	}

	if message == "unsubsribe-all" {
		service.unsubscribe_all(id)
		return
	}

	_, err = fmt.Sscanf(message, "publish:%s:%s", &channel, &body)
	if err == nil {
		service.publish(context, channel, body)
		return
	}

	context.Close(id)
}

func (service *Service) subscribe(id uint64, ch string) {
	if service.idChMap[id] == nil {
		service.idChMap[id] = make(map[string]struct{})
	}
	service.idChMap[id][ch] = struct{}{}

	if service.chIdMap[ch] == nil {
		service.chIdMap[ch] = make(map[uint64]struct{})
	}
	service.chIdMap[ch][id] = struct{}{}
}

func (service *Service) unsubscribe(id uint64, ch string) {
	channels, ok := service.idChMap[id]
	if ok {
		delete(channels, ch)
	}

	ids, ok := service.chIdMap[ch]
	if ok {
		delete(ids, id)
	}
}

func (service *Service) unsubscribe_all(id uint64) {
	channels, ok := service.idChMap[id]
	if ok {
		for ch, _ := range channels {
			delete(service.chIdMap[ch], id)
		}
		delete(service.idChMap, id)
	}
}

func (service *Service) publish(context *gowsev.EvContext, ch string, message string) {
	ids, ok := service.chIdMap[ch]
	if ok {
		for id, _ := range ids {
			context.Write(id, []byte(message))
		}
	}
}

func main() {

	var service gowsev.Handler

	service = &Service{nil, nil}

	context := gowsev.MakeContext(&service)

	context.ListenAndServe("9000")
	context.EventLoop()

}
