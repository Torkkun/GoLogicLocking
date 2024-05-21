package parser

import (
	"bufio"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
)

type GateType int

const (
	_ GateType = iota
	Xor
	And
	Or
)

type prevType int

const (
	_ prevType = iota
	Lvalue
	Rvalue
	Gate
)

type Node struct {
	In1  string
	In2  string
	Gate GateType
	Out  string
}

type Decl struct {
	IOPorts *[]IOPort
	Wires   *[]Wire
}

type ParseResult struct {
	Nodes       map[int]Node
	Declaration *Decl
}

func NewParse() ParseResult {
	file := NewReader("test.txt")
	scanner := NewScanner(file)
	// Nodes
	newNodes := map[int]Node{}
	// Declsetup
	newDecl := new(Decl)
	newDecl.IOPorts = new([]IOPort)
	newDecl.Wires = new([]Wire)

	ioportlist := []string{} //このリストをもとにWireのIOを除外するかどうか決める

	var pt prevType

	for scanner.Scan() {
		txt := scanner.Text()
		strline := strings.Fields(txt)

		switch strline[0] {
		case "Portlist:":
			// 次の値を取得
			strline = peekScan(scanner)
			isPort := strline[0]
			for isPort == "Port:" {
				portname := strline[1]
				portname = removechar(portname)
				ioportlist = append(ioportlist, portname)
				// 次の値を取得して更新
				strline = peekScan(scanner)
				isPort = strline[0]
			}
		case "Decl:":
			strline = peekScan(scanner)
			newDecl.parseIOPortlist(strline, ioportlist)

		case "Assign:":
			nextstr := deparseStr(strline[1:])
			from, to, err := parseFromTo(nextstr)
			if err != nil {
				log.Fatalln(err)
			}
			if from != to {
				log.Fatalln("Assign width is not one")
			}

		case "Lvalue:":
			nextstr := deparseStr(strline[1:])
			at, err := parseAt(nextstr)
			if err != nil {
				log.Fatalln(err)
			}
			newNodes[at] = Node{}
			pt = Lvalue

		case "Rvalue:":
			// peek scan
			strline := peekScan(scanner)
			// ゲートが来るのを期待
			nextstr := deparseStr(strline[1:])
			at, err := parseAt(nextstr)
			if err != nil {
				log.Fatalln(err)
			}
			mapval := newNodes[at]
			mapval.Gate, err = parseGate(strline[0])
			if err != nil {
				log.Fatalln(err)
			}
			newNodes[at] = mapval
			pt = Gate

		case "Identifier:":
			switch pt {
			case Lvalue:
				nextstr := deparseStr(strline[2:])
				at, err := parseAt(nextstr)
				if err != nil {
					log.Fatalln(err)
				}
				mapval := newNodes[at]
				mapval.Out = strline[1]
				newNodes[at] = mapval

			case Gate:
				nextstr := deparseStr(strline[2:])
				at, err := parseAt(nextstr)
				if err != nil {
					log.Fatalln(err)
				}
				mapval := newNodes[at]
				mapval.In1 = strline[1]
				strline := peekScan(scanner)
				mapval.In2 = strline[1]
				newNodes[at] = mapval
			}

		default:
			// ひとまず考えない部分はスキップ
			continue
		}
	}
	return ParseResult{
		Nodes:       newNodes,
		Declaration: newDecl,
	}

}

type IOPortType string

const (
	IN  IOPortType = "IN"
	OUT IOPortType = "OUT"
)

type IOPort struct {
	Type IOPortType
	Name string
}

type Wire struct {
	Name string
}

func (d Decl) parseIOPortlist(declstr []string, list []string) error {
	declType := declstr[0]
	switch declType {
	case "Wire:":
		name := declstr[1]
		name = removechar(name)
		if !slices.Contains(list, name) {
			*d.Wires = append(*d.Wires, Wire{Name: name})
		}
		// なかったら何もせずスキップ

	case "Input:":
		name := declstr[1]
		name = removechar(name)
		*d.IOPorts = append(*d.IOPorts, IOPort{Type: IN, Name: name})

	case "Output:":
		name := declstr[1]
		name = removechar(name)
		*d.IOPorts = append(*d.IOPorts, IOPort{Type: OUT, Name: name})
	default:
		return fmt.Errorf("decl is not recognize %s", declType)
	}
	return nil
}

// (from N to N)をパースしてfromとtoを取得
func parseFromTo(str string) (int, int, error) {
	str = strings.Replace(str, "(", "", -1)
	str = strings.Replace(str, ")", "", -1)
	splitstr := strings.Fields(str)

	from, err := strconv.Atoi(splitstr[1])
	if err != nil {
		log.Fatalln(splitstr)
		return -1, -1, err
	}
	to, err := strconv.Atoi(splitstr[3])
	if err != nil {
		return -1, -1, err
	}

	return from, to, nil
}

// (at N)をパースしてNを取得
func parseAt(str string) (int, error) {
	str = strings.Replace(str, "(", "", -1)
	str = strings.Replace(str, ")", "", -1)
	splitstr := strings.Fields(str)

	at, err := strconv.Atoi(splitstr[1])
	if err != nil {
		return -1, err
	}
	return at, nil
}

func parseGate(str string) (GateType, error) {
	splitstr := strings.Fields(str)

	gate := splitstr[0]
	switch gate {
	case "Xor:":
		return Xor, nil
	case "Or:":
		return Or, nil
	case "And:":
		return And, nil
	default:
		return 0, fmt.Errorf("GateType is not implementence: %s", gate)
	}
}

func deparseStr(strarry []string) string {
	tmpstr := strings.Join(strarry, " ")
	return tmpstr
}

func peekScan(scanner *bufio.Scanner) []string {
	scanner.Scan()
	txt := scanner.Text()
	strline := strings.Fields(txt)
	return strline
}

// パース時に出てくる不必要な文字を除去する
func removechar(str string) string {
	// とりあえず　,　だけ
	return strings.Replace(str, ",", "", -1)
}
