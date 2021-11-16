package errors

import (
	"fmt"
	"runtime"
	"strings"
)

func newStackTrace(skip int) StackTrace {
	callers := callers(skip)

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
	pc := f - 1

	return Frame{
		SystemFileName: NewFileNameFromPC(pc),
		FunctionName:   NewFunctionFromPC(pc),
	}
}

func NewFileNameFromPC(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	file, _ := fn.FileLine(pc)

	return file
}

type Frame struct {
	// not very useful, since it is system dependant
	SystemFileName string

	FunctionName Function
}

func (f Frame) String() string {
	return f.FunctionName.String()
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

func NewFunctionFromPC(pc uintptr) Function {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return Function{
			NameWithPath:    "unknown",
			NameWithPackage: "unknown",
			Name:            "unknown",
		}
	}

	withPath := fn.Name()
	i := strings.LastIndex(withPath, "/")
	withPackage := withPath[i+1:]
	j := strings.Index(withPackage, ".")
	name := withPackage[j+1:]

	_, lineNum := fn.FileLine(pc)

	return Function{
		NameWithPath:    withPath,
		NameWithPackage: withPackage,
		Name:            name,
		LineNumber:      lineNum,
	}
}

type Function struct {
	NameWithPath    string
	NameWithPackage string
	Name            string
	LineNumber      int
}

func (f Function) String() string {
	return fmt.Sprintf("%s:%d", f.NameWithPath, f.LineNumber)
}

func callers(skip int) []uintptr {
	const depth = 32

	var pcs [depth]uintptr

	n := runtime.Callers(skip, pcs[:])
	st := pcs[0:n]

	return st
}
