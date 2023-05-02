// Copyright 2023 Nordcloud Oy or its affiliates. All Rights Reserved.

package errors

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

var (
	regexFuncName = regexp.MustCompile(`(([^/]+/)+)?([^/.]+)((\.[^/.]+)+)?`)
)

// getStackTraces returns custom-formatted and raw (in the form of program counters) stack trace
// for the purpose of initializing NCError struct
func getStackTraces() ([]string, *stack) {
	var formattedStack []string
	callStack := *callers()
	st := callStack[:len(callStack)-1]
	for _, f := range st {
		frame := frame(f)
		formattedStack = append(formattedStack, frame.formatContext())
	}

	return formattedStack, &callStack
}

// GetTrace return the simplified stack trace in the format file_name(func_name):line. It also contains the current goroutine entrypoint.
func GetTrace() []string {
	stack, _ := getStackTraces()

	return stack
}

type frame uintptr

func (f frame) pc() uintptr { return uintptr(f) - 1 }
func (f frame) getContext() (fileName, funcName string, lineNum int) {
	fn := runtime.FuncForPC(f.pc())
	if fn == nil {
		return "", "", 0
	}

	funcNamePath := fn.Name()
	pathParts := strings.Split(funcNamePath, "/")
	funcPkg := pathParts[len(pathParts)-1]
	if len(pathParts) > 1 {
		pathParts = pathParts[:len(pathParts)-1]
	}
	pkgParts := strings.SplitN(funcPkg, ".", 2)
	if len(pkgParts) < 2 {
		return "", "", 0
	}
	funcPkg = strings.Join(pathParts, "/") + "/" + pkgParts[0]
	funcName = pkgParts[1]

	var pos int
	fileName, lineNum = fn.FileLine(f.pc())
	pos = strings.LastIndex(fileName, "/")
	if pos >= 0 {
		fileName = fileName[pos+1:]
	}
	fileName = fmt.Sprintf("%s/%s", funcPkg, fileName)

	return
}

func (f frame) formatContext() string {
	fileName, funcName, line := f.getContext()
	return fmt.Sprintf("%s(%s):%d", fileName, funcName, line)
}

// GetRuntimeContext returns function name and code line.
func GetRuntimeContext() (fileName, funcName string, line int) {
	st := *callers()
	frame := frame(st[1])
	fileName, funcName, line = frame.getContext()
	return
}

type stack []uintptr

func (s *stack) StackTrace() errors.StackTrace {
	f := make([]errors.Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = errors.Frame((*s)[i])
	}
	return f
}

func callers() *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}
