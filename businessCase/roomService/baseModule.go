package main

import (
	"fmt"
	"time"
)

type BaseModule struct {
	Service    interface{}      // Something like nats, redis client or grpc ...
	Settings   interface{}      // Something like heartbeat, connection timeout or etc. ...
	InstanceId string           // Service instance id (uuid)
	GameMap    map[string]*Game // Game list
}

/** 建立 BaseModule
 *
 * @return *BaseModule 物件實例 */
func NewBaseModule() *BaseModule {
	b := &BaseModule{
		InstanceId: fmt.Sprintf("%v", time.Now().UnixNano()),
		GameMap:    map[string]*Game{},
	}

	return b
}

/** 初始化 */
func (b *BaseModule) Init() {
	game1 := NewGame("g01", "g01", "1.0.1")
	game1.Init()

	b.GameMap[game1.GameId] = game1
}
