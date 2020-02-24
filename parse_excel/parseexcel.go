package parse_excel

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io/ioutil"
	"os"
	"strings"
)


func ReadExcel(path, sheet string) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Errorf(err.Error())
	}

	rows := f.GetRows(sheet)

	sql1, sql2 := "",""

	for rowIdx, row := range rows {
		if rowIdx == 0 {
			for i, title := range row {
				fmt.Print(i, title, "\t")
			}
			fmt.Println()
		} else {
			s1,s2 := updateUserCardIndustry(row)
			sql1 += s1
			sql2 += s2
		}
	}
	fmt.Println("excel 共", len(rows), "行")

	sql := sql1 + sql2
	err = ioutil.WriteFile("/Users/zhangxinjie/Downloads/update_card_industry.sql",[]byte(sql), os.ModePerm)
	if err != nil {
		fmt.Errorf(err.Error())
	}
}

func updateUserCardIndustry(row []string) (string,string) {
	sql := ""
	sql2 := ""

	var cardId,companyType,industryNames,mainProduct string
	for i, cell := range row {
		val := strings.Trim(cell," ")
		switch i {
		case 0:
			cardId = val
		case 1:
			// userId = val
		case 2:
			companyType = val
		case 3:
			if val != "" {
				names := strings.Split(val,",")
				nameList := make([]string,0)
				for _, name := range names {
					nameList = append(nameList,"'" + strings.Trim(name," ") + "'")
				}
				industryNames = strings.Join(nameList, ",")
			}
		case 4:
			mainProduct = val
		}
	}

	sql += "UPDATE user_card SET company_type='" + companyType + "',main_product='" + mainProduct + "' WHERE id=" + cardId + ";\n"
	if industryNames != "" {
		sql2 += "INSERT INTO industry_card(industry_id,card_id) SELECT id," + cardId + " AS card_id FROM industry WHERE industry_name IN (" + industryNames + ");\n"
	}

	return sql, sql2
}