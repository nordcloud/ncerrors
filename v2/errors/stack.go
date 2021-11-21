package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type stackTrace struct {
	Frames []frame
}

func (t stackTrace) stringStack() []string {
	var trace []string
	for _, f := range t.Frames {
		trace = append(trace, f.String())
	}

	return trace
}

func (t stackTrace) String() string {
	return strings.Join(t.stringStack(), "\n")
}

func newStackTrace(skip int) stackTrace {
	callers := callers(skip)

	var frames []frame
	for _, c := range callers {
		frames = append(frames, newFrame(c))
	}

	return stackTrace{
		Frames: frames,
	}
}

func newFrame(f uintptr) frame {
	pc := f - 1

	return frame{
		systemFileName: newFileNameFromPC(pc),
		functionName:   newFunctionFromPC(pc),
	}
}

func newFileNameFromPC(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	file, _ := fn.FileLine(pc)

	return file
}

type frame struct {
	// not very useful, since it is system dependant
	systemFileName string

	functionName function
}

func (f frame) String() string {
	return f.functionName.String()
}

func newFunctionFromPC(pc uintptr) function {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return function{
			nameWithPath:    "unknown",
			nameWithPackage: "unknown",
			name:            "unknown",
		}
	}

	withPath := fn.Name()
	i := strings.LastIndex(withPath, "/")
	withPackage := withPath[i+1:]
	j := strings.Index(withPackage, ".")
	name := withPackage[j+1:]

	_, lineNum := fn.FileLine(pc)

	return function{
		nameWithPath:    withPath,
		nameWithPackage: withPackage,
		name:            name,
		lineNumber:      lineNum,
	}
}

type function struct {
	nameWithPath    string
	nameWithPackage string
	name            string
	lineNumber      int
}

func (f function) String() string {
	return fmt.Sprintf("%s:%d", f.nameWithPath, f.lineNumber)
}

func callers(skip int) []uintptr {
	const depth = 32

	var pcs [depth]uintptr

	n := runtime.Callers(skip, pcs[:])
	st := pcs[0:n]

	return st
}
