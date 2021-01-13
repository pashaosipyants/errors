package errors

import "runtime"

// Maximum depth of stack recorded to an error.
// One can change it, but be cautious.
// Set value only in initialization part of a program before any error is created
// not to cause concurrency problems.
const ErrStackMaxDepth = 32

type stack []uintptr

// callers with skip==0 returns stack of program counters starting from caller of callers
// with greater skip it skips more stack frames starting from upper level of invocations
func callers(skip int) stack {
	var pcs [ErrStackMaxDepth]uintptr
	n := runtime.Callers(2+skip, pcs[:])
	pcsReduced := make([]uintptr, n)
	copy(pcsReduced, pcs[:])
	return pcsReduced
}
