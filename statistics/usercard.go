package main

import (
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"myTestGo/util"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
统计有色网员工递名片
*/
const (
	user   = ""
	pwd    = ""
	ip     = ""
	port   = 3306
	dbName = ""
)

var db *sqlx.DB

func init() {
	mysqlLink := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local", user, pwd, ip, port, dbName)
	var err error
	db, err = sqlx.Open("mysql", mysqlLink)
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Second * 300)
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("db init successful")
}

type Employee struct {
	Number    string
	WorkCity  string
	Name      string
	Cellphone string
	Email     string
	UserId    int64
}
type result struct {
	UserId int64 `json:"user_id"`
}

func main() {
	f, err := excelize.OpenFile("/Users/zhangxinjie/Downloads/有色网员工手机信息.xlsx")
	if err != nil {
		fmt.Errorf(err.Error())
	}

	employees := make([]Employee, 0)
	rows := f.GetRows("岗位表")
	for i, row := range rows {
		if i == 0 {
			continue
		}
		number := strings.Trim(row[0], " ")
		workCity := strings.Trim(row[1], " ")
		name := strings.Trim(row[2], " ")
		cellphone := strings.Trim(row[15], " ")
		email := strings.Trim(row[16], " ")
		emp := Employee{
			Number:    number,
			WorkCity:  workCity,
			Name:      name,
			Cellphone: cellphone,
			Email:     email,
			UserId:    0,
		}
		if cellphone == "" {
			fmt.Println("手机号为空 >> ", name, cellphone, email)
			employees = append(employees, emp)
			continue
		}
		var response util.HttpResponse
		err = util.HttpGetCall("https://platform.smm.cn/usercenter/5.0/get_user_id_by_cellphone",
			map[string]string{"cellphone": cellphone}, nil, &response)
		if err != nil {
			panic(err.Error())
		}
		if response.Code != 0 {
			panic(response.Msg)
		}
		var res result
		err = json.Unmarshal(*response.Data, &res)
		if err != nil {
			panic(err.Error())
		}
		emp.UserId = res.UserId
		employees = append(employees, emp)
		time.Sleep(time.Millisecond * 50)
	}

	userIds := make([]int64, 0)
	for _, emp := range employees {
		if emp.UserId > 0 {
			userIds = append(userIds, emp.UserId)
		}
	}
	cardIds, err := getUserCard(userIds)
	if err != nil {
		panic(err.Error())
	}
	cardIdList := make([]int64, 0)
	for _, cardid := range cardIds {
		cardIdList = append(cardIdList, cardid)
	}
	begin := time.Date(2021,8,2,0,0,0,0,time.Local)
	end := time.Date(2021,8,8,23,59,59,0,time.Local)
	dateStr := begin.Format("01/02") + "-" + end.Format("01/02")
	shareSelfTotal, err := getShareSelf(begin.Unix(),end.Unix(),userIds)
	if err != nil {
		panic(err.Error())
	}
	submitTotal, err := getSubmitCard(begin.Unix(),end.Unix(),userIds)
	if err != nil {
		panic(err.Error())
	}
	friendTotal, err := getFriends(begin.Unix(),end.Unix(),userIds)
	if err != nil {
		panic(err.Error())
	}
	lookMeTotal, err := getLookMe(begin.Unix(),end.Unix(),cardIdList)
	if err != nil {
		panic(err.Error())
	}

	file := excelize.NewFile()
	sheetname := "名片统计"
	file.SetActiveSheet(file.NewSheet(sheetname))

	title := []string{
		"姓名","手机号","用户ID","名片ID",dateStr+"分享名片",dateStr+"递名片",dateStr+"新增好友数",dateStr+"名片被看数",
	}
	for i, s := range title {
		file.SetCellValue(sheetname,excelize.ToAlphaString(i)+"1", s)
	}
	for i, emp := range employees {
		rownum := strconv.Itoa(i+2)
		file.SetCellValue(sheetname,excelize.ToAlphaString(0)+rownum, emp.Name)
		file.SetCellValue(sheetname,excelize.ToAlphaString(1)+rownum, emp.Cellphone)
		file.SetCellValue(sheetname,excelize.ToAlphaString(2)+rownum, emp.UserId)
		if cardid, ok := cardIds[emp.UserId]; ok {
			file.SetCellValue(sheetname,excelize.ToAlphaString(3)+rownum,cardid)
		} else {
			file.SetCellValue(sheetname,excelize.ToAlphaString(3)+rownum,"")
		}
		if shareTotal, ok := shareSelfTotal[emp.UserId]; ok {
			file.SetCellValue(sheetname,excelize.ToAlphaString(4)+rownum,shareTotal)
		} else {
			file.SetCellValue(sheetname,excelize.ToAlphaString(4)+rownum,"")
		}
		if submit, ok := submitTotal[emp.UserId]; ok {
			file.SetCellValue(sheetname,excelize.ToAlphaString(5)+rownum,submit)
		} else {
			file.SetCellValue(sheetname,excelize.ToAlphaString(5)+rownum,"")
		}
		if friend, ok := friendTotal[emp.UserId]; ok {
			file.SetCellValue(sheetname,excelize.ToAlphaString(6)+rownum,friend)
		} else {
			file.SetCellValue(sheetname,excelize.ToAlphaString(6)+rownum,"")
		}
		if look, ok := lookMeTotal[cardIds[emp.UserId]]; ok {
			file.SetCellValue(sheetname,excelize.ToAlphaString(7)+rownum,look)
		} else {
			file.SetCellValue(sheetname,excelize.ToAlphaString(7)+rownum,"")
		}
	}
	b, _ := file.WriteToBuffer()
	ioutil.WriteFile("/Users/zhangxinjie/Downloads/有色网员工名片统计信息.xlsx",b.Bytes(), os.ModePerm)

	fmt.Println("结束！！！")
}

