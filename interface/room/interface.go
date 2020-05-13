package room

type IRoom interface {
	Start()
	Update()
}

var (
	VoidPeriod = FSMState("VoidPeriod")
	IdlePeriod = FSMState("IdlePeriod")
	EndPeriod  = FSMState("EndPeriod")

	VoidEvent = FSMEvent("進入 VoidPeriod")
	IdleEvent = FSMEvent("進入 IdlePeriod")
	EndEvent  = FSMEvent("進入 EndPeriod")
)
