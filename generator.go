package main

import (
	"bytes"
	"html/template"
)

var (
	tplBaseDir = "templates/"

	fileTpl = "file.tpl"

	funcsTplMap = map[string][]string{
		typePersist: {},
		typeExpire:  {},
		typeHash:    {},
	}
)

type FileGenerator struct {
	sheetConfigs []*SheetConfig
	sheetData    *SheetTplData

	funcGen *FuncGenerator
}

type SheetTplData struct {
	SheetConfigs map[string][]SheetConfig

	Sheets    map[string]string
	LoadFuncs []*LoadFunc
}

type LoadFunc struct {
	SheetName string
	FuncName  string
	FuncBody  string
}

func NewGenerator(root string, sheetConfigs []*SheetConfig) *FileGenerator {
	if len(root) != 0 {
		tplBaseDir = root + tplBaseDir
	}
	return &FileGenerator{
		sheetConfigs: sheetConfigs,
		sheetData:    &SheetTplData{},
	}
}

func (g *FileGenerator) Dump() (buff []byte, err error) {
	err = g.fileConfigs()
	if err != nil {
		return
	}

	err = g.loadFuncs()
	if err != nil {
		return
	}

	tmpl := template.Must(template.New("file").Funcs(template.FuncMap{
		"unescaped": g.unescaped,
	}).ParseFiles(tplBaseDir + fileTpl))

	b := bytes.NewBufferString("")
	err = tmpl.Execute(b, g.sheetData)
	if err != nil {
		return
	}
	buff = b.Bytes()

	return
}

func (g *FileGenerator) fileConfigs() (err error) {
	sheetConfigs := make(map[string][]SheetConfig)

	for _, config := range g.sheetConfigs {

		if sheetConfigs[config.VarType] == nil {
			sheetConfigs[config.VarType] = make([]SheetConfig, 0, 50)
		}

		v := *config

		sheetConfigs[config.VarType] = append(sheetConfigs[config.VarType], v)
	}

	g.sheetData.SheetConfigs = sheetConfigs

	return
}

func (g *FileGenerator) loadFuncs() (err error) {
	funcs := make([]*LoadFunc, 0, len(g.sheetConfigs))

	g.sheetData.LoadFuncs = funcs

	return
}

func (g *FileGenerator) unescaped(x string) interface{} {
	return template.HTML(x)
}

type FuncGenerator struct {
	tmpl map[string]*template.Template
}

func (fg *FuncGenerator) Init(tplDir string) (err error) {
	tmpl := make(map[string]*template.Template)

	for typ, tplsName := range funcsTplMap {
		tplsFile := make([]string, 0, 3)

		for _, tn := range tplsName {
			tplsFile = append(tplsFile, tplDir+tn)
		}

		tmpl[typ] = template.Must(template.ParseFiles(tplsFile...))
	}

	fg.tmpl = tmpl

	return
}

func (fg *FuncGenerator) Dump(config *SheetConfig) (funcStr string, err error) {
	t, ok := fg.tmpl[config.VarType]
	if ok {
		buff := bytes.NewBufferString("")
		err = t.Execute(buff, config)
		if err != nil {
			return
		}

		funcStr = buff.String()
	}

	return
}
