package vgenerator

import (
	"goll/vgenerator/funcmap"
	"os"
	"text/template"
)

type GenerateData struct {
	ModuleName string
	PortList   []string
	Parameter  []string // not implement
	PortDecl   []*funcmap.PortDecl
	Register   []string              // not implement
	EventDecl  []string              // not implement
	NetDecl    []*funcmap.NetDecl    // not implement
	AssignDecl []*funcmap.AssignDecl // not implement
}

func NewGenerator(val GenerateData, dryrun bool) error {
	tplmap := funcmap.NewFuncMap()
	tpl, err := template.New("module.tmpl").Funcs(tplmap).ParseFiles("template/module.tmpl")
	if err != nil {
		return err
	}
	// dryrun trueでStdOut
	if dryrun {
		if err := tpl.Execute(os.Stdout, val); err != nil {
			return err
		}
	} else {
		file, err := os.Create("test.v")
		if err != nil {
			return err
		}
		if err := tpl.Execute(file, val); err != nil {
			return err
		}

	}
	return nil
}
