package rpc

import (
	"errors"
	"net/rpc"
	"time"

	"github.com/chuanmoon/utils/mjson"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

const defaultTimeoutSeconds = 20

// vars
var (
	ErrShutdown = rpc.ErrShutdown
)

// Client nat Client
type NatClient struct {
	url    string
	conn   *nats.Conn
	logger *zap.Logger
}

// Dial Dial
func Dial(url string, _logger *zap.Logger) (Client, error) {
	conn, err := nats.Connect(url,
		nats.MaxReconnects(3),
		nats.ReconnectWait(time.Second))
	if err != nil {
		return nil, errors.New("Failed to connect to MQServer: " + err.Error())
	}
	return NewClientWithConn(conn, url, _logger), nil
}

// NewClientWithConn NewClientWithConn
func NewClientWithConn(conn *nats.Conn, url string, _logger *zap.Logger) Client {
	if _logger == nil {
		panic("logger is nil")
	}
	return &NatClient{
		url:    url,
		conn:   conn,
		logger: _logger,
	}
}

// JSONCall Call
func (client *NatClient) JSONCall(queue, serviceMethod string, args *[]byte, reply *[]byte) error {
	return client.JSONCallWithTimeout(queue, serviceMethod, defaultTimeoutSeconds, args, reply)
}
func (client *NatClient) JSONCallWithTimeout(queue, serviceMethod string, timeoutSeconds int64, args *[]byte, reply *[]byte) error {
	msg := &nats.Msg{
		Subject: queue,
		Data:    *args,
		Header: nats.Header{
			"ServiceMethod": []string{serviceMethod},
		},
		Sub: &nats.Subscription{
			Queue: queue,
		},
	}
	returnMsg, err := client.conn.RequestMsg(msg, time.Duration(timeoutSeconds)*time.Second)
	if err != nil {
		return err
	}
	remoteErr := returnMsg.Header.Get("Error")
	if remoteErr != "" {
		return errors.New(remoteErr)
	}
	*reply = returnMsg.Data
	return nil
}

// Call Call
func (client *NatClient) Call(subject, method string, args, receiver interface{}) error {
	return client.CallWithTimeout(subject, method, defaultTimeoutSeconds, args, receiver)
}
func (client *NatClient) CallWithTimeout(subject, method string, timeoutSeconds int64, args, receiver interface{}) error {
	startTime := time.Now()
	var argsbytes = []byte{}
	var replyBytes = []byte{}
	var err error
	defer func() {
		replyStr := string(replyBytes)
		if len(replyStr) > 1024 {
			replyStr = replyStr[:1024]
		}
		client.logger.Info("nats client call:",
			zap.String("used", time.Since(startTime).String()),
			zap.String("queue", subject),
			zap.String("method", method),
			zap.ByteString("args", argsbytes),
			zap.String("reply", replyStr),
		)
	}()
	if client == nil {
		return errors.New("[NATS Call] client is nil")
	}
	argsbytes, err = mjson.Marshal(args)
	if err != nil {
		return errors.New("[NATS Call] args Marshal error: " + err.Error())
	}

	err = client.JSONCallWithTimeout(subject, method, timeoutSeconds, &argsbytes, &replyBytes)
	if err != nil {
		return errors.New("[NATS Call] remote error: " + err.Error())
	}

	err = mjson.Unmarshal(replyBytes, receiver)
	if err != nil {
		return errors.New("[NATS Call] reply Unmarshal error: " + err.Error())
	}
	return nil
}
