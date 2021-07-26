package main

import "strings"

var (
	dataBuilders map[string]IDataBuilder
)

type SheetConfig struct {
	VarType string

	TypeName  string
	ClassName string
	KeyFormat string
	KeyInputs []string
	KeyParams []string
}

func init() {
	dataBuilders = make(map[string]IDataBuilder)

	dataBuilders["@PE"] = &PersistDataBuilder{}
	dataBuilders["@EX"] = &ExpireDataBuilder{}
	dataBuilders["@HA"] = &HashDataBuilder{}

}

type IDataBuilder interface {
	Build(msg string, params []string) (config *SheetConfig, err error)
}

type BaseDataBuilder struct {
	Builder interface{}
}

// @PE
type PersistDataBuilder struct {
	BaseDataBuilder
}

func (p *PersistDataBuilder) Build(typeName string, params []string) (config *SheetConfig,
	err error) {

	config = &SheetConfig{
		VarType:   typePersist,
		ClassName: params[1],
		TypeName:  typeName,
		KeyFormat: params[2],
	}

	for i := 3; i < len(params)-1; i += 2 {
		config.KeyParams = append(config.KeyParams, params[i])
		config.KeyInputs = append(config.KeyInputs, strings.Join(params[i:i+2], " "))
	}

	return
}

// @EX
type ExpireDataBuilder struct {
	BaseDataBuilder
}

func (p *ExpireDataBuilder) Build(typeName string, params []string) (config *SheetConfig,
	err error) {

	config = &SheetConfig{
		VarType:   typeExpire,
		ClassName: params[1],
		TypeName:  typeName,
		KeyFormat: params[2],
	}

	for i := 3; i < len(params)-1; i += 2 {
		config.KeyParams = append(config.KeyParams, params[i])
		config.KeyInputs = append(config.KeyInputs, strings.Join(params[i:i+2], " "))
	}

	return
}

// @HA
type HashDataBuilder struct {
	BaseDataBuilder
}

func (p *HashDataBuilder) Build(typeName string, params []string) (config *SheetConfig,
	err error) {

	config = &SheetConfig{
		VarType:   typeHash,
		ClassName: params[1],
		TypeName:  typeName,
		KeyFormat: params[2],
	}

	for i := 3; i < len(params)-1; i += 2 {
		config.KeyParams = append(config.KeyParams, params[i])
		config.KeyInputs = append(config.KeyInputs, strings.Join(params[i:i+2], " "))
	}

	return
}
