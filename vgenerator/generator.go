package vgenerator

import (
	"fmt"
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
		return fmt.Errorf("not yet implement")
	}
	return nil
}
