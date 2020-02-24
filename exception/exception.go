package exception

import (
	"fmt"
	"runtime/debug"
)

type Exception struct {
	Msg string
	Stack []byte
}

func NewException(code int, msg string) *Exception {
	return &Exception{Msg:msg, Stack: debug.Stack()}
}

func NewSysErrorException(errMsgLog string) *Exception {
	return &Exception{Msg: errMsgLog, Stack: debug.Stack()}
}

func ToError(err error) {
	if err == nil {
		return
	}
	excep, ok := err.(*Exception)
	if ok {
		fmt.Println(excep.Error())
	} else {
		fmt.Println(err.Error())
	}
}

func (e *Exception) Error() string {
	lineBytes := make([]byte, 0)
	lineNum := 1
	stackStr := "[clubcenter error begin] >>>>> ErrorMsg: " + e.Msg
	for _, b := range e.Stack {
		lineBytes = append(lineBytes, b)
		if b == 10 {
			if lineNum < 2 || lineNum > 5 {
				stackStr += string(lineBytes)
			}
			lineNum++
			lineBytes = make([]byte, 0)
		}
	}
	stackStr += "[clubcenter error end] <<<<<"
	return stackStr

	//stackStr := string(e.Stack)

	//loop := 1
	//for true {
	//	stackStr = stackStr[strings.Index(stackStr, "\n") + 1:]
	//	loop++
	//	if loop >= 6 {
	//		break
	//	}
	//}
	//fmt.Println(stackStr)

	//stackList := strings.Split(stackStr, "\n")
	//res := ""
	//for i, line := range stackList {
	//	if i > 0 && i < 5 {
	//		continue
	//	}
	//	res += line + "\n"
	//}
	//fmt.Println(res)
}

