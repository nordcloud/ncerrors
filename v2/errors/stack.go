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
		FileName:     NewFileNameFromPC(pc),
		FunctionName: NewFunctionNameFromPC(pc),
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
	FileName     string
	FunctionName FunctionName
}

func (f Frame) String() string {
	return f.FileName + " " + f.FunctionName.WithPath
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

func NewFunctionNameFromPC(pc uintptr) FunctionName {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return FunctionName{
			WithPath:    "unknown",
			WithPackage: "unknown",
			Name:        "unknown",
		}
	}

	withPath := fn.Name()
	i := strings.LastIndex(withPath, "/")
	withPackage := withPath[i+1:]
	j := strings.Index(withPackage, ".")
	name := withPackage[j+1:]

	return FunctionName{
		WithPath:    withPath,
		WithPackage: withPackage,
		Name:        name,
	}
}

type FunctionName struct {
	WithPath    string
	WithPackage string
	Name        string
}

func callers(skip int) []uintptr {
	const depth = 32

	var pcs [depth]uintptr

	n := runtime.Callers(skip, pcs[:])
	st := pcs[0:n]

	return st
}
