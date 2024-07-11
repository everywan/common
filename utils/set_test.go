package utils

import (
	"sort"
	"testing"
)

func TestSet(t *testing.T) {
	set := NewSet[int]()
	set.Add(1)
	set.Add(1)
	set.Add(2)
	set.Add(3)
	l := set.ToList()
	if len(l) != 3 {
		t.Errorf("set.ToList() len=%d, want 3", len(l))
	}
	sort.Slice(l, func(i, j int) bool {
		return l[i] < l[j]
	})
	if l[0] != 1 || l[1] != 2 || l[2] != 3 {
		t.Errorf("set.ToList()=%v, want [1,2,3]", l)
	}
}
