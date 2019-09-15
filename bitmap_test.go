package bitmap_test

import (
	"bitmap"
	"testing"
)

func TestAdd(t *testing.T) {
	b := bitmap.New()
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if b.Len() != 4 && b.String() != "{0 1 2 10000}" {
		t.Errorf("TestAdd failed. Expected {0 1 2 10000}, Got %s", b.String())
	}
}

func TestHas(t *testing.T) {
	b := bitmap.New()
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if b.Has(-1) || !b.Has(0) || !b.Has(1) || !b.Has(2) || b.Has(3) || !b.Has(10000) {
		t.Errorf("TestHas failed.")
	}
}

func TestRemove(t *testing.T) {
	b := bitmap.New()
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(2)
	b.Add(4)
	b.Remove(3)
	if b.String() != "{0 1 2 4}" {
		t.Errorf("TestRemove failed. Expected {0 1 2 4}, Got %s", b.String())
	}
	b.Remove(4)
	if b.Has(4) {
		t.Errorf("TestRemove failed. Expected false, Got true")
	}
	b.Remove(-1)
	b.Remove(1)
	b.Remove(2)
	b.Remove(0)
	if b.String() != "{}" {
		t.Errorf("TestRemove failed. Expected {}, Got %s", b.String())
	}
}

func TestClear(t *testing.T) {
	b := bitmap.New()
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	b.Clear()
	if b.Has(0) {
		t.Errorf("TestClear failed")
	}
}

func TestLen(t *testing.T) {
	b := bitmap.New()
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if b.Len() != 4 {
		t.Errorf("TestLen failed.")
	}
}

func TestCopy(t *testing.T) {
	b := bitmap.New()
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	c := b.Copy()
	if !c.Has(0) || !c.Has(1) || !c.Has(2) {
		t.Errorf("TestCopy failed.")
	}
	b.Remove(2)
	if !c.Has(2) {
		t.Errorf("TestCopy failed.")
	}
}

func TestSets(t *testing.T) {
	b := bitmap.New()
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(3)
	b.Add(4)
	c := bitmap.New()
	c.Add(3)
	c.Add(4)
	c.Add(5)
	c.Add(6)
	c.Add(10000)
	bb := b.Copy()
	bb.Union(c)
	if bb.String() != "{0 1 2 3 4 5 6 10000}" {
		t.Errorf("TestSets Union failed. Expected {0 1 2 3 4 5 6 10000}, Got %s", bb.String())
	}
	bb = b.Copy()
	bb.Add(100001)
	bb.Intersect(c)
	if bb.String() != "{3 4}" {
		t.Errorf("TestSets Intersect failed. Expected {3 4}, Got %s", bb.String())
	}
	bb = b.Copy()
	bb.Except(c)
	if bb.String() != "{0 1 2}" {
		t.Errorf("TestSets Except failed. Expected {0 1 2}, Got %s", bb.String())
	}
	bb = b.Copy()
	bb.SymExcept(c)
	if bb.String() != "{0 1 2 5 6 10000}" {
		t.Errorf("TestSets SymExcept failed. Expected {0 1 2 5 6 10000}, Got %s", bb.String())
	}
}

func BenchmarkBitmap(b *testing.B) {
	bm := bitmap.New()
	const memory = 1000000000
	for i := 0; i < b.N; i++ {
		bm.Add(i % memory)
		bm.Has(i % memory)
		bm.Remove(i % memory)
	}
}
