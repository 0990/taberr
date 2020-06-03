package printer

type jsonPrinter struct {
}

func (self *jsonPrinter) Run(g *Global) *Stream {
	bf := NewStream()
	bf.Printf("{\n")

	bf.Printf("	\"Tool\": \"github.com/0990/taberr\",\n")
	bf.Printf("	\"Version\": \"%s\",\n", g.Version)

	printTableJson(g, bf)
	bf.Printf("\n}")
	return bf
}

func printTableJson(g *Global, stream *Stream) bool {
	stream.Printf("    \"%s\":[\n", g.PackageName)

	for rowIndex, rowData := range g.Data {
		// 每一行开始
		stream.Printf("		{ ")
		stream.Printf("\"%s\":%d", g.ErrIDLabel, rowData.ErrID)
		stream.Printf(",")
		stream.Printf("\"%s\":\"%s\"", g.ErrTypeLabel, rowData.ErrType)
		stream.Printf(",")
		stream.Printf("\"%s\":\"%s\"", g.ErrMsgLabel, rowData.ErrMsg)
		// 每一行结束
		stream.Printf(" }")
		if rowIndex < len(g.Data)-1 {
			stream.Printf(",")
		}
		stream.Printf("\n")
	}
	stream.Printf("	]")
	return true
}

func init() {
	RegisterPrinter("json", &jsonPrinter{})
}
