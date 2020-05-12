package main

type Game struct {
	GameName string            // e.g. Holdem
	GameId   string            // e.g. g-01-001
	Version  string            // e.g. 1.0.1
	ThemeMap map[string]*Theme // Something like sub server
}

/** 建立 Game
 *
 * @param gameName
 * @param gameId
 * @param version
 * @return *Game 物件實例 */
func NewGame(gameName, gameId, version string) *Game {
	g := &Game{
		GameName: gameName,
		GameId:   gameId,
		Version:  version,
		ThemeMap: map[string]*Theme{},
	}

	return g
}

/** 初始化
 * 建立兩個 theme, 以及各自的 rooms */
func (g *Game) Init() {
	theme1Id := "1"
	theme1 := NewTheme(theme1Id, NewRoom1)
	g.ThemeMap[theme1Id] = theme1
	g.ThemeMap[theme1Id].RoomMap["1"] = g.ThemeMap[theme1Id].NewRoom("1")
	g.ThemeMap[theme1Id].RoomMap["2"] = g.ThemeMap[theme1Id].NewRoom("2")

	theme2Id := "2"
	theme2 := NewTheme(theme2Id, NewRoom2)
	g.ThemeMap[theme2Id] = theme2
	g.ThemeMap[theme2Id].RoomMap["1"] = g.ThemeMap[theme2Id].NewRoom("1")
	g.ThemeMap[theme2Id].RoomMap["2"] = g.ThemeMap[theme2Id].NewRoom("2")
}
