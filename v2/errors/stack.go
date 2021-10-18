package errors

import (
	"fmt"
	"runtime"
	"strings"
)

func newStackTrace() StackTrace {
	callers := callers(4)

	var frames []Frame
	for _, c := range callers {
		frames = append(frames, NewFrame(c))
	}

	return StackTrace{
		Frames: frames,
	}
}

type StackTrace struct {
	Frames []Frame
}

func (t StackTrace) StringStack() []string {
	var trace []string
	for _, f := range t.Frames {
		trace = append(trace, f.String())

	}

	return trace
}

func (t StackTrace) String() string {
	return strings.Join(t.StringStack(), "\n")
}

func NewFrame(f uintptr) Frame {
	return Frame{pc: f - 1}
}

type Frame struct {
	pc uintptr
}

func (f Frame) String() string {
	return f.FileName() + " " + f.FunctionName()
}

func (f Frame) Format(s fmt.State, verb rune) {
	//	filename := f.file()
	//	_, _ = io.WriteString(s, "f: ")
	//	_, _ = io.WriteString(s, filename)
	//	_, _ = io.WriteString(s, ": ")
	// _, _ = io.WriteString(s, funcname(f.name()))
	//	_, _ = io.WriteString(s, "\n")
}

//func (f Frame) pc() uintptr { return uintptr(f) - 1 }

// file returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) FileName() string {
	fn := runtime.FuncForPC(f.pc)
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.pc)
	return file
}

// name returns the name of this function, if known.
func (f Frame) FunctionName() string {
	fn := runtime.FuncForPC(f.pc)
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}

//func funcname(name string) string {
//	i := strings.LastIndex(name, "/")
//	name = name[i+1:]
//	i = strings.Index(name, ".")
//	return name[i+1:]
//}

func callers(skip int) []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip, pcs[:])
	var st []uintptr = pcs[0:n]
	return st
}
