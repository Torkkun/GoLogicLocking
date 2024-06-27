package yosysjson

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestParseandConvertToDB(t *testing.T) {
	var path = "C:\\Users\\onigi\\projects\\GoLogicLocking\\testverilog\\workspace\\fulladd\\testfulladd.json"
	var topmodule = "fulladd"
	//var path = "C:\\Users\\onigi\\projects\\GoLogicLocking\\testverilog\\workspace\\ltika2\\ltika2.json"
	//var topmodule = "blink"
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	yosysjson := new(YosysJson)
	if err = json.Unmarshal(file, yosysjson); err != nil {
		log.Fatalln(err)
	}
	fulladdmod := yosysjson.Modules[topmodule]
	gport := PortsConvertToStoreGraph(fulladdmod.Ports)
	gname := NetsConvertToStoreGraph(fulladdmod.NetName)
	for _, v := range gport {
		fmt.Println(v)
	}
	for _, v := range gname {
		fmt.Println(v)
	}
}
