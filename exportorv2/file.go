package exportorv2

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/tealeg/xlsx"
)

type valueRepeatData struct {
	fd    *model.FieldDefine
	value string
}

type File struct {
	TypeSet  *model.BuildInTypeSet // 1个类型描述表
	Name     string                // 电子表格的名称(文件名无后缀)
	FileName string

	valueRepByKey map[valueRepeatData]bool
}

func (self *File) Export(filename string) *model.Table {

	self.Name = makeTableName(filename)
	self.FileName = filename

	log.Infof("%s\n", filepath.Base(filename))

	file, err := xlsx.OpenFile(filename)

	if err != nil {
		log.Errorln(err.Error())

		return nil
	}

	// 解析类型表
	for _, rawSheet := range file.Sheets {

		if isTypeSheet(rawSheet.Name) {
			if self.TypeSet != nil {
				log.Errorln("@Types sheet should keep one!")
				return nil
			}

			typeSheet := newTypeSheet(NewSheet(self, rawSheet))

			// 从cell添加类型
			if !typeSheet.Parse() {
				return nil
			}

			self.TypeSet = typeSheet.BuildInTypeSet

		}
	}

	// 没有这个类型, 默认填一个
	if self.TypeSet == nil {
		log.Errorln("'@Types' sheet not found in ", filename)
		return nil
	}

	var tab model.Table

	var needAddRowType bool = true

	for _, rawSheet := range file.Sheets {

		if !isTypeSheet(rawSheet.Name) {

			dSheet := newDataSheet(NewSheet(self, rawSheet))

			if !dSheet.Valid() {
				continue
			}

			log.Infof("            %s", rawSheet.Name)

			// TODO 只使用第一个sheet中的protoheader定义
			// TODO 其他Sheet可以在顶部定义一个标记@RefProtoHeader, 引用前面的protoheader
			if !dSheet.Export(self, &tab) {
				return nil
			}

			if needAddRowType {
				self.makeRowBuildInType(self.TypeSet, dSheet.headerFields)
				needAddRowType = false
			}

		}
	}

	return &tab
}

func (self *File) makeRowBuildInType(ts *model.BuildInTypeSet, rootField []*model.FieldDefine) {

	rowType := model.NewBuildInType()
	rowType.Name = fmt.Sprintf("%sDefine", ts.Pragma.TableName)
	rowType.Kind = model.BuildInTypeKind_Struct
	self.TypeSet.Add(rowType)

	// 将表格中的列添加到类型中, 方便导出
	for _, field := range rootField {

		rowType.Add(field)
	}

	fileType := model.NewBuildInType()
	fileType.RootFile = true

	// 强制命名文件类型名
	if ts.Pragma.FileTypeName != "" {
		fileType.Name = ts.Pragma.FileTypeName
	} else {
		// 文件类型名: Table名+File
		fileType.Name = fmt.Sprintf("%sFile", ts.Pragma.TableName)
	}

	fileType.Kind = model.BuildInTypeKind_Struct

	var rowTypeField model.FieldDefine
	rowTypeField.Name = ts.Pragma.TableName
	rowTypeField.Type = model.FieldType_Struct
	rowTypeField.IsRepeated = true
	rowTypeField.BuildInType = rowType
	rowTypeField.Comment = ts.Pragma.TableName
	fileType.Add(&rowTypeField)

	self.TypeSet.FileType = fileType

	self.TypeSet.Add(fileType)
}

func (self *File) checkValueRepeat(fd *model.FieldDefine, value string) bool {

	key := valueRepeatData{
		fd:    fd,
		value: value,
	}

	if _, ok := self.valueRepByKey[key]; ok {
		return false
	}

	self.valueRepByKey[key] = true

	return true
}

func isTypeSheet(name string) bool {
	return strings.TrimSpace(name) == "@Types"
}

// 制作表名
func makeTableName(filename string) string {
	baseName := path.Base(filename)
	return strings.TrimSuffix(baseName, path.Ext(baseName))
}

func NewFile() *File {

	self := &File{
		valueRepByKey: make(map[valueRepeatData]bool),
	}

	return self
}
