package main

import (
	"flag"
	"fmt"
	"github.com/0990/taberr/printer"
	"github.com/tealeg/xlsx"
	"os"
)

var (
	paramPackageName = flag.String("package", "emsg", "package name")
	paramEnumName    = flag.String("enum_name", "Err", "enum name")
	paramProtoOut    = flag.String("proto_out", "", "output protobuf define (*.proto)")
	paramLuaOut      = flag.String("lua_out", "", "output lua code (*.lua)")
)

func main() {
	flag.Parse()

	var fileName string
	for _, v := range flag.Args() {
		fileName = v
	}

	g := printer.NewGlobal()
	g.PackageName = *paramPackageName
	g.EnumName = *paramEnumName
	if *paramProtoOut != "" {
		g.AddOutputType("proto", *paramProtoOut)
	}

	if *paramLuaOut != "" {
		g.AddOutputType("lua", *paramLuaOut)
	}

	if !GetXLSXData(g, fileName) {
		os.Exit(1)
	}

	g.Print()
}

func GetXLSXData(g *printer.Global, fileName string) bool {
	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		fmt.Printf("file('%s') read error,%s", fileName, err.Error())
		return false
	}
	g.FileName = fileName
	for _, sheet := range xlFile.Sheets {
		fmt.Printf("read sheet:%s,rowCount:%d\n", sheet.Name, len(sheet.Rows))

		g.ErrIDLabel = sheet.Rows[0].Cells[0].String()
		g.ErrTypeLabel = sheet.Rows[0].Cells[1].String()
		g.ErrMsgLabel = sheet.Rows[0].Cells[2].String()
		//每一行
		for rowIndex, row := range sheet.Rows {
			if rowIndex == 0 {
				continue
			}
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
					validData = true
				case 2:
					rowData.ErrMsg = cell.String()
				}
			}
			if validData {
				if g.CheckErrIDRepeate(rowData.ErrID) {
					fmt.Printf("errID repeat:%d row:%d,sheet:%s", rowData.ErrID, rowIndex+1, sheet.Name)
					return false
				}
				if g.CheckErrTypeRepeate(rowData.ErrType) {
					fmt.Printf("errType repeat:%s row:%d sheet:%s", rowData.ErrType, rowIndex+1, sheet.Name)
					return false
				}
				g.Data = append(g.Data, rowData)
			}
		}
	}
	return true
}
