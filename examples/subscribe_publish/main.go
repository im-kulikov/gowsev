package main

import (
	"fmt"
	"github.com/morten-krogh/gowsev/gowsev"
	"strings"
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

/* implementation of the gowsev Handler interface */
func (service *Service) ConnAccepted(context *gowsev.Context, id uint64) {
	fmt.Printf("Connection accepted %d\n", id)
}

func (service *Service) ConnClosed(context *gowsev.Context, id uint64) {
	fmt.Printf("Connection closed %d\n", id)
	service.unsubscribe_all(id)
}

func (service *Service) EventLoopTimeout(context *gowsev.Context) {
	fmt.Printf("Timeout\n")
}

func (service *Service) MessageReceived(context *gowsev.Context, id uint64, message []byte) {
	messageStr := string(message)
	service.handle_message(context, id, messageStr)
	fmt.Printf("Connection %d sent message %s\n", id, messageStr)
}

func (service *Service) handle_message(context *gowsev.Context, id uint64, message string) {
	var channel, body string

	_, err := fmt.Sscanf(message, "subscribe:%s", &channel)
	if err == nil {
		context.Write(id, []byte("OK"))
		service.subscribe(id, channel)
		return
	}

	_, err = fmt.Sscanf(message, "unsubscribe:%s", &channel)
	if err == nil {
		context.Write(id, []byte("OK"))
		service.unsubscribe(id, channel)
		return
	}

	if message == "unsubscribe-all" {
		context.Write(id, []byte("OK"))
		service.unsubscribe_all(id)

		return
	}

	components := strings.SplitN(message, ":", 3)
	if len(components) >= 3 && components[0] == "publish" {
		context.Write(id, []byte("OK"))
		channel = components[1]
		body = components[2]
		service.publish(context, channel, body)
		return
	}
	
	context.Write(id, []byte("Invalid message"))
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

func (service *Service) publish(context *gowsev.Context, ch string, message string) {
	ids, ok := service.chIdMap[ch]
	if ok {
		for id, _ := range ids {
			context.Write(id, []byte(message))
		}
	}
}

func main() {

	var service gowsev.Handler

	idChMap := make(map[uint64](map[string]struct{}))
	chIdMap := make(map[string](map[uint64]struct{}))
	service = &Service{idChMap, chIdMap}

	context := gowsev.MakeContext(&service)

	context.ListenAndServe("9000")
	context.ListenAndServeTLS("9001", "localhost.crt", "localhost.key")
	fmt.Printf("The subscribe publish service is running on port 9000\n")
	context.EventLoop()
}
