package main

import (
	"fmt"
	"myTestGo/exception"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case error:
				println(exception.NewException(500, x.Error()).Error())
			case exception.Exception:
				println(x.Error())
			default:
				println("Unknown Panic")
			}
		}
	}()
	var err error

	timeFormat := "2006-01-02 15:04:05"
	now := time.Now()
	fmt.Println("Go Version:", runtime.Version(), " | ", now.Format(timeFormat), " | ", now.Unix())

	startTime := time.Date(now.Year(),now.Month(),now.Day(),0,0,0,0,time.Local)
	endTime := time.Date(now.Year(),now.Month(),now.Day(),23,59,59,0,time.Local)

	fmt.Println("今天的开始和结束时间戳：", startTime.Unix(), endTime.Unix())
	fmt.Println()

	// 帖子迁移
	//ClubPostMigrate()

	// 加密
	//shanghaiPay.SignData("")


	// 删除帖子索引
	//delpost.DoDel(1, 100)

	//canvas.CanvasDrawImg()

	//img.ImageDraw()

	// 生成postman markdown
	//err = postmantomd.PostManJSONToMarddown("/Users/zhangxinjie/Downloads/商机2.0.postman_collection.json","/Users/zhangxinjie/Downloads/clubcenter_2.0.md")
	//err = postmantomd.PostManJSONToMarddown("/Users/zhangxinjie/Downloads/通讯录card_center.postman_collection.json","/Users/zhangxinjie/Downloads/cardcenter-2020.md")
	if err != nil {
		exception.ToError(err)
	}

	// 读取excel
	//parse_excel.ReadExcel("/Users/zhangxinjie/Documents/企业微信/WXWork Files/File/2019-12/user_card的副本.xlsx","user_card")

	//tonghuashun.LessId()

}