// 获取分享自己名片统计
func getShareSelf(begin, end int64, userIds []int64) (map[int64]int, error) {
	sql := "SELECT master_user_id AS user_id,COUNT(*) AS count FROM submit_card_record WHERE 1=1 "
	params := make([]interface{}, 0)
	if len(userIds) > 0 {
		sql += " AND master_user_id IN (?) "
		params = append(params, userIds)
	}
	sql += " AND submit_type=2 AND share_type='self_card' "
	if begin > 0 {
		sql += " AND create_time >= ? "
		params = append(params, begin)
	}
	if end > 0 {
		sql += " AND create_time <= ? "
		params = append(params, end)
	}
	sql += " GROUP BY master_user_id"
	type data struct {
		UserId int64 `db:"user_id"`
		Count  int   `db:"count"`
	}
	list := make([]data, 0)
	sql, args, err := sqlx.In(sql, params...)
	if err != nil {
		return nil, err
	}
	sql = db.Rebind(sql)
	err = db.Select(&list, sql, args...)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]int)
	for _, d := range list {
		result[d.UserId] = d.Count
	}
	return result,nil
}

// 获取递名片统计
func getSubmitCard(begin, end int64, userIds []int64) (map[int64]int, error) {
	sql := "SELECT master_user_id AS user_id,COUNT(*) AS count FROM submit_card_record WHERE 1=1 "
	params := make([]interface{}, 0)
	if len(userIds) > 0 {
		sql += " AND master_user_id IN (?) "
		params = append(params, userIds)
	}
	sql += " AND submit_type=1 "
	if begin > 0 {
		sql += " AND create_time >= ? "
		params = append(params, begin)
	}
	if end > 0 {
		sql += " AND create_time <= ? "
		params = append(params, end)
	}
	sql += " GROUP BY master_user_id"
	type data struct {
		UserId int64 `db:"user_id"`
		Count  int   `db:"count"`
	}
	list := make([]data, 0)
	sql, args, err := sqlx.In(sql, params...)
	if err != nil {
		return nil, err
	}
	sql = db.Rebind(sql)
	err = db.Select(&list, sql, args...)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]int)
	for _, d := range list {
		result[d.UserId] = d.Count
	}
	return result,nil
}

// 获取好友数统计
func getFriends(begin, end int64, userIds []int64) (map[int64]int, error) {
	sql := "SELECT user_master AS user_id,COUNT(*) AS count FROM user_card_clip WHERE 1=1 "
	params := make([]interface{}, 0)
	if len(userIds) > 0 {
		sql += " AND user_master IN (?) "
		params = append(params, userIds)
	}
	sql += " AND relationship=1 "
	if begin > 0 {
		sql += " AND add_time >= ? "
		params = append(params, begin)
	}
	if end > 0 {
		sql += " AND add_time <= ? "
		params = append(params, end)
	}
	sql += " GROUP BY user_master"
	type data struct {
		UserId int64 `db:"user_id"`
		Count  int   `db:"count"`
	}
	list := make([]data, 0)
	sql, args, err := sqlx.In(sql, params...)
	if err != nil {
		return nil, err
	}
	sql = db.Rebind(sql)
	err = db.Select(&list, sql, args...)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]int)
	for _, d := range list {
		result[d.UserId] = d.Count
	}
	return result,nil
}

// 获取名片被看数统计
func getLookMe(begin, end int64, cardIds []int64) (map[int64]int, error) {
	sql := "SELECT seen_card_id AS card_id,COUNT(*) AS count FROM look_me WHERE 1=1 "
	params := make([]interface{}, 0)
	if len(cardIds) > 0 {
		sql += " AND seen_card_id IN (?) "
		params = append(params, cardIds)
	}
	if begin > 0 {
		sql += " AND look_time >= ? "
		params = append(params, begin)
	}
	if end > 0 {
		sql += " AND look_time <= ? "
		params = append(params, end)
	}
	sql += " GROUP BY seen_card_id"
	type data struct {
		CardId int64 `db:"card_id"`
		Count  int   `db:"count"`
	}
	list := make([]data, 0)
	sql, args, err := sqlx.In(sql, params...)
	if err != nil {
		return nil, err
	}
	sql = db.Rebind(sql)
	err = db.Select(&list, sql, args...)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]int)
	for _, d := range list {
		result[d.CardId] = d.Count
	}
	return result,nil
}

// 获取用户对应的名片ID
func getUserCard(userIds []int64) (map[int64]int64, error) {
	if len(userIds) <= 0 {
		return make(map[int64]int64), nil
	}
	sql := "SELECT user_id,id FROM user_card WHERE user_id IN (?)"
	sql, args, err := sqlx.In(sql,userIds)
	if err != nil {
		return nil, err
	}
	sql = db.Rebind(sql)
	type data struct {
		UserId int64 `db:"user_id"`
		CardId int64 `db:"id"`
	}
	list := make([]data, 0)
	err = db.Select(&list,sql,args...)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]int64)
	for _, d := range list {
		result[d.UserId] = d.CardId
	}
	return result, nil
}