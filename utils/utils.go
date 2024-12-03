package utils

import (
	"errors"
	"math/rand"
)

// ErrEmptySlice slice is empty.
var ErrEmptySlice = errors.New("slice is empty")

// ErrSizeTooLarge size is greater than slice.
var ErrSizeTooLarge = errors.New("size is greater than slice")

// ErrNegativeSize size is negative.
var ErrNegativeSize = errors.New("size is negative")

func Choice[T any](in []T) (T, error) {
	if len(in) == 0 {
		var v T // zero value
		return v, ErrEmptySlice
	}

	return in[rand.Intn(len(in))], nil
}

func Choices[T any](in []T, size int) ([]T, error) {
	if len(in) == 0 {
		return nil, ErrEmptySlice
	}

	if size < 0 {
		return nil, ErrNegativeSize
	}

	out := make([]T, 0, size)
	for i := 0; i < size; i++ {
		v, _ := Choice(in)
		out = append(out, v)
	}

	return out, nil
}

func Sample[T any](in []T, size int) ([]T, error) {
	if len(in) == 0 {
		return nil, ErrEmptySlice
	}

	if size < 0 {
		return nil, ErrNegativeSize
	}

	if len(in) < size {
		return nil, ErrSizeTooLarge
	}

	idx := make([]int, len(in))
	for i := 0; i < len(in); i++ {
		idx[i] = i
	}

	rand.Shuffle(len(idx), func(i, j int) {
		idx[i], idx[j] = idx[j], idx[i]
	})

	out := make([]T, 0, size)
	for i := 0; i < size; i++ {
		out = append(out, in[idx[i]])
	}

	return out, nil
}

func RandomNumbers(num, count int) []int {
	numbers := rand.Perm(num)
	result := numbers[:count]
	return result
}

func WriteSpace() string {
	return " "
}

func WriteIndentTab() string {
	return "  "
}

func WriteNewLine() string {
	return "\n"
}

func WriteSemiColon() string {
	return ";"
}

func ProductBool(repeat int) func() []bool {
	bools := []bool{false, true}
	p := make([]bool, repeat)
	x := make([]int, len(p))

	return func() []bool {
		p := p[:len(x)]
		for i, xi := range x {
			p[i] = bools[xi]
		}
		for i := len(x) - 1; i >= 0; i-- {
			x[i]++
			if x[i] < len(bools) {
				break
			}
			x[i] = 0
			if i <= 0 {
				x = x[0:0]
				break
			}
		}
		return p
	}

}
