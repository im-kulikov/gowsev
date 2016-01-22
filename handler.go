package gowsev

type Handler interface {
	connAccepted(context *EvContext, id uint64)
	connClosed(context *EvContext, id uint64)
	timerFired(context *EvContext)
	messageReceived(context *EvContext, id uint64, message []byte) 
}
