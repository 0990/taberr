package main

import (
	"flag"
	"fmt"
	"github.com/0990/taberr/printer"
	"github.com/tealeg/xlsx"
)

var (
	paramPackageName = flag.String("package", "emsg", "package name")
	paramEnumName    = flag.String("enum_name", "Err", "enum name")
	paramProtoOut    = flag.String("proto_out", "game_err.proto", "output protobuf define (*.proto)")
	paramLuaOut      = flag.String("lua_out", "game_err.lua", "output lua code (*.lua)")
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
	g.AddOutputType("proto", *paramProtoOut)
	g.AddOutputType("lua", *paramLuaOut)
	g.AddOutputType("lua1", "xujialong1.lua")
	if !GetXLSXData(g, fileName) {
		return
	}

	g.Print()
}

func GetXLSXData(g *printer.Global, fileName string) bool {
	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	//rowDatas := make([]printer.Data, 0)
	for _, sheet := range xlFile.Sheets {
		fmt.Println("sheet name:", sheet.Name)
		fmt.Println("rowCount:", len(sheet.Rows))

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
				case 2:
					rowData.ErrMsg = cell.String()
					validData = true
				}
			}
			if validData {
				g.Data = append(g.Data, rowData)
			}
		}
	}
	return true
}
