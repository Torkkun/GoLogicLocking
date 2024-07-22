package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"testing"
)

func TestEnvironment(t *testing.T) {
	if runtime.GOOS == "windows" {
		out, err := exec.Command("echo", "Hello").Output()
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(out)
	}
}
