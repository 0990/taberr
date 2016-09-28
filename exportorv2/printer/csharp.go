package printer

import (
	"fmt"
	"text/template"

	"github.com/davyxu/tabtoy/exportorv2/model"
)

// TODO pbmeta解析换rune的lexer [tabtoy] {{.Comment}}
const csharpTemplate = `// Generated by github.com/davyxu/tabtoy
// Version: {{.ToolVersion}}
// DO NOT EDIT!!
using System.Collections.Generic;
using System.IO;

namespace {{.Namespace}}{{$globalIndex:=.Indexes}}
{
	{{range .Enums}}
	public enum {{.Name}}
	{
	{{range .Fields}}	
		// {{.Comment}}
		{{.FieldDescriptor.Name}} = {{.FieldDescriptor.EnumValue}},
	{{end}}
	}
	{{end}}
	{{range .Classes}}
	public partial class {{.Name}} : tabtoy.DataObject
	{	
	{{range .Fields}}	
		// {{.Comment}}
		{{.TypeCode}}
	{{end}}
	
	{{if .IsCombine}}{{range $globalIndex}}
	 	Dictionary<{{.IndexType}}, {{.RowType}}> _{{.RowName}}By{{.IndexName}} = new Dictionary<{{.IndexType}}, {{.RowType}}>();
        public {{.RowType}} Get{{.RowName}}By{{.IndexName}}({{.IndexType}} {{.IndexName}})
        {
            {{.RowType}} ret;
            if ( _{{.RowName}}By{{.IndexName}}.TryGetValue( {{.IndexName}}, out ret ) )
            {
                return ret;
            }

            return null;
        }
	{{end}}{{end}}
		public void Deserialize( tabtoy.DataReader reader )
		{
			{{range .Fields}}
			// {{.Comment}}
			if ( reader.MatchTag({{.Tag}}) )
			{
				{{.ReadCode}}
			}
			{{end}}
			{{if .IsCombine}}{{range $a, $row :=.IndexedFields}}
			// Build {{$row.FieldDescriptor.Name}} Index
            for( int i = 0;i< this.{{$row.FieldDescriptor.Name}}.Count;i++)
            {
                var element = this.{{$row.FieldDescriptor.Name}}[i];
				{{range $b, $key := .IndexKeys}}
                _{{$row.FieldDescriptor.Name}}By{{$key.Name}}.Add(element.{{$key.Name}}, element);                
				{{end}}
            }
			{{end}}{{end}}
		}
	}
	{{end}}

}
`

type indexField struct {
	TableIndex
}

func (self indexField) IndexName() string {
	return self.Index.Name
}

func (self indexField) RowType() string {
	return self.Row.Complex.Name
}

func (self indexField) RowName() string {
	return self.Row.Name
}

func (self indexField) IndexType() string {

	switch self.Index.Type {
	case model.FieldType_Int32:
		return "int"
	case model.FieldType_UInt32:
		return "uint"
	case model.FieldType_Int64:
		return "long"
	case model.FieldType_UInt64:
		return "ulong"
	case model.FieldType_String:
		return "string"
	case model.FieldType_Float:
		return "float"
	case model.FieldType_Bool:
		return "bool"
	case model.FieldType_Enum:

		return self.Index.Complex.Name
	default:
		log.Errorf("%s can not be index ", self.Index.String())
	}

	return "unknown"
}

type csharpField struct {
	*model.FieldDescriptor

	IndexKeys []*model.FieldDescriptor
}

func (self csharpField) Comment() string {

	if self.FieldDescriptor.Comment != "" {
		return self.FieldDescriptor.Comment
	}

	return self.FieldDescriptor.Meta.Alias
}

func (self csharpField) ReadCode() string {

	var baseType string

	switch self.Type {
	case model.FieldType_Int32:
		baseType = "Int32"
	case model.FieldType_UInt32:
		baseType = "UInt32"
	case model.FieldType_Int64:
		baseType = "Int64"
	case model.FieldType_UInt64:
		baseType = "UInt64"
	case model.FieldType_String:
		baseType = "String"
	case model.FieldType_Float:
		baseType = "Float"
	case model.FieldType_Bool:
		baseType = "Bool"
	case model.FieldType_Enum:

		if self.Complex == nil {
			return "unknown"
		}

		baseType = fmt.Sprintf("Enum<%s>", self.Complex.Name)

	case model.FieldType_Struct:
		if self.Complex == nil {
			return "unknown"
		}

		baseType = fmt.Sprintf("Struct<%s>", self.Complex.Name)

	}

	if self.IsRepeated {
		return fmt.Sprintf("reader.ReadList_%s( this.%s );", baseType, self.Name)
	} else {
		return fmt.Sprintf("this.%s = reader.Read%s( );", self.Name, baseType)
	}

}

