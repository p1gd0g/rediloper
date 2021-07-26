package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	typePersist = "PE"
	typeExpire  = "EX"
	typeHash    = "HA"
)

const (
	Separator = "|"
)

func main() {

	var err error

	log.SetFlags(log.Lshortfile | log.Ltime)

	filePath := os.Args[1]
	outputPath := os.Args[2]
	tplPath := os.Args[3]

	ps := NewParser(filePath)
	err = ps.ParseComments()
	if err != nil {
		panic(err.Error())
	}
	comments := ps.Export()
	if len(comments) == 0 {
		panic("have no comments")
	}

	codeBytes, err := BuildLoader(tplPath, comments)
	if err != nil {
		panic(err.Error())
	}

	codeBytes, err = format.Source(codeBytes)
	if err != nil {
		panic(err.Error())
	}

	fileName := outputPath + "auto_" + filepath.Base(filePath)
	err = ioutil.WriteFile(fileName, codeBytes, 0644)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Generate successfully!")
}

func BuildLoader(tplPath string, typeInfo []*info) (buff []byte, err error) {
	sheetConfigs := make([]*SheetConfig, 0, 1)
	for _, typeComments := range typeInfo {
		for _, comment := range typeComments.list {
			params := strings.Split(comment, Separator)

			params[0] = strings.Trim(params[0], "// ")

			declares := strings.Split(params[0], ":")
			if builder, ok := dataBuilders[declares[0]]; ok {
				config, err := builder.Build(typeComments.name, params)
				if err != nil {
					continue
				}

				sheetConfigs = append(sheetConfigs, config)
			}
		}
	}

	buff, err = NewGenerator(tplPath, sheetConfigs).Dump()
	if err != nil {
		return
	}

	return
}
