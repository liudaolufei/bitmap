package bitmap_test

import (
	"bitmap"
	"testing"
)

func TestCAdd(t *testing.T) {
	b := bitmap.NewC(3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(10000)
	if b.Len() != 4 && b.String() != "{0 1 2 10000}" {
		t.Errorf("TestAdd failed. Expected {0 1 2 10000}, Got %s", b.String())
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
		t.Errorf("TestHas failed.")
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
		t.Errorf("TestRemove failed. Expected {0 1 2 4}, Got %s", b.String())
	}
	b.Remove(4)
	if b.Has(4) {
		t.Errorf("TestRemove failed. Expected false, Got true")
	}
	b.Remove(-1)
	b.Remove(1)
	b.Remove(2)
	b.Remove(2)
	b.Remove(0)
	if b.String() != "{}" {
		t.Errorf("TestRemove failed. Expected {}, Got %s", b.String())
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
		t.Errorf("TestClear failed")
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
		t.Errorf("TestLen Add failed.")
	}
	b.Remove(2)
	b.Remove((3))
	if b.Len() != 3 {
		t.Errorf("TestLen Remove failed.")
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
		t.Errorf("TestCopy failed.")
	}
	b.Remove(2)
	if !c.Has(2) {
		t.Errorf("TestCopy failed.")
	}
}

func TestCount(t *testing.T) {
	b := bitmap.NewC(3)
	b.Add(-1)
	b.Add(0)
	b.Add(1)
	b.Add(2)
	b.Add(2)
	b.Add(2)
	b.Add(1)
	b.Add(10000)
	if b.Count(-1) != 0 || b.Count(2) != 3 || b.Count(1) != 2 {
		t.Errorf("TestCount Add failed.")
	}
	b.Remove(2)
	if b.Count(2) != 2 {
		t.Errorf("TestCount Remove failed.")
	}
	b.RemoveAll(1)
	if b.Count(1) != 0 {
		t.Errorf("TestCount RemoveAll failed.")
	}
}

func BenchmarkCBitmap(b *testing.B) {
	bm := bitmap.NewC(3)
	const memory = 1000000000
	for i := 0; i < b.N; i++ {
		bm.Add(i % memory)
		bm.Has(i % memory)
		bm.Remove(i % memory)
	}
}
