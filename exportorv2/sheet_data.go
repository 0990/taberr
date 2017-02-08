package exportorv2

import (
	"strings"

	"github.com/davyxu/tabtoy/exportorv2/filter"
	"github.com/davyxu/tabtoy/exportorv2/i18n"
	"github.com/davyxu/tabtoy/exportorv2/model"
	"github.com/davyxu/tabtoy/util"
)

/*
	Sheet数据表单的处理

*/

type DataSheet struct {
	*Sheet
}

func (self *DataSheet) Valid() bool {

	name := strings.TrimSpace(self.Sheet.Name)
	if name != "" && name[0] == '#' {
		return false
	}

	return self.GetCellData(0, 0) != ""
}

func dataProcessor(file *File, fd *model.FieldDescriptor, raw string, node *model.Node) bool {

	// 列表
	if fd.IsRepeated {

		spliter := fd.ListSpliter()

		// 使用多格子实现的repeated
		if spliter == "" {

			if _, ok := filter.ConvertValue(fd, raw, file.GlobalFD, node); !ok {
				goto ConvertError
			}

		} else {
			// 一个格子切割的repeated

			valueList := strings.Split(raw, spliter)

			for _, v := range valueList {

				if _, ok := filter.ConvertValue(fd, v, file.GlobalFD, node); !ok {
					goto ConvertError
				}
			}

		}

	} else {

		// 单值
		if cv, ok := filter.ConvertValue(fd, raw, file.GlobalFD, node); !ok {
			goto ConvertError

		} else {

			// 值重复检查
			if fd.Meta.GetBool("RepeatCheck") && !file.checkValueRepeat(fd, cv) {
				log.Errorf("%s, %s raw: '%s'", i18n.String(i18n.DataSheet_ValueRepeated), fd.String(), cv)
				return false
			}
		}

	}

	return true

ConvertError:

	log.Errorf("%s, %s raw: '%s'", i18n.String(i18n.DataSheet_ValueConvertError), fd.String(), raw)

	return false
}

func (self *DataSheet) Export(file *File, tab *model.Table, dataHeader *DataHeader) bool {

	// 是否继续读行
	var readingLine bool = true

	var meetEmptyLine bool

	var warningAfterEmptyLineDataOnce bool

	// 遍历每一行
	for self.Row = DataSheetHeader_DataBegin; readingLine; self.Row++ {

		// 整行都是空的
		if self.IsFullRowEmpty(self.Row, dataHeader.RawFieldCount()) {

			// 再次碰空行, 表示确实是空的
			if meetEmptyLine {
				break

			} else {
				meetEmptyLine = true
			}

			continue

		} else {

			//已经碰过空行, 这里又碰到数据, 说明有人为隔出的空行, 做warning提醒, 防止数据没导出
			if meetEmptyLine && !warningAfterEmptyLineDataOnce {
				r, _ := self.GetRC()

				log.Warnf("%s %s|%s(%s)", i18n.String(i18n.DataSheet_RowDataSplitedByEmptyLine), self.file.FileName, self.Name, util.ConvR1C1toA1(r, 1))

				warningAfterEmptyLineDataOnce = true
			}

		}

		record := model.NewRecord()

		// 遍历每一列
		for self.Column = 0; self.Column < dataHeader.RawFieldCount(); self.Column++ {

			fieldDef := dataHeader.RawField(self.Column)

			// 数据大于列头时, 结束这个列
			if fieldDef == nil {
				break
			}

			// #开头表示注释, 跳过
			if strings.Index(fieldDef.Name, "#") == 0 {
				continue
			}

			rawValue := self.GetCellData(self.Row, self.Column)

			// repeated的, 没有填充的, 直接跳过, 不生成数据
			if rawValue == "" && fieldDef.Meta.GetString("Default") == "" {
				continue
			}

			node := record.NewNodeByDefine(fieldDef)

			// 结构体要多添加一个节点, 处理repeated 结构体情况
			if fieldDef.Type == model.FieldType_Struct {

				node.StructRoot = true
				node = node.AddKey(fieldDef)

			}

			//log.Debugf("raw: %v  r:%d c: %d", rawValue, self.Row, self.Column)

			if !dataProcessor(file, fieldDef, rawValue, node) {
				goto ErrorStop
			}

		}

		tab.Add(record)

	}

	return true

ErrorStop:

	r, c := self.GetRC()

	log.Errorf("%s|%s(%s)", self.file.FileName, self.Name, util.ConvR1C1toA1(r, c))
	return false
}

func newDataSheet(sheet *Sheet) *DataSheet {

	return &DataSheet{
		Sheet: sheet,
	}
}
