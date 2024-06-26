package funcmap

import (
	"bytes"
	"fmt"
	"goll/parser/yosysjson"
	"goll/utils"
	"text/template"
)

func assignDeclElement(assigndecls []*AssignDecl) (string, error) {
	var returnstr string
	var err error
	var data interface{}
	for _, assigndecl := range assigndecls {
		var templatename string
		var globpath string

		switch assigndecl.ExpressionType {
		case yosysjson.BUF:
			templatename = "buf.tmpl"
			globpath = "buf"
			data = assigndecl.Connection.(BuforNot)
		case yosysjson.NOT:
			templatename = "not.tmpl"
			globpath = "logic"
			data = assigndecl.Connection.(BuforNot)
		case yosysjson.AND:
			templatename = "and.tmpl"
			globpath = "logic"
			data = assigndecl.Connection.(Logic)
		case yosysjson.NAND:
			templatename = "nand.tmpl"
			globpath = "logic"
			data = assigndecl.Connection.(Logic)
		//case yosysjson.ANDNOT:
		//	t, err = t.ParseFiles("vgenerator/template/cells/logic/andnot.tmpl")
		case yosysjson.OR:
			templatename = "or.tmpl"
			globpath = "logic"
			data = assigndecl.Connection.(Logic)
		case yosysjson.NOR:
			templatename = "nor.tmpl"
			globpath = "logic"
			data = assigndecl.Connection.(Logic)
		//case yosysjson.ORNOT:
		//	t, err = t.ParseFiles("vgenerator/template/cells/logic/ornot.tmpl")
		case yosysjson.XOR:
			templatename = "xor.tmpl"
			globpath = "logic"
			data = assigndecl.Connection.(Logic)
		case yosysjson.XNOR:
			templatename = "xnor.tmpl"
			globpath = "logic"
			data = assigndecl.Connection.(Logic)
		default:
			return "", fmt.Errorf("not implement this type %s", assigndecl.ExpressionType)
		}
		t := template.New(templatename)
		t, err = t.ParseGlob(fmt.Sprintf("template/cells/%s/*.tmpl", globpath))
		if err != nil {
			return "", err
		}
		var tpl bytes.Buffer
		if err = t.Execute(&tpl, data); err != nil {
			return "", err
		}
		returnstr += utils.WriteIndentTab() + "assign" + utils.WriteSpace() + tpl.String() + utils.WriteNewLine()
	}
	return returnstr, nil
}
