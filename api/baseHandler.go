package api

type BaseHandler struct {
	addresses [][]string // Arrays are already pointers. We don't need to pass in pointer here
}

func NewBaseHandler(addresses [][]string) *BaseHandler {
	return &BaseHandler{
		addresses: addresses,
	}
}
