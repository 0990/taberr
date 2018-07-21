package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"github.com/0990/taberr/printer"
)


func main() {
	g:=printer.NewGlobal()
	g.Data =GetXLSXData()
	g.AddOutputType("proto", "xujialong.proto")
	g.Print()
}


func GetXLSXData()[]printer.Data{
	xlFile, err := xlsx.OpenFile("Item.xlsx")
	if err != nil {
		fmt.Println(err.Error())
	}
	rowDatas := make([]printer.Data, 0)
	for _, sheet := range xlFile.Sheets {
		fmt.Println("sheet name:", sheet.Name)
		fmt.Println("rowCount:", len(sheet.Rows))
		//每一行
		for _, row := range sheet.Rows {
			//每个单元
			rowData := printer.Data{}
			validData := false
			for i, cell := range row.Cells {
				switch i {
				case 0:
					intValue, err := cell.Int64()
					if err != nil {
						break
					}
					rowData.ErrID = int32(intValue)
				case 1:
					text := cell.String()
					if text == "" {
						break
					}
					rowData.ErrType = text
				case 2:
					rowData.ErrMsg = cell.String()
					validData = true
				}
			}
			if validData {
				rowDatas = append(rowDatas, rowData)
			}
		}
	}
	return rowDatas
}
