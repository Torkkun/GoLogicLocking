package utils

import (
	"fmt"
	"testing"
)

func TestSample(t *testing.T) {
	test := []string{"A", "B", "C"}
	for i := 0; i < 10; i++ {
		cht, err := Sample(test, 2)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(cht)
	}
}
