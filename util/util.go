package util

import (
	"encoding/json"
	"fmt"
	"github.com/mozillazg/go-pinyin"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/levigross/grequests"
)

func NewRequestOptions() *grequests.RequestOptions {
	return &grequests.RequestOptions{
		DialTimeout:         5 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
		RequestTimeout:      5 * time.Second,
		InsecureSkipVerify:  true,
	}
}

type HttpResponse struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Data *json.RawMessage `json:"data"`
}

func HttpPostCall(url string, params map[string]string, res interface{}) error {
	op := NewRequestOptions()
	op.Data = params
	resp, err := grequests.Post(url, op)
	if err != nil {
		return err
	}

	if !resp.Ok {
		return fmt.Errorf("Response failed with status code: %d.", resp.StatusCode)
	}

	if res != nil {
		return resp.JSON(res)
	}
	return nil
}
func HttpGetCall(url string, params map[string]string, cooikes []*http.Cookie, res interface{}) error {
	op := NewRequestOptions()
	op.Params = params
	op.Cookies = cooikes
	resp, err := grequests.Get(url, op)
	if err != nil {
		return err
	}

	if !resp.Ok {
		return fmt.Errorf("Response failed with status code: %d.", resp.StatusCode)
	}

	if res != nil {
		return resp.JSON(res)
	}
	return nil
}

func GetRandNum(begin, end int) int {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	if (begin == 0 && end == 0) || begin == end {
		return 0
	}

	if begin > end {
		temp := begin
		begin = end
		end = temp
	}
	return r.Intn(end-begin) + begin
}

func IntSliceToStr(slice []int, sep string) string {
	var strSlice []string
	for _, v := range slice {
		if v != 0 {
			strSlice = append(strSlice, strconv.Itoa(v))
		}
	}
	return strings.Join(strSlice, sep)
}

func IntSliceToInt64Slice(slice []int) []int64 {
	var list []int64
	for _, i := range slice {
		list = append(list, int64(i))
	}
	return list
}
func Int64SliceToIntSlice(slice []int64) []int {
	var list []int
	for _, i := range slice {
		list = append(list, int(i))
	}
	return list
}

func ToPinYin(str string) string {
	py := pinyin.Pinyin(str, pinyin.NewArgs())
	res := ""
	for _, i := range py {
		res += strings.Join(i, "")
	}
	return res
}

const (
	CELLPHONE_REGULAR = "^[0-9]{11}$"
	EMAIL_REGULAR     = `^[a-z0-9A-Z]+([\-_\.][a-z0-9A-Z]+)*@([a-z0-9A-Z]+(-[a-z0-9A-Z]+)*\.)+[a-zA-Z]+$`
)

func Validate(s string, regular string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(s)
}

func FirstIndexOf(str, s string) int {
	idx := strings.Index(str, s)
	if idx >= 0 {
		pre := []byte(str)[0:idx]
		rs := []rune(string(pre))
		idx = len(rs)
	}
	return idx
}

func LastIndexOf(str, s string) int {
	idx := strings.LastIndex(str, s)
	if idx >= 0 {
		pre := []byte(str)[0:idx]
		rs := []rune(string(pre))
		idx = len(rs)
	}
	return idx
}

func SubString(str string, start, end int) string {
	return string([]rune(str)[start:end])
}
