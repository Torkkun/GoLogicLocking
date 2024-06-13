package funcmap

import (
	"fmt"
)

func portListElement(io []string) string {
	// ( 引数への変換 );
	var arg string
	for i, v := range io {
		if len(io)-1 == i {
			arg += v
		} else {
			arg += fmt.Sprintf("%s, ", v)
		}
	}
	return arg
}
