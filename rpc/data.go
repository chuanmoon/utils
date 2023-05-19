package rpc

type EncodeType string

const (
	EncodeTypeJson    EncodeType = "1"
	EncodeTypeMsgpack EncodeType = "2"
)

type SubscribeData struct {
	Table    string
	Id       int64
	ErrorMsg string
	ExtInfo  map[string]string
}
