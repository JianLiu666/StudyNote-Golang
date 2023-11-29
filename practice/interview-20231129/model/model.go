package model

type User struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	Gender   int    `json:"gender"`
	NumDates int    `json:"numDates"`
}
