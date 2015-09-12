package unsafer

import (
	"runtime"
)

var stopTheWorld func(string)
var startTheWorld func()

func StopTheWorld(reason string) {
	stopTheWorld(reason)
}

func StartTheWorld() {
	startTheWorld()
}

func init() {
	runtime.Stack(nil, false)
	InsertFunction("runtime.stopTheWorld", &stopTheWorld)
	InsertFunction("runtime.startTheWorld", &startTheWorld)
}
