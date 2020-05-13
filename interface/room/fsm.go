package room

import (
	"log"
	"sync"
)

type FSMState string            // 定義狀態
type FSMEvent string            // 定義事件
type FSMHandler func() FSMState // 定義處理方法

// 有限狀態機
type FSM struct {
	mu       sync.Mutex                           // 排他鎖
	state    FSMState                             // 狀態機當前狀態
	handlers map[FSMState]map[FSMEvent]FSMHandler // 狀態機圖
}

/** 取得狀態機當前狀態
 *
 * @return FSMState 狀態機當前狀態 */
func (f *FSM) GetState() FSMState {
	return f.state
}

/** 變更狀態機為指定狀態
 *
 * @param newState 指定狀態 */
func (f *FSM) SetState(newState FSMState) {
	f.state = newState
}

/** 在目標狀態上加入新事件/事件處理
 *
 * @param state 目標狀態
 * @param event 欲加入的新事件
 * @param handler 欲加入的新事件處理
 * @return FSM 更新後的有限狀態機 */
func (f *FSM) AddHandler(state FSMState, event FSMEvent, handler FSMHandler) *FSM {
	if _, ok := f.handlers[state]; !ok {
		f.handlers[state] = make(map[FSMEvent]FSMHandler)
	}
	if _, ok := f.handlers[state][event]; ok {
		log.Printf("[警告] 狀態(%s)事件(%s)已定義過", state, event)
	}
	f.handlers[state][event] = handler
	return f
}

/** 向狀態機請求事件處理
 *
 * @param event 事件定義
 * @return FSMState 狀態機當前新狀態 */
func (f *FSM) Call(event FSMEvent) FSMState {
	f.mu.Lock()
	defer f.mu.Unlock()
	events := f.handlers[f.GetState()]
	if events == nil {
		return f.GetState()
	}
	if fn, ok := events[event]; ok {
		// oldState := f.GetState()
		f.SetState(fn())
		// newState := f.GetState()
		// log.Printf("狀態從 [%v] 變成 [%v]", string(oldState), string(newState))
	} else {
		log.Print("切換狀態失敗")
	}
	return f.GetState()
}

/** 實體化有限狀態機
 *
 * @param initState 設定狀態機初始狀態
 * @return FSM 有限狀態機 */
func NewFSM(initState FSMState) *FSM {
	return &FSM{
		state:    initState,
		handlers: make(map[FSMState]map[FSMEvent]FSMHandler),
	}
}

/** 確認當前狀態有無在指定狀態內
 *
 * @param  a 指定狀態列表
 * @return  是否有在狀態內 */
func (f *FSM) CheckState(a ...FSMState) bool {
	for _, fsm := range a {
		if fsm == f.GetState() {
			return true
		}
	}
	return false
}
