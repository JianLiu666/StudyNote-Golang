package main

import "github.com/JianLiu666/LearningNote-Golang/businessCase/roomService/room"

type Theme struct {
	Settings map[string]interface{}     // Something like ante, or business logic ...
	RoomMap  map[string]room.IRoom      // Room instances
	NewRoom  func(id string) room.IRoom // Callback
}

/** 建立 Theme
 *
 * @param ante Business logic condition
 * @param newRoomCallback 建立房間的方法
 * @return *Theme 物件實例 */
func NewTheme(ante string, newRoomCallback func(id string) room.IRoom) *Theme {
	t := &Theme{
		RoomMap: map[string]room.IRoom{},
		NewRoom: newRoomCallback,
	}

	return t
}
