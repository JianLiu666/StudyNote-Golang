package main

// 1. 不要把 Context 放在 struct 內，要用參數的方始傳遞
// 2. 以參數傳遞 Context 時，應該要把 Context 參數放在第一位
// 3. 需要傳遞 Conext 時如果不知道要傳什麼，寧願傳 Conext.TODO 也不要 nil
// 4. Conext 在多執行序中是安全的，可以放心在多個 Goroutine 中被傳遞
func main() {
	multiRoutine()
	// withValue()
}
