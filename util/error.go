package util

import (
	"errors"
	"strconv"
	"strings"
)

//define application level error
type ErrorCode int

const (
	SYSTEM_ERROR   int = 0
	ADDR_FORBIDDEN     = 1
)

func ConstructError(str string, errCode int) error {
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(errCode))
	sb.WriteRune(':')
	sb.WriteString(str)
	return errors.New(sb.String())
}

func DeErrCode(err error) int {
	str := err.Error()
	strArr := strings.Split(str, ":")
	if errCodeExist(strArr[0]) {
		code, _ := strconv.Atoi(strArr[0])
		return code
	}
	return SYSTEM_ERROR
}
func errCodeExist(str string) bool {
	_, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	return true
}
