package vgenerator

import (
	"goll/vgenerator/funcmap"
	"os"
	"text/template"
)

func NewGenerator(val any, dryrun bool) error {
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
