package basal

import (
	"reflect"
)

type SortedSet struct {
	buf []interface{}
	cmp func(min, max interface{}) int // max > min return 1, max == min return 0, max < min return -1
}

func NewSortedSet(cmp func(min, max interface{}) int) *SortedSet {
	return &SortedSet{cmp: cmp}
}

func (my *SortedSet) Add(v interface{}) {
	index, found := my.binaryFind(v)
	if found {
		return
	}
	my.buf = append(my.buf, v)
	copy(my.buf[index+1:], my.buf[index:])
	my.buf[index] = v
}

func (my *SortedSet) binaryFind(value interface{}) (index int, found bool) {
	length := len(my.buf)
	if length == 0 {
		return 0, false
	}
	start := 0
	end := length - 1
	var cmp int
	for {
		index = start + (end-start)/2
		cmp = my.cmp(value, my.buf[index])
		if cmp == 1 {
			end = index - 1
		} else if cmp == -1 {
			start = index + 1
		} else {
			return index, true
		}
		if start > end {
			return start, false
		}
	}
}

func (my *SortedSet) Front() (front interface{}, found bool) {
	if len(my.buf) > 0 {
		return my.buf[0], true
	} else {
		return nil, false
	}
}

func (my *SortedSet) Back() (back interface{}, found bool) {
	length := len(my.buf)
	if length > 0 {
		return my.buf[length-1], true
	} else {
		return nil, false
	}
}

func (my *SortedSet) Get(index int) (v interface{}, found bool) {
	length := len(my.buf)
	if length > 0 {
		if index < 0 || index >= length {
			return nil, false
		}
		return my.buf[index], true
	} else {
		return nil, false
	}
}

func (my *SortedSet) Del(index int) bool {
	if index < 0 || index >= len(my.buf) {
		return false
	}
	my.buf = append(my.buf[:index], my.buf[index+1:]...)
	return true
}

func (my *SortedSet) Find(v interface{}) (index int, found bool) {
	return my.binaryFind(v)
}

func (my *SortedSet) Remove(v interface{}) bool {
	index, found := my.binaryFind(v)
	if found {
		my.buf = append(my.buf[:index], my.buf[index+1:]...)
		return true
	} else {
		return false
	}
}

func (my *SortedSet) Clear() {
	my.buf = my.buf[:0]
}

func (my *SortedSet) can(b *SortedSet) bool {
	aLen, bLen := len(my.buf), len(b.buf)
	if bLen < 1 {
		return true
	}
	if aLen > 0 {
		if reflect.TypeOf(my.buf[0]) == reflect.TypeOf(b.buf[0]) {
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}

func (my *SortedSet) Union(b *SortedSet) *SortedSet {
	if !my.can(b) {
		return nil
	}
	c := DeepCopy(my).(*SortedSet)
	for _, value := range b.buf {
		c.Add(value)
	}
	return c
}

func (my *SortedSet) Difference(b *SortedSet) *SortedSet {
	if !my.can(b) {
		return nil
	}
	c := DeepCopy(my).(*SortedSet)
	for _, value := range b.buf {
		c.Remove(value)
	}
	return c
}

func (my *SortedSet) Intersection(b *SortedSet) *SortedSet {
	if !my.can(b) {
		return nil
	}
	c := NewSortedSet(my.cmp)
	for _, value := range my.buf {
		_, found := b.Find(value)
		if found {
			c.Add(value)
		}
	}
	return c
}

type IntSortedSet struct {
	SortedSet
}

func NewIntSortedSet(reverse bool) *IntSortedSet {
	if reverse {
		return &IntSortedSet{SortedSet{cmp: func(min, max interface{}) int {
			x, y := min.(int), max.(int)
			if x < y {
				return 1
			} else if x > y {
				return -1
			} else {
				return 0
			}
		}}}
	} else {
		return &IntSortedSet{SortedSet{cmp: func(min, max interface{}) int {
			x, y := min.(int), max.(int)
			if x < y {
				return -1
			} else if x > y {
				return 1
			} else {
				return 0
			}
		}}}
	}
}
