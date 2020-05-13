package room

import (
	"log"
	"time"
)

type Room1 struct {
	IRoom
	fsm              *FSM
	enabled          bool
	id               string // uuid
	currentFrame     int64
	voidPeriodFrame  int64
	idlePeriodFrame  int64
	endPeriodFrame   int64
	periodTimer      map[string]int64
	void2IdleHandler FSMHandler
	idle2EndHandler  FSMHandler
}

/** 建立 Room1
 *
 * @param id uuid
 * @return IRoom 共用的抽象界面 */
func NewRoom1(id string) IRoom {
	r := &Room1{
		fsm:             NewFSM(VoidPeriod),
		enabled:         false,
		id:              id,
		currentFrame:    0,
		voidPeriodFrame: 0,
		idlePeriodFrame: 0,
		endPeriodFrame:  0,
		periodTimer: map[string]int64{
			"VoidPeriod": 10,
			"IdlePeriod": 20,
			"EndPeriod":  20},
	}
	r.InitFSM()

	return r
}

func (r *Room1) InitFSM() {
	r.void2IdleHandler = FSMHandler(func() FSMState {
		r.idlePeriodFrame = r.currentFrame
		r.Event1()

		log.Printf("[Room1-%s] 已經從 VoidPeriod 進入 IdlePeriod", r.id)
		return IdlePeriod
	})

	r.idle2EndHandler = FSMHandler(func() FSMState {
		r.endPeriodFrame = r.currentFrame
		r.Event2()

		log.Printf("[Room1-%s] 已經從 IdlePeriod 進入 EndPeriod", r.id)
		return EndPeriod
	})

	r.fsm.AddHandler(VoidPeriod, IdleEvent, r.void2IdleHandler)
	r.fsm.AddHandler(IdlePeriod, EndEvent, r.idle2EndHandler)
}

/** Room1 sample function */
func (r *Room1) Start() {
	log.Printf("[Room1-%s] start ...", r.id)

	go func() {
		tick := time.NewTicker(time.Duration(100) * time.Millisecond)
		r.enabled = true

		defer func() {
			log.Printf("[Room1-%s] end.", r.id)
			tick.Stop()
		}()

		for r.enabled {
			select {
			case <-tick.C:
				r.Update()
			default:
			}
		}
	}()
}

func (r *Room1) Update() {
	r.currentFrame++
	r.StateSwitch()
}

func (r *Room1) StateSwitch() {
	switch r.fsm.GetState() {
	case VoidPeriod:
		if r.currentFrame-r.voidPeriodFrame == r.periodTimer["VoidPeriod"] {
			r.fsm.Call(IdleEvent)
		}

	case IdlePeriod:
		if r.currentFrame-r.idlePeriodFrame == r.periodTimer["IdlePeriod"] {
			r.fsm.Call(EndEvent)
		}

	case EndPeriod:
		if r.currentFrame-r.endPeriodFrame == r.periodTimer["EndPeriod"] {
			r.enabled = false
		}
	}
}

func (r *Room1) Event1() {
	log.Printf("[Room1-%s] Execute Event1", r.id)
}

func (r *Room1) Event2() {
	log.Printf("[Room1-%s] Execute Event2", r.id)
}
