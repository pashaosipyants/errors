package errors

/*
	sprint an error with all the details
*/

import (
	"fmt"
	"io"
	"runtime"
	"strings"
)

// Returns detailed error description.
// With error msg, stacktrace, annotations, and suppressed errors.
func SprintE(err error) string {
	var b strings.Builder

	_, _ = fmt.Fprintln(&b, "ERROR:")
	_, _ = fmt.Fprintln(&b, err)

	var (
		errAnno     *_errorAnnotation
		annotations = make(map[string]string)
		as          annotatedStack
	)
	for errIteration := err; AsE(errIteration, &errAnno); errIteration = errAnno.error {
		annotations[errAnno.where] = errAnno.annotation
	}
	as.annotations = annotations

	var errStack *_errorStack
	if AsE(err, &errStack) {
		_, _ = fmt.Fprint(&b, "\n")
		_, _ = fmt.Fprintln(&b, "STACK:")
		as.stack = errStack.stack
	}
	_, _ = fmt.Fprint(&b, as)

	supps := SuppressedE(err)
	for _, s := range supps {
		_, _ = fmt.Fprint(&b, "\n")
		_, _ = fmt.Fprintln(&b, "SUPPRESSED:")
		_, _ = fmt.Fprintln(&b, SprintE(s))
	}

	return b.String()
}

type annotatedStack struct {
	stack       []uintptr
	annotations map[string]string
}

func (s annotatedStack) Format(st fmt.State, verb rune) {
	for _, pc := range s.stack {
		f := frame(pc)
		_, _ = fmt.Fprint(st, f, "\n")
		name := runtime.FuncForPC(f.pc()).Name()
		if msgs, ok := s.annotations[name]; ok {
			_, _ = fmt.Fprint(st, "\tANNOTATION: ")
			_, _ = fmt.Fprint(st, msgs, "\n")
			delete(s.annotations, name)
		}
	}
	if len(s.annotations) != 0 {
		_, _ = fmt.Fprint(st, "\nELSE ANNOTATIONS:")
	}
	for funcname, msgs := range s.annotations {
		_, _ = fmt.Fprintf(st, "\n%s:", funcname)
		_, _ = fmt.Fprint(st, msgs)
	}
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
