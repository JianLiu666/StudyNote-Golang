package sdk

import (
	"fmt"
	"messagebroker/pkg/broker"
	"messagebroker/pkg/packet"
	"reflect"
)

type IClient interface {
	// Subscribe 提供 User 訂閱指定主題, 並自動開始接收該主題內的所有訊息
	// @param topic 指定主題
	//
	// @return int status code
	Subscribe(topic string) (status int)

	// Publish 提供 User 推送訊息到指定主題的功能
	// @param topic 指定主題
	// @param msg 訊息內容
	//
	// @return int status code
	Publish(topic string, msg string) (status int)
}

// NewClient 對指定的 Message Broker 建立一個新的 Client 實例
// @param 提供 client 後續指定的 message broker server
//
// @return IClient instance
func NewClient(broker broker.IMessageBroker) IClient {
	client := &Client{
		broker: broker,
		topics: make([]reflect.SelectCase, 0, 4),
	}
	go client.listen()

	return client
}

type Client struct {
	broker broker.IMessageBroker
	topics []reflect.SelectCase
}

func (c *Client) Subscribe(topic string) (status int) {
	subscriber := c.broker.Subscribe(topic)
	c.topics = append(c.topics, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(subscriber),
	})

	return broker.StatusOK
}

func (c *Client) Publish(topic string, msg string) (status int) {
	return c.broker.Publish(topic, &packet.Payload{
		Message: msg,
	})
}

func (c *Client) listen() {
	for {
		_, recv, _ := reflect.Select(c.topics)
		fmt.Println(recv.Interface().(*packet.Payload).Message)
	}
}
