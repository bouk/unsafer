package unsafer

import (
	"runtime"
)

var (
	StopTheWorld  func(string)
	StartTheWorld func()
)

func init() {
	runtime.Stack(nil, false)
	InsertFunction("runtime.stopTheWorld", &StopTheWorld)
	InsertFunction("runtime.startTheWorld", &StartTheWorld)
}
