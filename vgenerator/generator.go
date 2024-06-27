package vgenerator

import (
	"goll/vgenerator/funcmap"
	"os"
	"path/filepath"
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
	tmplpath, err := filepath.Abs("../template/module.tmpl")
	if err != nil {
		return err
	}

	tpl, err := template.New("module.tmpl").Funcs(tplmap).ParseFiles(tmplpath)
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
