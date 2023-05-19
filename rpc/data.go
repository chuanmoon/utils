package rpc

type SubscribeData struct {
	Table    string
	Id       int64
	ErrorMsg string
	ExtInfo  map[string]string
}
