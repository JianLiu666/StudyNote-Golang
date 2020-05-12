package main

import "log"

type IRoom interface {
	Start()
}

type Room1 struct {
	IRoom
	Id string // uuid
}

/** 建立 Room1
 *
 * @param id uuid
 * @return IRoom 共用的抽象界面 */
func NewRoom1(id string) IRoom {
	r := &Room1{
		Id: id,
	}

	return r
}

/** Room1 sample function */
func (r Room1) Start() {
	log.Printf("[Room1-%s] start ...", r.Id)
}

type Room2 struct {
	IRoom
	Id string // uuid
}

/** 建立 Room2
 *
 * @param id uuid
 * @return IRoom 共用的抽象界面 */
func NewRoom2(id string) IRoom {
	r := &Room2{
		Id: id,
	}

	return r
}

/** Room2 sample function */
func (r Room2) Start() {
	log.Printf("[Room2-%s] start ...", r.Id)
}
