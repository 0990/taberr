package printer

import (
	"fmt"
)

func valueWrapperLua(str string) string {
	return fmt.Sprintf("\"%s\"", str)
}

type luaPrinter struct {
}

func (self *luaPrinter) Run(g *Global) *Stream {

	stream := NewStream()

	stream.Printf("-- Generated by github.com/0990/taberr\n")
	stream.Printf("-- Version: %s\n", g.Version)

	stream.Printf("\nlocal tab = {\n")

	printTableLua(g, stream)
	stream.Printf("\n\n")
	stream.Printf("}\n\n")

	if !genLuaIndexCode(g, stream) {
		return stream
	}

	// 生成枚举
	if !genLuaEnumCode(g, stream) {
		return stream
	}

	stream.Printf("\nreturn tab")

	return stream
}

func printTableLua(g *Global, stream *Stream) bool {

	stream.Printf("	%s = {\n", g.PackageName)

	// 遍历每一行
	for rowIndex, rowData := range g.Data {

		// 每一行开始
		stream.Printf("		{ ")
		stream.Printf("%s = %s", g.ErrTypeLabel, valueWrapperLua(rowData.ErrType))
		//stream.Printf("%s")
		stream.Printf(", ")
		stream.Printf("%s = %s", g.ErrMsgLabel, rowData.ErrMsg)
		//stream.Printf("")

		// 每一行结束
		stream.Printf(" 	}")

		if rowIndex < len(g.Data)-1 {
			stream.Printf(",")
		}
		stream.Printf("\n")
	}

	stream.Printf("	}")

	return true

}

// 收集需要构建的索引的类型
func genLuaEnumCode(g *Global, stream *Stream) bool {

	stream.Printf("\ntab.Enum = {\n")
	stream.Printf("	%s = {\n", g.EnumName)
	// 遍历字段
	for _, rowData := range g.Data {
		stream.Printf("		%s = %d,\n", rowData.ErrType, rowData.ErrID)
	}
	stream.Printf("	},\n")
	stream.Printf("}\n")

	return true

}

// 收集需要构建的索引的类型
func genLuaIndexCode(g *Global, stream *Stream) bool {

	mapperVarName := fmt.Sprintf("tab.%sBy%s", g.PackageName, g.ErrTypeLabel)

	stream.Printf("\n-- %s\n", g.ErrTypeLabel)
	stream.Printf("%s = {}\n", mapperVarName)
	stream.Printf("for _, rec in pairs(tab.%s) do\n", g.PackageName)
	stream.Printf("\t%s[rec.%s] = rec\n", mapperVarName, g.ErrTypeLabel)
	stream.Printf("end\n")
	return true
}

func init() {

	RegisterPrinter("lua", &luaPrinter1{})

}
