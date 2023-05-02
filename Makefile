GOCMD = go
GOTEST = $(GOCMD) test

GOOS ?= linux
GOARCH ?= amd64
export GO111MODULE ?= on
export GOPROXY ?= direct
export GOSUMDB ?= off
LDFLAGS ?= -s -w -extldflags "-static"
export CGO_ENABLED ?= 0

test:
	TZ=UTC $(GOTEST) ./... -count=1 -p 1 -cover
