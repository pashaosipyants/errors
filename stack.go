package errors

import (
	"fmt"
	"io"
	"path"
	"runtime"
	"strings"
)

// annotatedStack represents a stack of program counters and messages bound to functions
type annotatedStack struct {
	stack []uintptr
	messages map[string]string
}

func (s annotatedStack) addMsg(pc uintptr, msg string) {
	name := runtime.FuncForPC(pc).Name()
	s.messages[name] = s.messages[name] + "\n" + msg
}

func (s annotatedStack) Format(st fmt.State, verb rune) {
	var messagesCopy map[string]string
	copy_map_msgs(s.messages, &messagesCopy)

	switch verb {
	case 'v':
		for _, pc := range s.stack {
			f := frame(pc)
			fmt.Fprintf(st, "\n%+v", f)
			name := runtime.FuncForPC(pc).Name()
			if msgs, ok := messagesCopy[name]; ok {
				fmt.Fprintf(st, "\nMESSAGES:")
				fmt.Fprintf(st, msgs)
				delete(messagesCopy, name)
			}
		}
		if len(messagesCopy) != 0 {
			fmt.Fprintf(st, "\n\nELSE MESSAGES:")
		}
		for funcname, msgs := range messagesCopy {
			fmt.Fprintf(st, "\n%s:", funcname)
			fmt.Fprintf(st, msgs)
		}
	}
}

func copy_map_msgs(src map[string]string, dst *map[string]string) {
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
type frame uintptr

// pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f frame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this frame's pc.
func (f frame) file() string {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc())
	return file
}

// line returns the line number of source code of the
// function for this frame's pc.
func (f frame) line() int {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.pc())
	return line
}

// Format formats the frame according to the fmt.Formatter interface.
//
//    %s    source file
//    %d    source line
//    %n    function name
//    %v    equivalent to %s:%d
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+s   function name and path of source file relative to the compile time
//          GOPATH separated by \n\t (<funcname>\n\t<path>)
//    %+v   equivalent to %+s:%d
func (f frame) Format(s fmt.State, verb rune) {
	switch verb {
	case 's':
		switch {
		case s.Flag('+'):
			pc := f.pc()
			fn := runtime.FuncForPC(pc)
			if fn == nil {
				io.WriteString(s, "unknown")
			} else {
				file, _ := fn.FileLine(pc)
				fmt.Fprintf(s, "%s\n\t%s", fn.Name(), file)
			}
		default:
			io.WriteString(s, path.Base(f.file()))
		}
	case 'd':
		fmt.Fprintf(s, "%d", f.line())
	case 'n':
		name := runtime.FuncForPC(f.pc()).Name()
		io.WriteString(s, funcname(name))
	case 'v':
		f.Format(s, 's')
		io.WriteString(s, ":")
		f.Format(s, 'd')
	}
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}
