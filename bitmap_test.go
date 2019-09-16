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
		t.Errorf("TestLen Add failed.")
	}
	b.Remove(2)
	b.Remove(3)
	if b.Len() != 3 {
		t.Errorf("TestLen Remove failed.")
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
	const memory = 100000000
	for i := 0; i < b.N; i++ {
		bm.Add(i % memory)
		bm.Has(i % memory)
		bm.Remove(i % memory)
	}
}

func TestCAdd(t *testing.T) {
	b := bitmap.NewC(3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if b.Len() != 4 && b.String() != "{0 1 2 10000}" {
		t.Errorf("TestCAdd failed. Expected {0 1 2 10000}, Got %s", b.String())
	}
}

func TestCHas(t *testing.T) {
	b := bitmap.NewC(3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if b.Has(-1) || !b.Has(0) || !b.Has(1) || !b.Has(2) || b.Has(3) || !b.Has(10000) {
		t.Errorf("TestCHas failed.")
	}
}

func TestCRemove(t *testing.T) {
	b := bitmap.NewC(3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(2)
	b.Add(4)
	b.Remove(3)
	if b.String() != "{0 1 2 4}" {
		t.Errorf("TestCRemove failed. Expected {0 1 2 4}, Got %s", b.String())
	}
	b.Remove(4)
	if b.Has(4) {
		t.Errorf("TestCRemove failed. Expected false, Got true")
	}
	b.Remove(-1)
	b.Remove(1)
	b.Remove(2)
	b.Remove(2)
	b.Remove(0)
	if b.String() != "{}" {
		t.Errorf("TestCRemove failed. Expected {}, Got %s", b.String())
	}
}

func TestCClear(t *testing.T) {
	b := bitmap.NewC(3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	b.Clear()
	if b.Has(0) {
		t.Errorf("TestCClear failed")
	}
}

func TestCLen(t *testing.T) {
	b := bitmap.NewC(3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if b.Len() != 4 {
		t.Errorf("TestCLen Add failed.")
	}
	b.Remove(2)
	b.Remove(3)
	if b.Len() != 3 {
		t.Errorf("TestCLen Remove failed.")
	}
}

func TestCCopy(t *testing.T) {
	b := bitmap.NewC(3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	c := b.Copy()
	if !c.Has(0) || !c.Has(1) || !c.Has(2) {
		t.Errorf("TestCCopy failed.")
	}
	b.Remove(2)
	if !c.Has(2) {
		t.Errorf("TestCCopy failed.")
	}
}

func TestCount(t *testing.T) {
	b := bitmap.NewC(4)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(2)
	b.Add(2)
	b.Add(2)
	b.Add(2)
	b.Add(2)
	b.Add(1)
	b.Add(10000)
	if b.Count(-1) != 0 || b.Count(2) != 4 || b.Count(1) != 2 {
		t.Errorf("TestCount Add failed.")
	}
	b.Remove(2)
	if b.Count(2) != 3 {
		t.Errorf("TestCount Remove failed.")
	}
	b.RemoveAll(1)
	if b.Count(1) != 0 {
		t.Errorf("TestCount RemoveAll failed.")
	}
}

func BenchmarkCBitmap(b *testing.B) {
	bm := bitmap.NewC(3)
	const memory = 100000000
	for i := 0; i < b.N; i++ {
		bm.Add(i % memory)
		bm.Has(i % memory)
		bm.Remove(i % memory)
	}
}

func TestRAdd(t *testing.T) {
	b := bitmap.NewR(-1, 3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if b.Len() != 4 && b.String() != "{-1 0 1 2}" {
		t.Errorf("TestRAdd failed. Expected {-1 0 1 2}, Got %s", b.String())
	}
}

func TestRHas(t *testing.T) {
	b := bitmap.NewR(-2, 10000)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if !b.Has(-1) || !b.Has(0) || !b.Has(1) || !b.Has(2) || b.Has(3) || b.Has(10000) {
		t.Errorf("TestRHas failed.")
	}
}

func TestRRemove(t *testing.T) {
	b := bitmap.NewR(-1, 5)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(2)
	b.Add(4)
	b.Remove(3)
	if b.String() != "{-1 0 1 2 4}" {
		t.Errorf("TestRRemove failed. Expected {-1 0 1 2 4}, Got %s", b.String())
	}
	b.Remove(4)
	if b.Has(4) {
		t.Errorf("TestRRemove failed. Expected false, Got true")
	}
	b.Remove(-1)
	b.Remove(1)
	b.Remove(2)
	b.Remove(0)
	if b.String() != "{}" {
		t.Errorf("TestRRemove failed. Expected {}, Got %s", b.String())
	}
}

func TestRClear(t *testing.T) {
	b := bitmap.NewR(-1, 10)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	b.Clear()
	if b.Has(0) {
		t.Errorf("TestRClear failed")
	}
}

func TestRLen(t *testing.T) {
	b := bitmap.NewR(-1, 3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if b.Len() != 4 {
		t.Errorf("TestRLen Add failed.")
	}
	b.Remove(2)
	b.Remove(3)
	if b.Len() != 3 {
		t.Errorf("TestRLen Remove failed.")
	}
}

func TestRCopy(t *testing.T) {
	b := bitmap.NewR(-1, 3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	c := b.Copy()
	if !c.Has(0) || !c.Has(1) || !c.Has(2) {
		t.Errorf("TestRCopy failed.")
	}
	b.Remove(2)
	if !c.Has(2) {
		t.Errorf("TestRCopy failed.")
	}
}

func TestRSets(t *testing.T) {
	b := bitmap.NewR(-1, 10001)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(3)
	b.Add(4)
	c := bitmap.NewR(-1, 10001)
	c.Add(3)
	c.Add(4)
	c.Add(5)
	c.Add(6)
	c.Add(10000)
	bb := b.Copy()
	bb.Union(c)
	if bb.String() != "{-1 0 1 2 3 4 5 6 10000}" {
		t.Errorf("TestRSets Union failed. Expected {-1 0 1 2 3 4 5 6 10000}, Got %s", bb.String())
	}
	bb = b.Copy()
	bb.Add(100001)
	bb.Intersect(c)
	if bb.String() != "{3 4}" {
		t.Errorf("TestRSets Intersect failed. Expected {3 4}, Got %s", bb.String())
	}
	bb = b.Copy()
	bb.Except(c)
	if bb.String() != "{-1 0 1 2}" {
		t.Errorf("TestRSets Except failed. Expected {-1 0 1 2}, Got %s", bb.String())
	}
	bb = b.Copy()
	bb.SymExcept(c)
	if bb.String() != "{-1 0 1 2 5 6 10000}" {
		t.Errorf("TestRSets SymExcept failed. Expected {-1 0 1 2 5 6 10000}, Got %s", bb.String())
	}
}

func BenchmarkRBitmap(b *testing.B) {
	bm := bitmap.New()
	const memory = 100000000
	for i := 0; i < b.N; i++ {
		bm.Add(i % memory)
		bm.Has(i % memory)
		bm.Remove(i % memory)
	}
}

func TestRCAdd(t *testing.T) {
	b := bitmap.NewRC(-1, 3, 3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if b.Len() != 4 && b.String() != "{-1 0 1 2}" {
		t.Errorf("TestRCAdd failed. Expected {-1 0 1 2}, Got %s", b.String())
	}
}

func TestRCHas(t *testing.T) {
	b := bitmap.NewRC(-1, 3, 3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if !b.Has(-1) || !b.Has(0) || !b.Has(1) || !b.Has(2) || b.Has(3) || b.Has(10000) {
		t.Errorf("TestRCHas failed.")
	}
}

func TestRCRemove(t *testing.T) {
	b := bitmap.NewRC(0, 5, 3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(2)
	b.Add(4)
	b.Remove(3)
	if b.String() != "{0 1 2 4}" {
		t.Errorf("TestRCRemove failed. Expected {0 1 2 4}, Got %s", b.String())
	}
	b.Remove(4)
	if b.Has(4) {
		t.Errorf("TestRCRemove failed. Expected false, Got true")
	}
	b.Remove(-1)
	b.Remove(1)
	b.Remove(2)
	b.Remove(2)
	b.Remove(0)
	if b.String() != "{}" {
		t.Errorf("TestRCRemove failed. Expected {}, Got %s", b.String())
	}
}

func TestRCClear(t *testing.T) {
	b := bitmap.NewRC(0, 20, 3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	b.Clear()
	if b.Has(0) {
		t.Errorf("TestRCClear failed")
	}
}

func TestRCLen(t *testing.T) {
	b := bitmap.NewRC(-1, 3, 3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if b.Len() != 4 {
		t.Errorf("TestRCLen Add failed.")
	}
	b.Remove(2)
	b.Remove(3)
	if b.Len() != 3 {
		t.Errorf("TestRCLen Remove failed.")
	}
}

func TestRCCopy(t *testing.T) {
	b := bitmap.NewRC(-1, 3, 3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	c := b.Copy()
	if !c.Has(0) || !c.Has(1) || !c.Has(2) {
		t.Errorf("TestRCCopy failed.")
	}
	b.Remove(2)
	if !c.Has(2) {
		t.Errorf("TestRCCopy failed.")
	}
}

func TestRCount(t *testing.T) {
	b := bitmap.NewRC(-1, 3, 3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(2)
	b.Add(2)
	b.Add(1)
	b.Add(10000)
	if b.Count(-1) != 1 || b.Count(2) != 3 || b.Count(1) != 2 {
		t.Errorf("TestRCount Add failed.")
	}
	b.Remove(2)
	if b.Count(2) != 2 {
		t.Errorf("TestRCount Remove failed.")
	}
	b.RemoveAll(1)
	if b.Count(1) != 0 {
		t.Errorf("TestRCount RemoveAll failed.")
	}
}

func BenchmarkRCBitmap(b *testing.B) {
	const memory = 100000000
	bm := bitmap.NewRC(0, memory, 3)
	for i := 0; i < b.N; i++ {
		bm.Add(i % memory)
		bm.Has(i % memory)
		bm.Remove(i % memory)
	}
}
