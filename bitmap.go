package bitmap

import (
	"bytes"
	"fmt"
	"math"
)

type bitInt uint

const (
	bitSize    = 32 << (^bitInt(0) >> 63)
	bitmapSize = 4
)

// Bitmap is the interface of bitSet
type Bitmap interface {
	Add(x int)    // add x to bitmap
	Has(x int)    // return true if x is in bitmap
	Remove(x int) // remove x in bitmap
	Len() int     // return length of bitmap
	Clear()       // clear bitmap to free memory
}

// NBitmap is a normal bitSet
type NBitmap struct {
	len   int
	words []bitInt
}

// New return a new bitmap
func New() *NBitmap {
	return &NBitmap{
		len:   0,
		words: make([]bitInt, bitmapSize),
	}
}

// String return formated string of bitmap
func (b *NBitmap) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range b.words {
		if word != 0 {
			for j := 0; j < bitSize; j++ {
				if word&(1<<bitInt(j)) != 0 {
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
func (b *NBitmap) Len() int {
	return b.len
}

// Has return true if x is in the bitmap
func (b *NBitmap) Has(x int) bool {
	if x < 0 {
		return false
	}
	word, bit := x/bitSize, bitInt(x%bitSize)
	return word < len(b.words) && b.words[word]&(1<<bit) != 0
}

// Add add x to the bitmap
func (b *NBitmap) Add(x int) {
	if x < 0 {
		return
	}
	word, bit := x/bitSize, bitInt(x%bitSize)
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
func (b *NBitmap) Remove(x int) {
	if x < 0 {
		return
	}
	word, bit := x/bitSize, bitInt(x%bitSize)
	if word < len(b.words) {
		num := bitInt(1 << bit)
		if b.words[word]&num != 0 {
			b.len--
			b.words[word] &^= num
		}
	}
}

// Clear make the bitmap empty
func (b *NBitmap) Clear() {
	*b = *New()
}

// Copy return a copy bitmap
func (b *NBitmap) Copy() *NBitmap {
	new := NBitmap{}
	new.len = b.len
	new.words = make([]bitInt, len(b.words))
	copy(new.words, b.words)
	return &new
}

// Union b = b | c
// elements in b or c
func (b *NBitmap) Union(c *NBitmap) {
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
func (b *NBitmap) Intersect(c *NBitmap) {
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
func (b *NBitmap) Except(c *NBitmap) {
	for i, cwords := range c.words {
		if i >= len(b.words) {
			break
		}
		b.words[i] &^= cwords
	}
}

// SymExcept b = (b - c) | (c - b)
// elements only in b or only in c
func (b *NBitmap) SymExcept(c *NBitmap) {
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

// CBitmap is a bitSet
type CBitmap struct {
	len     int
	n       int
	bitSize int
	mask    int
	numSize int
	words   []bitInt
}

// NewC return a new bitmap
func NewC(n int) *CBitmap {
	if n <= 0 || n > (1<<(bitSize-1)-1) {
		return nil
	}
	numSize := int(math.Log2(float64(n)))
	cb := CBitmap{}
	cb.len = 0
	cb.n = n
	cb.numSize = numSize + 1
	cb.bitSize = bitSize / cb.numSize
	cb.mask = (1 << bitInt(cb.numSize)) - 1
	cb.words = make([]bitInt, bitmapSize)
	return &cb
}

// String return formated string of bitmap
func (cb *CBitmap) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range cb.words {
		if word != 0 {
			for j := 0; j < cb.bitSize; j += cb.numSize {
				if word&(bitInt(cb.mask<<bitInt(j))) != 0 {
					if buf.Len() > len("{") {
						buf.WriteByte(' ')
					}
					fmt.Fprintf(&buf, "%d", cb.bitSize*i+j/cb.numSize)
				}
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Len return numbers in bitmap
func (cb *CBitmap) Len() int {
	return cb.len
}

// Has return true if x is in the bitmap
func (cb *CBitmap) Has(x int) bool {
	if x < 0 {
		return false
	}
	word, bit := x/cb.bitSize, bitInt(x%cb.bitSize*cb.numSize)
	return word < len(cb.words) && cb.words[word]&(1<<bit) != 0
}

// Add add x to the bitmap
func (cb *CBitmap) Add(x int) {
	if x < 0 {
		return
	}
	word, bit := x/cb.bitSize, bitInt(x%cb.bitSize*cb.numSize)
	if word >= len(cb.words) {
		cb.words = append(cb.words, make([]bitInt, word+1-len(cb.words))...)
	}
	numSize := bitInt(cb.mask << bit)
	if cb.words[word]&numSize == 0 {
		cb.len++
	}
	if ((cb.words[word] & numSize) >> bit) != bitInt(cb.mask) {
		cb.words[word] += 1 << bit
	}
}

// Remove remove x in bitmap
func (cb *CBitmap) Remove(x int) {
	if x < 0 {
		return
	}
	word, bit := x/cb.bitSize, bitInt(x%cb.bitSize*cb.numSize)
	if word < len(cb.words) {
		numSize := bitInt(cb.mask << bit)
		if cb.words[word]&numSize != 0 {
			cb.words[word] -= 1 << bit
			if cb.words[word]&numSize == 0 {
				cb.len--
			}
		}
	}
}

// Count return the numSize of x elements
func (cb *CBitmap) Count(x int) int {
	if x < 0 {
		return 0
	}
	word, bit := x/cb.bitSize, bitInt(x%cb.bitSize*cb.numSize)
	if word >= len(cb.words) {
		return 0
	}
	numSize := bitInt(cb.mask << bit)
	return int((numSize & cb.words[word]) >> bit)
}

// RemoveAll remove x in bitmap
func (cb *CBitmap) RemoveAll(x int) {
	if x < 0 {
		return
	}
	word, bit := x/cb.bitSize, bitInt(x%cb.bitSize*cb.numSize)
	if word < len(cb.words) {
		numSize := bitInt(cb.mask << bit)
		if cb.words[word]&numSize != 0 {
			cb.len--
			cb.words[word] &^= numSize
		}
	}
}

// Clear make the bitmap empty
func (cb *CBitmap) Clear() {
	*cb = *NewC(cb.n)
}

// Copy return a copy bitmap
func (cb *CBitmap) Copy() *CBitmap {
	new := CBitmap{}
	new.len = cb.len
	new.n = cb.n
	new.mask = cb.mask
	new.bitSize = cb.bitSize
	new.numSize = cb.numSize
	new.words = make([]bitInt, len(cb.words))
	copy(new.words, cb.words)
	return &new
}
