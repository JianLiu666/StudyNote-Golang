package main

import (
	"time"

	"messagebroker/pkg/broker"
	"messagebroker/pkg/sdk"
)

const orderTopic = "orderCreated"

func main() {
	broker := broker.NewMsgBroker()

	client1 := sdk.NewClient(broker)
	client2 := sdk.NewClient(broker)
	client3 := sdk.NewClient(broker)

	// 客戶訂閱主題，收到 broker 訊息後送進通道
	client1.Subscribe(orderTopic)
	client2.Subscribe(orderTopic)

	// 客戶對指定主題送出訊息
	client3.Publish(orderTopic, "Hello")

	time.Sleep(5 * time.Second)
}
