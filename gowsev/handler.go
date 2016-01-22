package gowsev

type Handler interface {
	ConnAccepted(context *EvContext, id uint64)
	ConnClosed(context *EvContext, id uint64)
	EventLoopTimeout(context *EvContext)
	MessageReceived(context *EvContext, id uint64, message []byte)
}
