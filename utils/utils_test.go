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

func TestRandomNumbers(t *testing.T) {
	test := RandomNumbers(5, 3)
	fmt.Println(test)
}

func TestProduct(t *testing.T) {
	product := ProductBool(3)
	for {
		product := product()
		if len(product) == 0 {
			break
		}
		fmt.Println(product)
	}
}
