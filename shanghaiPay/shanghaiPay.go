package shanghaiPay

import (
	"fmt"
	"os"
	"strings"
)

const (
	pfxPath string = "/Users/zhangxinjie/Downloads/demo/API_DEMO/315310018000398.pfx"
	pfxPwd string = "123123"
	errCode int = 500
)

var (
	file *os.File
	eofFound bool
	limit int64
)

func SignData(data string)  {
	signInfo := make([]string, 2)

	fmt.Println(signInfo)
	initSignCertAndKey(pfxPath, pfxPwd)
}

func initSignCertAndKey(pfxFileName, password string) int {
	if len(strings.Trim(password, " ")) <= 0 {
		println("Password cannot be empty")
		return errCode
	}

	file, err := os.Open(pfxPath)
	defer file.Close()

	if err != nil {
		println(err.Error())
		return errCode
	}
	bytes := make([]byte, 1)
	_, err = file.Read(bytes)
	if err != nil {
		println(err.Error())
		return errCode
	}

	tag := int(bytes[0])
	if tag != 48 {
		println("stream does not represent a PKCS12 key store")
		return errCode
	}

	readObject(tag)


	return 0
}

func readObject(tag int) int {
	eofFound = false
	limit = 2147483647

	if tag == -1 {
		if eofFound {
			println("attempt to read past end of file")
			return errCode
		} else {
			eofFound = true
			return 0
		}
	} else {
		tagNo := 0
		fmt.Println(tagNo)
		if (tag & 128) != 0 {
			tagNo, code := readTagNumber(tag)
			if code == errCode {
				return code
			}
			println(tagNo)
		}
	}
	return 0
}

func readTagNumber(tag int) (int, int) {
	tagNo := tag & 31
	if tagNo == 31 {
		bytes := make([]byte, 1)
		_, err := file.Read(bytes)
		b := int(bytes[0])
		if err != nil {
			println(err.Error())
			return 0, errCode
		}
		tagNo = 0
		for b >= 0 && (b & 128) != 0 {
			tagNo |= b & 127
			tagNo <<= 7

			_, err = file.Read(bytes)
			if err != nil {
				println(err.Error())
				return 0, errCode
			}
			b = int(bytes[0])
			println(b)
		}
		if b < 0 {
			eofFound = true
			println("EOF found inside tag value.")
			return 0, errCode
		}
		tagNo |= b & 127
	}
	return tagNo, 0
}