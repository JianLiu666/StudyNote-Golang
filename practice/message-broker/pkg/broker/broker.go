package broker

import (
	"fmt"
	"messagebroker/pkg/packet"
	"sync"
	"time"
)

const (
	StatusOK = iota
	StatusTopicNotFound
)

type IMessageBroker interface {
	// Subscribe 提供 client 訂閱指定主題
	// @param topic 指定主題
	//
	// @return chan *Flow 指定主題的消息隊列
	Subscribe(topic string) (subscriber chan *packet.Payload)

	// Publish 提供 Client 推送訊息到指定主題
	// @param topic 指定主題
	// @param data 訊息封包
	//
	// @return int status code
	Publish(topic string, data *packet.Payload) (status int)
}

// NewMsgBroker 建立 Message Broker 實例
//
// @return IMessageBroker instance
func NewMsgBroker() IMessageBroker {
	return &Broker{}
}

type Broker struct {
	topics sync.Map
}

func (b *Broker) Subscribe(topic string) (subscriber chan *packet.Payload) {
	// TODO: magic number
	subscriber = make(chan *packet.Payload, 100)

	if val, ok := b.topics.Load(topic); ok {
		b.topics.Store(topic, append(val.([]chan *packet.Payload), subscriber))
	} else {
		b.topics.Store(topic, []chan *packet.Payload{subscriber})
	}

	return subscriber
}

func (b *Broker) Publish(topic string, data *packet.Payload) (status int) {
	val, ok := b.topics.Load(topic)

	// topic not found
	if !ok {
		return StatusTopicNotFound
	}

	// broadcast to all subscribers
	for _, subscriber := range val.([]chan *packet.Payload) {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("timeout")
		case subscriber <- data:
		}
	}

	return StatusOK
}
