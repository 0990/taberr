package printer

import (
	"text/template"
)

// TODO pbmeta解析换rune的lexer [tabtoy] {{.Comment}}
const protoTemplate = `// Generated by github.com/0990/taberr {{.XlsxName}}
// Version: {{.ToolVersion}}
// DO NOT EDIT!!
{{if ge .ProtoVersion 3}}
syntax = "proto3";
{{end}}
package {{.Package}};
{{range .Enums}}
enum {{.Name}}
{	
{{range .ProtoFields}}
	{{.Name}} = {{.Number}}; {{.Comment}}
{{end}}
}
{{end}}
{{range .Messages}}
// Defined in table: {{.DefinedTable}}
message {{.Name}}
{	
{{range .ProtoFields}}	
	{{.Alias}}
	{{.Label}}{{.TypeString}} {{.Name}} = {{.Number}}; {{.Comment}}
{{end}}
}
{{end}}
`

type protoFieldDescriptor struct {
	Name       string //枚举名
	CommentRaw string //注释原始值

	Number int
}

func (self protoFieldDescriptor) Comment() string {
	if self.CommentRaw == "" {
		return ""
	}
	return "// " + self.CommentRaw
}

type protoDescriptor struct {
	Name        string
	ProtoFields []protoFieldDescriptor
}

type protoFileModel struct {
	Package      string
	ProtoVersion int
	ToolVersion  string
	XlsxName     string
	Messages     []protoDescriptor
	Enums        []protoDescriptor
}

type protoPrinter struct {
}

func (self *protoPrinter) Run(g *Global) *Stream {

	tpl, err := template.New("proto").Parse(protoTemplate)
	if err != nil {
		log.Errorln(err)
		return nil
	}

	var m protoFileModel

	m.Package = g.PackageName
	m.ProtoVersion = g.ProtoVersion
	m.ToolVersion = g.Version
	m.XlsxName = g.FileName

	var protoD protoDescriptor
	protoD.Name = g.EnumName
	// 遍历所有
	for _, data := range g.Data {
		protoD.ProtoFields = append(protoD.ProtoFields, protoFieldDescriptor{
			Name:       data.ErrType,
			CommentRaw: data.ErrMsg,
			Number:     int(data.ErrID),
		})
	}

	m.Enums = append(m.Enums, protoD)

	bf := NewStream()

	err = tpl.Execute(bf.Buffer(), &m)
	if err != nil {
		log.Errorln(err)
		return nil
	}

	return bf
}

func init() {

	RegisterPrinter("proto", &protoPrinter{})

}
