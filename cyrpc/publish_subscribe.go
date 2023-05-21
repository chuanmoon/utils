package cyrpc

import (
	"errors"

	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
)

// Subscribe 订阅
func (client *NatClient) Subscribe(subj string, callback func(*SubscribeData)) error {
	_, err := client.conn.Subscribe(subj, func(m *nats.Msg) {
		var data SubscribeData
		err := msgpack.Unmarshal(m.Data, &data)
		if err != nil {
			data.ErrorMsg = "[NATS Subscribe] callback: " + err.Error()
		}
		callback(&data)
	})
	if err != nil {
		return errors.New("[NATS Subscribe] error: " + err.Error())
	}
	return nil
}

// Publish 发布
func (client *NatClient) Publish(subj string, data *SubscribeData) error {
	b, err := msgpack.Marshal(data)
	if err != nil {
		return errors.New("[NATS Publish] remote error: " + err.Error())
	}

	err = client.conn.Publish(subj, b)
	if err != nil {
		return errors.New("[NATS Publish] remote error: " + err.Error())
	}
	return nil
}
