package bitmap

import (
	"bytes"
	"fmt"
)

type bitInt uint

const (
	bitSize    = 32 << (^uint(0) >> 63)
	bitmapSize = 4
)

// Bitmap is a bitSet
type Bitmap struct {
	len   int
	words []bitInt
}

// New return a new bitmap
func New() *Bitmap {
	return &Bitmap{
		len:   0,
		words: make([]bitInt, bitmapSize),
	}
}

// String return formated string of bitmap
func (b *Bitmap) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range b.words {
		if word != 0 {
			for j := 0; j < bitSize; j++ {
				if word&(1<<uint(j)) != 0 {
					if buf.Len() > len("{") {
						buf.WriteByte(' ')
					}
					fmt.Fprintf(&buf, "%d", bitSize*i+j)
				}
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Len return numbers in bitmap
func (b *Bitmap) Len() int {
	return b.len
}

// Has return true if x is in the bitmap
func (b *Bitmap) Has(x int) bool {
	if x < 0 {
		return false
	}
	word, bit := x/bitSize, uint(x%bitSize)
	return word < len(b.words) && b.words[word]&(1<<bit) != 0
}

// Add add x to the bitmap
func (b *Bitmap) Add(x int) {
	if x < 0 {
		return
	}
	word, bit := x/bitSize, uint(x%bitSize)
	if word >= len(b.words) {
		b.words = append(b.words, make([]bitInt, word+1-len(b.words))...)
	}
	num := bitInt(1 << bit)
	if b.words[word]&num == 0 {
		b.len++
		b.words[word] |= num
	}
}

// Remove remove x in bitmap
func (b *Bitmap) Remove(x int) {
	if x < 0 {
		return
	}
	word, bit := x/bitSize, uint(x%bitSize)
	if word < len(b.words) {
		num := bitInt(1 << bit)
		if b.words[word]&num != 0 {
			b.len--
			b.words[word] &^= num
		}
	}
}

// Clear make the bitmap empty
func (b *Bitmap) Clear() {
	*b = *New()
}

// Copy return a copy bitmap
func (b *Bitmap) Copy() *Bitmap {
	new := Bitmap{}
	new.len = b.len
	new.words = make([]bitInt, len(b.words))
	copy(new.words, b.words)
	return &new
}

// Union b = b | c
// elements in b or c
func (b *Bitmap) Union(c *Bitmap) {
	length := 0
	for i, cword := range c.words {
		if i >= len(b.words) {
			length = i
			break
		}
		b.words[i] |= cword
	}
	if length != 0 {
		b.words = append(b.words, c.words[length:]...)
	}
}

// Intersect b = b & c
// elements both in b and c
func (b *Bitmap) Intersect(c *Bitmap) {
	for i, cwords := range c.words {
		if i >= len(b.words) {
			break
		}
		b.words[i] &= cwords
	}
	if len(c.words) < len(b.words) {
		b.words = b.words[:len(c.words)]
	}
}

// Except b = b - c
// elements only in b
func (b *Bitmap) Except(c *Bitmap) {
	for i, cwords := range c.words {
		if i >= len(b.words) {
			break
		}
		b.words[i] &^= cwords
	}
}

// SymExcept b = (b - c) | (c - b)
// elements only in b or only in c
func (b *Bitmap) SymExcept(c *Bitmap) {
	length := 0
	for i, cwords := range c.words {
		if i >= len(b.words) {
			length = i
			break
		}
		b.words[i] = (b.words[i] &^ cwords) | (cwords &^ b.words[i])
	}
	if length != 0 {
		b.words = append(b.words, c.words[length:]...)
	}
}
