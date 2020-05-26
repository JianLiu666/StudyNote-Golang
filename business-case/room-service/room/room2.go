package room

import (
	"log"
	"time"
)

/** 建立 Room2
 *
 * @param id uuid
 * @return IRoom 共用的抽象界面 */
func NewRoom2(id string) IRoom {
	r := &Room2{
		fsm:             NewFSM(VoidPeriod),
		enabled:         false,
		id:              id,
		currentFrame:    0,
		voidPeriodFrame: 0,
		endPeriodFrame:  0,
		periodTimer: map[string]int64{
			"VoidPeriod": 10,
			"EndPeriod":  20},
	}
	r.InitFSM()

	return r
}

type Room2 struct {
	IRoom
	fsm             *FSM
	enabled         bool
	id              string // uuid
	currentFrame    int64
	voidPeriodFrame int64
	endPeriodFrame  int64
	periodTimer     map[string]int64
	void2EndHandler FSMHandler
}

func (r *Room2) InitFSM() {
	r.void2EndHandler = FSMHandler(func() FSMState {
		r.endPeriodFrame = r.currentFrame
		r.Event1()
		r.Event2()

		log.Printf("[Room2-%s] 已經從 VoidPeriod 進入 EndPeriod", r.id)
		return EndPeriod
	})

	r.fsm.AddHandler(VoidPeriod, EndEvent, r.void2EndHandler)
}

/** Room2 sample function */
func (r *Room2) Start() {
	log.Printf("[Room2-%s] start ...", r.id)

	go func() {
		tick := time.NewTicker(time.Duration(100) * time.Millisecond)
		r.enabled = true

		defer func() {
			log.Printf("[Room2-%s] end.", r.id)
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

func (r *Room2) Update() {
	r.currentFrame++
	r.StateSwitch()
}

func (r *Room2) StateSwitch() {
	switch r.fsm.GetState() {
	case VoidPeriod:
		if r.currentFrame-r.voidPeriodFrame == r.periodTimer["VoidPeriod"] {
			r.fsm.Call(EndEvent)
		}

	case EndPeriod:
		if r.currentFrame-r.endPeriodFrame == r.periodTimer["EndPeriod"] {
			r.enabled = false
		}
	}
}

func (r *Room2) Event1() {
	log.Printf("[Room2-%s] Execute Event1", r.id)
}

func (r *Room2) Event2() {
	log.Printf("[Room2-%s] Execute Event2", r.id)
}
