package cyrpc

import (
	"context"
	"sync"
	"time"

	"github.com/chuanmoon/utils/cyconfig"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Client interface {
	CallJson(subject, method string, args, receiver interface{}) error
	CallJsonWithTimeout(subject, method string, timeoutSeconds int64, args, receiver interface{}) error
	CallBytes(subject, method string, args, receiver *[]byte, encodeType EncodeType) error
	CallBytesWithTimeout(subject, method string, timeoutSeconds int64, args, receiver *[]byte, encodeType EncodeType) error
	CallMsgpack(subject, method string, args, receiver interface{}) error
	CallMsgpackWithTimeout(subject, method string, timeoutSeconds int64, args, receiver interface{}) error
	Subscribe(subj string, callback func(*SubscribeData)) error
	Publish(subj string, data *SubscribeData) error
}

type Tool struct {
	natsConn *nats.Conn
	natsUrl  string
	appName  string
	logger   *zap.Logger
	client   Client
	mutex    sync.Mutex
}

func NewTool(appName string, logger *zap.Logger) *Tool {
	var err error
	if logger == nil {
		logger, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	}

	t := &Tool{
		appName: appName,
		logger:  logger,
	}

	t.natsUrl = cyconfig.String("nats_url", "0.0.0.0:4222")
	t.natsConn, err = nats.Connect(t.natsUrl,
		nats.Name(appName),
		nats.DontRandomize(),
		nats.MaxReconnects(3),
		nats.ReconnectWait(2*time.Second),
		nats.ClosedHandler(func(conn *nats.Conn) {
			conn.Close()
		}),
	)
	if err != nil {
		panic(err)
	}
	return t
}

func (t *Tool) GetClient() Client {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.client != nil {
		return t.client
	}

	t.client = NewClientWithConn(t.natsConn, t.natsUrl, t.logger)
	return t.client
}

func (t *Tool) StartServer(objs ...interface{}) {
	server := NewServer(t.logger)
	for _, obj := range objs {
		server.Register(obj)
	}

	t.logger.Info(t.appName + " started...")
	server.ServeConn(context.Background(), t.natsConn, t.appName)
}
