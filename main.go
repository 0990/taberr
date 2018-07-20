package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
)

type Data struct {
	errID   int32
	errType string
	errMsg  string
}

func main() {
	xlFile, err := xlsx.OpenFile("Item.xlsx")
	if err != nil {
		fmt.Println(err.Error())
	}
	rowDatas := make([]Data, 0)
	for _, sheet := range xlFile.Sheets {
		fmt.Println("sheet name:", sheet.Name)
		fmt.Println("rowCount:", len(sheet.Rows))
		//每一行
		for _, row := range sheet.Rows {
			//每个单元
			rowData := Data{}
			validData := false
			for i, cell := range row.Cells {
				switch i {
				case 0:
					intValue, err := cell.Int64()
					if err != nil {
						break
					}
					rowData.errID = int32(intValue)
				case 1:
					text := cell.String()
					if text == "" {
						break
					}
					rowData.errType = text
				case 2:
					rowData.errMsg = cell.String()
					validData = true
				}
			}
			if validData {
				rowDatas = append(rowDatas, rowData)
			}
		}

	}
	fmt.Println(rowDatas)
}
