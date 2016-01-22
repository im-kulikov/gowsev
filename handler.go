package gowsev

type Handler interface {
	connAccepted(context *EvContext, id uint64)
	connClosed(context *EvContext, id uint64)
	eventLoopTimeout(context *EvContext)
	messageReceived(context *EvContext, id uint64, message []byte) 
}
