package gowsev

type Handler interface {
	ConnAccepted(context *Context, id uint64)
	ConnClosed(context *Context, id uint64)
	EventLoopTimeout(context *Context)
	MessageReceived(context *Context, id uint64, message []byte)
}
