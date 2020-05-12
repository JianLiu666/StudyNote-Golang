package main

type Theme struct {
	Pass    string                // yeah you already know
	RoomMap map[string]IRoom      // yeah you already know
	NewRoom func(id string) IRoom // callback function
}

/** 建立 Theme
 *
 * @param pass guess what it is
 * @param newRoomCallback 建立房間的方法
 * @return *Theme 物件實例 */
func NewTheme(pass string, newRoomCallback func(id string) IRoom) *Theme {
	t := &Theme{
		Pass:    pass,
		RoomMap: map[string]IRoom{},
		NewRoom: newRoomCallback,
	}

	return t
}
