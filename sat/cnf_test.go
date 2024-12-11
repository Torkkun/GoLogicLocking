package sat

import (
	"fmt"
	"testing"
)

/*
1
2
11
19
20
*/

func TestIDPool(t *testing.T) {
	vpool := NewIDPool(1, [][]int{{12, 18}, {3, 10}})
	for i := range 5 {
		fmt.Println(vpool.Id(fmt.Sprintf("v%d", i+1)))
	}
}
