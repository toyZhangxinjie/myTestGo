package tonghuashun

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func readExcel(path, sheet string) ([][]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	return f.GetRows(sheet), nil
}

func LessId() {
	excel, err := readExcel("/Users/zhangxinjie/Downloads/有色网同花顺互换数据汇总-更新20210205.xlsx", "同花顺提供补充数据汇总")
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	db, err := readExcel("/Users/zhangxinjie/Downloads/cn_macro_index_basic_info.xlsx", "cn_macro_index_basic_info")
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	noExist := make([][]string, 0)
	for excelIdx, excelitem := range excel {
		if excelIdx == 0 {
			continue
		}
		idexist := false
		nameexist := false
		idcel := strings.Trim(excelitem[4], " ")
		namecel := strings.Trim(excelitem[5], " ")
		id := "S"
		for i := 0; i < 9-len(idcel); i++ {
			id += "0"
		}
		id += idcel
		for dbIdx, dbitem := range db {
			if dbIdx == 0 {
				continue
			}
			indicatorId := strings.Trim(dbitem[3], " ")
			indicatorName := strings.Trim(dbitem[4], " ")
			if id == indicatorId {
				idexist = true
			}
			if indicatorName == namecel {
				nameexist = true
			}
		}
		if !idexist || !nameexist {
			if !idexist && !nameexist {
				excelitem = append(excelitem, "指标ID和名称都不存在")
			} else if !idexist {
				excelitem = append(excelitem, "指标ID不存在")
			}else if !nameexist {
				excelitem = append(excelitem, "指标名称不存在")
			}
			noExist = append(noExist, excelitem)
		}
	}
	println("noExist len=", len(noExist))
	file := excelize.NewFile()
	sheetname := "Sheet1"
	sheetidx := file.NewSheet(sheetname)
	file.SetActiveSheet(sheetidx)
	for n, row := range noExist {
		i := n+1
		for ci, cel := range row {
			file.SetCellValue(sheetname, excelize.ToAlphaString(ci) + strconv.Itoa(i), cel)
		}
	}
	b, _ := file.WriteToBuffer()
	ioutil.WriteFile("/Users/zhangxinjie/Downloads/补充数据-中国宏观数据.xlsx", b.Bytes(), os.ModePerm)
	//println("count length=", count)
}