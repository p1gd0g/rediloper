package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

func NewParser(fp string) *CustomParser {
	return &CustomParser{filepath: fp, typeInfo: []*info{}}
}

type CustomParser struct {
	filepath string
	msgStart bool

	// typeInfo map[string][]string
	typeInfo []*info
}

type info struct {
	name string
	list []string
}

func (p *CustomParser) Export() []*info {
	return p.typeInfo
}

// 解析注释获取table结构声明
func (p *CustomParser) ParseComments() (err error) {
	buf, err := ioutil.ReadFile(p.filepath)
	// Create the AST by parsing src.
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "", string(buf), parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// fmt.Println(f)
	ast.Walk(p, f)

	return
}

func (p *CustomParser) Visit(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.Package:
		return p
	case *ast.File:
		return p
	case *ast.GenDecl:

		if n.Tok == token.TYPE {
			defer func() {
				p.msgStart = true

				if n.Doc != nil {
					for _, comment := range n.Doc.List {
						if strings.Contains(comment.Text, "@") {

							exist := false
							for _, typeInfo := range p.typeInfo {
								if typeInfo.name == n.Specs[0].(*ast.TypeSpec).Name.Name {
									typeInfo.list = append(typeInfo.list, comment.Text)
									exist = true
									break
								}
							}
							if !exist {
								p.typeInfo = append(p.typeInfo, &info{
									name: n.Specs[0].(*ast.TypeSpec).Name.Name,
									list: []string{comment.Text},
								})
							}
						}
					}
				}
			}()
			return p
		}
	}

	return nil
}
