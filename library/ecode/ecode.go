package ecode

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"sync/atomic"
)

var (
	_messages atomic.Value // NOTE: stored map[int]string
)

// Register register ecode message map.
func Register(cm map[int]string) {
	_messages.Store(cm)
}

type ECode interface {
	// sometimes Error return Code in string form
	// NOTE: don't use Error in monitor report even it also work for now
	Error() string
	// Code get error code.
	Code() int
	// Message get code message.
	Message() string
}

type ecode struct {
	code    int
	message string
}

func (e *ecode) Error() string {
	return e.message
}

func (e *ecode) Code() int {
	return e.code
}

func (e *ecode) Message() string {
	if cm, ok := _messages.Load().(map[int]string); ok {
		if msg, ok := cm[e.Code()]; ok {
			return msg
		}
	}
	if e.message == "" {
		return strconv.Itoa(e.Code())
	}
	return e.message
}

func Errorf(code int, args ...interface{}) *ecode {
	message := ""
	if len(args) > 0 {
		message = fmt.Sprintf(strings.Repeat("%v ", len(args)), args...)
	}
	return &ecode{code: code, message: message}
}

// Cause cause from error to ecode.
func Cause(e error) ECode {
	if e == nil {
		return &ecode{code: 0}
	}
	ec, ok := errors.Cause(e).(ECode)
	if ok {
		return ec
	}
	return &ecode{code: 500, message: e.Error()}
}
