package sat

import (
	"math"
	"slices"
	"sort"
)

type IDPool struct {
	top      int
	occupied [][]int
	obj2id   map[interface{}]int
	id2obj   map[int]interface{}
}

func NewIDPool(startFrom int, occupied [][]int) *IDPool {
	pool := new(IDPool)
	pool.Restart(startFrom, occupied)
	return pool
}

func (p *IDPool) Restart(startFrom int, occupied [][]int) {
	p.top = startFrom - 1
	p.occupied = append([][]int{}, occupied...)
	sort.Slice(p.occupied, func(i, j int) bool {
		return p.occupied[i][0] < p.occupied[j][0]
	})
	p.obj2id = make(map[interface{}]int)
	p.id2obj = make(map[int]interface{})
}

func (p *IDPool) Id(obj interface{}) int {
	if obj != nil {
		if id, exists := p.obj2id[obj]; exists {
			return id
		}
		id := p.Next()
		p.obj2id[obj] = id
		p.id2obj[id] = obj
		return id
	}
	return p.Next()
}

func (p *IDPool) Obj(vid int) interface{} {
	if obj, exists := p.id2obj[vid]; exists {
		return obj
	}
	return nil
}

func (p *IDPool) Occupy(start, stop int) {
	if stop >= start {
		if len(p.occupied) > 0 {
			last := p.occupied[len(p.occupied)-1]
			if last[0] >= start && last[1] <= stop {
				p.occupied = p.occupied[:len(p.occupied)-1]
			}
		}
		p.occupied = append(p.occupied, []int{start, stop})
		sort.Slice(p.occupied, func(i, j int) bool {
			return p.occupied[i][0] < p.occupied[j][0]
		})
	}
}

func (p *IDPool) Next() int {
	p.top++
	for len(p.occupied) > 0 && p.top >= p.occupied[0][0] {
		if p.top <= p.occupied[0][1] {
			p.top = p.occupied[0][1] + 1
		}
		p.occupied = p.occupied[1:]
	}
	return p.top
}

type Formula struct {
	Nv          int
	FromClouses [][]int
}

// int(id)のみ現在は受け付ける
func NewCNF() *Formula {
	return &Formula{
		Nv:          0,
		FromClouses: make([][]int, 0),
	}
}

func (form *Formula) Append(clause []int) {
	absclause := []int{}
	for _, c := range clause {
		absclause = append(absclause, int(math.Abs(float64(c))))
	}
	absclause = append(absclause, form.Nv)
	form.Nv = slices.Max(absclause)
	form.FromClouses = append(form.FromClouses, clause)
}
