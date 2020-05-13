package main

import (
	"log"
	"time"
)

func main() {
	app := NewBaseModule()
	app.Init()

	log.Print(app.InstanceId)
	for _, gameMap := range app.GameMap {
		for _, theme := range gameMap.ThemeMap {
			for _, room := range theme.RoomMap {
				room.Start()
			}
		}
	}

	time.Sleep(6 * time.Second)
}
