package scanner

import (
	"path"
	"strings"

	"github.com/davyxu/pbmeta"
	"github.com/davyxu/tabtoy/proto/tool"
	"github.com/golang/protobuf/proto"
	"github.com/tealeg/xlsx"
)

type File struct {
	*xlsx.File
	Sheets   []*Sheet // 多个sheet
	SheetMap map[string]*Sheet
	descpool *pbmeta.DescriptorPool // 协议描述池
	Name     string                 // 电子表格的名称(文件名无后缀)
	FileName string
}

// 制作表名
func makeTableName(filename string) string {
	baseName := path.Base(filename)
	return strings.TrimSuffix(baseName, path.Ext(baseName))
}

func NewFile(filename string, descpool *pbmeta.DescriptorPool) *File {

	self := &File{
		Sheets:   make([]*Sheet, 0),
		SheetMap: make(map[string]*Sheet),
		Name:     makeTableName(filename),
		FileName: filename,
		descpool: descpool,
	}

	var err error
	self.File, err = xlsx.OpenFile(filename)

	if err != nil {
		log.Errorln(err.Error())

		return nil
	}

	// 这里将所有sheet表都合并导出到一个pbt
	for _, sheet := range self.File.Sheets {

		header := getHeader(sheet)

		// 没有找到导出头,忽略
		if header == nil {
			continue
		}

		mySheet := newSheet(self, sheet, header)

		// TODO 添加命令行导出忽略
		self.Sheets = append(self.Sheets, mySheet)

		self.SheetMap[sheet.Name] = mySheet

	}

	return self
}

func getHeader(sheet *xlsx.Sheet) *tool.ExportHeader {

	headerString := sheet.Cell(0, 0).Value

	var header tool.ExportHeader

	if err := proto.UnmarshalText(headerString, &header); err != nil {
		return nil
	}

	return &header
}
