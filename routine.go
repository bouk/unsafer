package unsafer

import (
	"unsafe"
)

type GoBuf struct {
	// The offsets of sp, pc, and g are known to (hard-coded in) libmach.
	Sp   uintptr
	Pc   uintptr
	G    *Goroutine
	Ctxt unsafe.Pointer // this has to be a pointer so that gc scans it
	Ret  uintptr
	Lr   uintptr
	Bp   uintptr // for GOEXPERIMENT=framepointer
}

type Gcstats struct {
	Nhandoff    uint64
	Nhandoffcnt uint64
	Nprocyield  uint64
	Nosyield    uint64
	Nsleep      uint64
}

type M struct {
	G0      *Goroutine // goroutine with scheduling stack
	Morebuf GoBuf      // gobuf arg to morestack
	Divmod  uint32     // div/mod denominator for arm - known to liblink

	// Fields not known to debuggers.
	Procid        uint64     // for debuggers, but offset not hard-coded
	Gsignal       *Goroutine // signal-handling g
	Sigmask       [4]uintptr // storage for saved signal mask
	Tls           [4]uintptr // thread-local storage (for x86 extern register)
	Mstartfn      func()
	Curg          *Goroutine // current running goroutine
	Caughtsig     uintptr    // goroutine running during fatal signal
	P             uintptr    // attached p for executing go code (nil if not executing go code)
	Nextp         uintptr
	ID            int32
	Mallocing     int32
	Throwing      int32
	Preemptoff    string // if != "", keep curg running on this m
	Locks         int32
	Softfloat     int32
	Dying         int32
	Profilehz     int32
	Helpgc        int32
	Spinning      bool // m is out of work and is actively looking for work
	Blocked       bool // m is blocked on a note
	Inwb          bool // m is executing a write barrier
	Printlock     int8
	Fastrand      uint32
	Ncgocall      uint64 // number of cgo calls in total
	Ncgo          int32  // number of cgo calls currently in progress
	Park          uintptr
	Alllink       *M // on allm
	Schedlink     uintptr
	Machport      uint32 // return address for mach ipc (os x)
	Mcache        uintptr
	Lockedg       *Goroutine
	Createstack   [32]uintptr // stack that created this thread.
	Freglo        [16]uint32  // d[i] lsb and f[i]
	Freghi        [16]uint32  // d[i] msb and f[i+16]
	Fflag         uint32      // floating point compare flags
	Locked        uint32      // tracking for lockosthread
	Nextwaitm     uintptr     // next m waiting for lock
	Waitsema      uintptr     // semaphore for parking on locks
	Waitsemacount uint32
	Waitsemalock  uint32
	Gcstats       Gcstats
	Needextram    bool
	Traceback     uint8
	Waitunlockf   unsafe.Pointer // todo go func(*g, unsafe.pointer) bool
	Waitlock      unsafe.Pointer
	Waittraceev   byte
	Waittraceskip int
	Startingtrace bool
	Syscalltick   uint32
}

type Goroutine struct {
	// Stack parameters.
	// stack describes the actual stack memory: [stack.lo, stack.hi).
	// stackguard0 is the stack pointer compared in the Go stack growth prologue.
	// It is stack.lo+StackGuard normally, but can be StackPreempt to trigger a preemption.
	// stackguard1 is the stack pointer compared in the C stack growth prologue.
	// It is stack.lo+StackGuard on g0 and gsignal stacks.
	// It is ~0 on other goroutine stacks, to trigger a call to morestackc (and crash).
	Stack struct {
		Lo, Hi uintptr
	} // offset known to runtime/cgo
	Stackguard0 uintptr // offset known to liblink
	Stackguard1 uintptr // offset known to liblink

	Panic      uintptr // innermost panic - offset known to liblink
	Defer      *Defer  // innermost defer
	M          *M      // current m; offset known to arm liblink
	StackAlloc uintptr // stack allocation is [stack.lo,stack.lo+stackAlloc)
	Sched      GoBuf
	Syscallsp  uintptr // if status==Gsyscall, syscallsp = sched.sp to use during gc
	Syscallpc  uintptr // if status==Gsyscall, syscallpc = sched.pc to use during gc
	Stkbar     []struct {
		SavedLRPtr, SavedLRVal uintptr // stack barriers, from low to high
	}
	StkbarPos      uintptr        // index of lowest stack barrier not hit
	Param          unsafe.Pointer // passed parameter on wakeup
	Atomicstatus   uint32
	StackLock      uint32 // sigprof/scang lock; TODO: fold in to atomicstatus
	Goid           int64
	Waitsince      int64  // approx time when the g become blocked
	Waitreason     string // if status==Gwaiting
	Schedlink      *Goroutine
	Preempt        bool   // preemption signal, duplicates stackguard0 = stackpreempt
	Paniconfault   bool   // panic (instead of crash) on unexpected fault address
	Preemptscan    bool   // preempted g does scan for gc
	Gcscandone     bool   // g has scanned stack; protected by _Gscan bit in status
	Gcscanvalid    bool   // false at start of gc cycle, true if G has not run since last scan
	Throwsplit     bool   // must not split stack
	Raceignore     int8   // ignore race detection events
	Sysblocktraced bool   // StartTrace has emitted EvGoInSyscall about this goroutine
	Sysexitticks   int64  // cputicks when syscall has returned (for tracing)
	Sysexitseq     uint64 // trace seq when syscall has returned (for tracing)
	Lockedm        uintptr
	Sig            uint32
	Writebuf       []byte
	Sigcode0       uintptr
	Sigcode1       uintptr
	Sigpc          uintptr
	Gopc           uintptr // pc of go statement that created this goroutine
	Startpc        uintptr // pc of goroutine function
	Racectx        uintptr
	Waiting        uintptr // sudog structures this g is waiting on (that have a valid elem ptr)
	Readyg         uintptr // scratch for readyExecute

	// Per-G gcController state
	Gcalloc    uintptr // bytes allocated during this GC cycle
	Gcscanwork int64   // scan work done (or stolen) this GC cycle
}

var tracebackothers func(*Goroutine)
var systemstack func(func())

func TracebackOthers() {
	systemstack(func() {
		tracebackothers(nil)
	})
}

func GetG() *Goroutine

func init() {
	InsertFunction("runtime.tracebackothers", &tracebackothers)
	InsertFunction("runtime.systemstack", &systemstack)
}
