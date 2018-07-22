package printer

type Data struct {
	ErrID   int32
	ErrType string
	ErrMsg  string
}

type Global struct {
	PackageName      string
	EnumName         string
	ProtoVersion     int
	Version          string
	ErrIDLabel       string
	ErrTypeLabel     string
	ErrMsgLabel      string
	Data             []Data
	ErrIDReapeatMap  map[int32]bool
	ErrTypeRepeatMap map[string]bool

	Printers []*PrinterContext
}

func NewGlobal() *Global {
	self := &Global{
		ErrIDLabel:       "ErrID",
		ErrTypeLabel:     "ErrType",
		ErrMsgLabel:      "ErrMsg",
		PackageName:      "emsg",
		ProtoVersion:     3,
		Version:          "0.0.1",
		EnumName:         "Err",
		ErrIDReapeatMap:  make(map[int32]bool),
		ErrTypeRepeatMap: make(map[string]bool),
	}

	return self
}

func (g *Global) CheckErrIDRepeate(errID int32) bool {
	if _, ok := g.ErrIDReapeatMap[errID]; ok {
		return true
	} else {
		g.ErrIDReapeatMap[errID] = true
	}
	return false
}

func (g *Global) CheckErrTypeRepeate(errType string) bool {
	if _, ok := g.ErrTypeRepeatMap[errType]; ok {
		return true
	} else {
		g.ErrTypeRepeatMap[errType] = true
	}
	return false
}

func (self *Global) AddOutputType(name string, outfile string) {

	if p, ok := printerByExt[name]; ok {
		self.Printers = append(self.Printers, &PrinterContext{
			p:       p,
			outFile: outfile,
			name:    name,
		})
	} else {
		panic("output type not found:" + name)
	}

}

func (self *Global) Print() bool {

	log.Infof("==========begin print==========")

	for _, p := range self.Printers {

		if !p.Start(self) {
			return false
		}
	}

	return true

}
