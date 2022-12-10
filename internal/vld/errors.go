package vld

import (
	"fmt"
)

type ErrorType int

const (
	ErrorTypeNil = ErrorType(iota)
	ErrorTypeMissRequired
	ErrorTypeNumOutOfRange
	ErrorTypeLengthOutOfRange
	ErrorTypeNotMatchRegexp
	ErrorTypeBadTimeValue
	ErrorTypeCanNotCastToNum
	ErrorTypeCanNotCastToBool
	ErrorTypeUserDefined
)

var (
	errReasonStrings []string
)

func (er ErrorType) String() string { return errReasonStrings[er] }

func init() {
	errReasonStrings = []string{
		"Undefined",
		"MissRequired",
		"NumOutOfRange",
		"LengthOutOfRange",
		"NotMatchRegexp",
		"BadTimeValue",
		"CanNotCastToNum",
		"CanNotCastToBool",
	}
}

type Error struct {
	Type ErrorType
	Rule *Rule
	msg  string
}

func MakeError(msg string) *Error {
	return &Error{
		msg:  msg,
		Type: ErrorTypeUserDefined,
	}
}

func (err *Error) Error() string {
	if err.Type == ErrorTypeUserDefined {
		return fmt.Sprintf("%s %s", err.msg, err.Rule.Name)
	}
	return fmt.Sprintf("%s %s", err.Type, err.Rule.Name)
}
