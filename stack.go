package errors

import (
	"fmt"
	"io"
	"runtime"
)

// annotatedStack represents a stack of program counters and annotations bounded to functions
type annotatedStack struct {
	stack       []uintptr
	annotations map[string]string
}

func (s annotatedStack) addAnnotation(skip int, annot string) {
	pc, _, _, _ := runtime.Caller(skip + 1)
	name := runtime.FuncForPC(pc).Name()
	s.annotations[name] = s.annotations[name] + "\n" + annot
}

func (s annotatedStack) Format(st fmt.State, verb rune) {
	var annotationsCopy map[string]string
	copy_map_annots(s.annotations, &annotationsCopy)

	for _, pc := range s.stack {
		f := frame(pc)
		fmt.Fprint(st, f, "\n")
		name := runtime.FuncForPC(f.pc()).Name()
		if msgs, ok := annotationsCopy[name]; ok {
			fmt.Fprint(st, "ANNOTATIONS:")
			fmt.Fprint(st, msgs, "\n")
			delete(annotationsCopy, name)
		}
	}
	if len(annotationsCopy) != 0 {
		fmt.Fprint(st, "\nELSE ANNOTATIONS:")
	}
	for funcname, msgs := range annotationsCopy {
		fmt.Fprintf(st, "\n%s:", funcname)
		fmt.Fprint(st, msgs)
	}
}

func copy_map_annots(src map[string]string, dst *map[string]string) {
	if src != nil {
		*dst = make(map[string]string, len(src))
		for k, v := range src {
			(*dst)[k] = v
		}
	}
}

// callers with skip==0 returns stack of program counters starting from caller of callers
// with greater skip it skips more stack frames starting from upper level of invocations
func callers(skip int) []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(2+skip, pcs[:])
	return pcs[0:n]
}

// frame represents a program counter inside a stack frame.
// For historical reasons if Frame is interpreted as a uintptr
// its value represents the program counter + 1.
type frame uintptr

func (f frame) pc() uintptr { return uintptr(f) - 1 }

func (f frame) Format(s fmt.State, verb rune) {
	pc := f.pc()
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		io.WriteString(s, "unknown")
	} else {
		file, line := fn.FileLine(pc)
		fmt.Fprintf(s, "%s\n\t%s:%d", fn.Name(), file, line)
	}
}