func (self csharpField) Tag() string {

	return fmt.Sprintf("0x%x", self.FieldDescriptor.Tag())
}

func (self csharpField) TypeCode() string {

	var raw string

	switch self.Type {
	case model.FieldType_Int32:
		raw = "int"
	case model.FieldType_UInt32:
		raw = "uint"
	case model.FieldType_Int64:
		raw = "long"
	case model.FieldType_UInt64:
		raw = "ulong"
	case model.FieldType_String:
		raw = "string"
	case model.FieldType_Float:
		raw = "float"
	case model.FieldType_Bool:
		raw = "bool"
	case model.FieldType_Enum:
		if self.Complex == nil {
			log.Errorln("unknown enum type ", self.Type)
			return "unknown"
		}

		raw = self.Complex.Name
	case model.FieldType_Struct:
		if self.Complex == nil {
			log.Errorln("unknown struct type ", self.Type)
			return "unknown"
		}

		raw = self.Complex.Name

		// 非repeated的结构体
		if !self.IsRepeated {
			return fmt.Sprintf("public %s %s = new %s();", raw, self.Name, raw)
		}

	default:
		raw = "unknown"
	}

	if self.IsRepeated {
		return fmt.Sprintf("public List<%s> %s = new List<%s>();", raw, self.Name, raw)
	}

	return fmt.Sprintf("public %s %s = %s;", raw, self.Name, wrapCSharpDefaultValue(self.FieldDescriptor))
}

func wrapCSharpDefaultValue(fd *model.FieldDescriptor) string {
	switch fd.Type {
	case model.FieldType_Enum:
		return fmt.Sprintf("%s.%s", fd.Complex.Name, fd.DefaultValue())
	case model.FieldType_String:
		return fmt.Sprintf("\"%s\"", fd.DefaultValue())
	}

	return fd.DefaultValue()
}

type structModel struct {
	*model.Descriptor
	Fields        []csharpField
	IndexedFields []csharpField // 与csharpField.IndexKeys组成树状的索引层次
}

func (self *structModel) Name() string {
	return self.Descriptor.Name
}

func (self *structModel) IsCombine() bool {
	return self.Descriptor.Usage == model.DescriptorUsage_CombineStruct
}

type csharpFileModel struct {
	Namespace   string
	ToolVersion string
	Classes     []*structModel
	Enums       []*structModel
	Indexes     []indexField // 全局的索引
}

type csharpPrinter struct {
}

func (self *csharpPrinter) Run(g *Globals) *BinaryFile {

	tpl, err := template.New("csharp").Parse(csharpTemplate)
	if err != nil {
		log.Errorln(err)
		return nil
	}

	var m csharpFileModel

	m.Namespace = g.FileDescriptor.Pragma.Package
	m.ToolVersion = g.Version

	// combinestruct的全局索引
	for _, ti := range g.GlobalIndexes {

		// 索引也限制
		if !ti.Index.Parent.File.MatchTag(".cs") {
			continue
		}

		m.Indexes = append(m.Indexes, indexField{TableIndex: ti})
	}

	// 遍历所有类型
	for _, d := range g.FileDescriptor.Descriptors {

		// 这给被限制输出
		if !d.File.MatchTag(".cs") {
			continue
		}

		var sm structModel
		sm.Descriptor = d

		switch d.Kind {
		case model.DescriptorKind_Struct:
			m.Classes = append(m.Classes, &sm)
		case model.DescriptorKind_Enum:
			m.Enums = append(m.Enums, &sm)
		}

		// 遍历字段
		for _, fd := range d.Fields {

			// 对CombineStruct的XXDefine对应的字段
			if d.Usage == model.DescriptorUsage_CombineStruct {

				// 这个字段被限制输出
				if fd.Complex != nil && !fd.Complex.File.MatchTag(".cs") {
					continue
				}

				// 这个结构有索引才创建
				if fd.Complex != nil && len(fd.Complex.Indexes) > 0 {

					// 被索引的结构
					indexedField := csharpField{FieldDescriptor: fd}

					// 索引字段
					for _, key := range fd.Complex.Indexes {
						indexedField.IndexKeys = append(indexedField.IndexKeys, key)
					}

					sm.IndexedFields = append(sm.IndexedFields, indexedField)
				}

			}

			csField := csharpField{FieldDescriptor: fd}

			sm.Fields = append(sm.Fields, csField)

		}

	}

	bf := NewBinaryFile()

	err = tpl.Execute(bf.Buffer(), &m)
	if err != nil {
		log.Errorln(err)
		return nil
	}

	return bf
}

func init() {

	RegisterPrinter(".cs", &csharpPrinter{})

}
