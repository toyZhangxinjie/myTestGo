package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	"strings"
	"time"
)


func Yarns() {
	//http://www.czce.com.cn/cn/DFSStaticFiles/Future/2021/20210610/FutureDataDailyCY.txt
	now := time.Now().In(time.Local)
	basePath := "http://www.czce.com.cn/cn/DFSStaticFiles/Future"

	dayArea := int64(365)
	for i := 0; i < 365; i++ {
		idxDay := now.AddDate(0,0,-i)

		url := fmt.Sprintf("%s/%v/%s/FutureDataDaily%s.txt", basePath,idxDay.Year(),idxDay.Format("20060102"),"CY")
		fmt.Println(url)
		content, err := download(url)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if content == nil || len(content) <= 0 {
			continue
		}

		stmt := fmt.Sprintf(`INSERT INTO yarns_czce(date,code,zuo_jie_suan,jin_kai_pan,zui_gao_jia,zui_di_jia,
jin_shou_pan,jin_jie_suan,zhang_die1,zhang_die2,cheng_jiao_liang,chi_cang_liang,zeng_jian_liang,cheng_jiao_e,jiao_ge_jie_suan_jia)
 VALUES`)
		params := make([]interface{}, 0)
		date := idxDay.Format("2006-01-02")
		for _, s := range content {
			s = strings.Trim(s," ")
			if s == "" || !strings.Contains(s, "|") || strings.HasPrefix(s,"合约代码") || strings.HasPrefix(s,"品种月份") {
				fmt.Println(s)
				continue
			}
			stmt += fmt.Sprintf(`(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?),`)
			col := strings.Split(s,"|")
			code := strings.Trim(col[0]," ")
			params = append(params, date, code)

			parseNum := func(str string) float64 {
				str = strings.Trim(str," ")
				if len(str) <= 0 {
					return 0
				} else {
					num, err := strconv.ParseFloat(strings.ReplaceAll(str,",",""),64)
					if err != nil {
						panic(err.Error())
					}
					return num
				}
			}
			params = append(params,
				parseNum(col[1]),
				parseNum(col[2]),
				parseNum(col[3]),
				parseNum(col[4]),
				parseNum(col[5]),
				parseNum(col[6]),
				parseNum(col[7]),
				parseNum(col[8]),
				parseNum(col[9]),
				parseNum(col[10]),
				parseNum(col[11]),
				parseNum(col[12]),
				parseNum(col[13]),
			)
		}
		stmt = stmt[:len(stmt)-1]
		err = Add(stmt, params)
		if err != nil {
			fmt.Printf(err.Error())
			println(stmt)
			for _, s := range content {
				println(s)
			}
			return
		}

		idxDay = idxDay.AddDate(0,0,-1)
		if (now.Unix() - idxDay.Unix()) / (60 * 60 * 24) >= dayArea {
			break
		}
	}

}

func Cotton() {
	//http://www.czce.com.cn/cn/DFSStaticFiles/Future/2021/20210412/FutureDataDailyCF.txt
	now := time.Now().In(time.Local)
	basePath := "http://www.czce.com.cn/cn/DFSStaticFiles/Future"

	dayArea := int64(365)
	for i := 0; i < 365; i++ {
		idxDay := now.AddDate(0,0,-i)

		url := fmt.Sprintf("%s/%v/%s/FutureDataDaily%s.txt", basePath,idxDay.Year(),idxDay.Format("20060102"),"CF")
		fmt.Println(url)
		content, err := download(url)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if content == nil || len(content) <= 0 {
			continue
		}

		stmt := fmt.Sprintf(`INSERT INTO cotton_czce(date,code,zuo_jie_suan,jin_kai_pan,zui_gao_jia,zui_di_jia,
jin_shou_pan,jin_jie_suan,zhang_die1,zhang_die2,cheng_jiao_liang,chi_cang_liang,zeng_jian_liang,cheng_jiao_e,jiao_ge_jie_suan_jia)
 VALUES`)
		params := make([]interface{}, 0)
		date := idxDay.Format("2006-01-02")
		for _, s := range content {
			s = strings.Trim(s," ")
			if s == "" || !strings.Contains(s, "|") || strings.HasPrefix(s,"合约代码") || strings.HasPrefix(s,"品种月份") {
				fmt.Println(s)
				continue
			}
			stmt += fmt.Sprintf(`(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?),`)
			col := strings.Split(s,"|")
			code := strings.Trim(col[0]," ")
			params = append(params, date, code)

			parseNum := func(str string) float64 {
				str = strings.Trim(str," ")
				if len(str) <= 0 {
					return 0
				} else {
					num, err := strconv.ParseFloat(strings.ReplaceAll(str,",",""),64)
					if err != nil {
						panic(err.Error())
					}
					return num
				}
			}
			params = append(params,
				parseNum(col[1]),
				parseNum(col[2]),
				parseNum(col[3]),
				parseNum(col[4]),
				parseNum(col[5]),
				parseNum(col[6]),
				parseNum(col[7]),
				parseNum(col[8]),
				parseNum(col[9]),
				parseNum(col[10]),
				parseNum(col[11]),
				parseNum(col[12]),
				parseNum(col[13]),
			)
		}
		stmt = stmt[:len(stmt)-1]
		err = Add(stmt, params)
		if err != nil {
			fmt.Printf(err.Error())
			println(stmt)
			for _, s := range content {
				println(s)
			}
			return
		}

		idxDay = idxDay.AddDate(0,0,-1)
		if (now.Unix() - idxDay.Unix()) / (60 * 60 * 24) >= dayArea {
			break
		}
	}

}

func download(filepath string) ([]string, error) {
	resp, err := http.Get(filepath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, nil
	}

	result := make([]string, 0)
	br := bufio.NewReader(resp.Body)
	for {
		line, _, err := br.ReadLine()
		if err != nil {
			break
		}
		result = append(result, string(line))
	}
	return result, nil
}

func Add(sqlStmt string, params []interface{}) error {
	user := ""
	pwd := ""
	ip := ""
	port := 3306
	dbName := ""
	mysqlLink := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True",user,pwd,ip,port,dbName)
	db, err := sql.Open("mysql",mysqlLink)
	if err != nil {
		return err
	}

	_, err = db.Exec(sqlStmt, params...)
	if err != nil {
		return err
	}
	return nil
}


func main() {
	Cotton()
}