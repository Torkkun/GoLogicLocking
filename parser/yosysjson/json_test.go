package yosysjson

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestParsefromFile(t *testing.T) {
	var path1 = "C:\\Users\\onigi\\projects\\GoLogicLocking\\testverilog\\workspace\\fulladd\\testfulladd.json"

	file, err := os.ReadFile(path1)
	if err != nil {
		log.Fatalln(err)
	}
	jsontest := new(YosysJson)
	if err = json.Unmarshal(file, jsontest); err != nil {
		log.Fatalln(err)
	}
	fmt.Println(jsontest)
}
