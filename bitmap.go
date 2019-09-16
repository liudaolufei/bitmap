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
func (n *NBitmap) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range n.words {
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
func (n *NBitmap) Len() int {
	return n.len
}

// Has return true if x is in the bitmap
func (n *NBitmap) Has(x int) bool {
	if x < 0 {
		return false
	}
	word, bit := x/bitSize, bitInt(x%bitSize)
	return word < len(n.words) && n.words[word]&(1<<bit) != 0
}

// Add add x to the bitmap
func (n *NBitmap) Add(x int) {
	if x < 0 {
		return
	}
	word, bit := x/bitSize, bitInt(x%bitSize)
	if word >= len(n.words) {
		n.words = append(n.words, make([]bitInt, word+1-len(n.words))...)
	}
	num := bitInt(1 << bit)
	if n.words[word]&num == 0 {
		n.len++
		n.words[word] |= num
	}
}

// Remove remove x in bitmap
func (n *NBitmap) Remove(x int) {
	if x < 0 {
		return
	}
	word, bit := x/bitSize, bitInt(x%bitSize)
	if word < len(n.words) {
		num := bitInt(1 << bit)
		if n.words[word]&num != 0 {
			n.len--
			n.words[word] &^= num
		}
	}
}

// Clear make the bitmap empty
func (n *NBitmap) Clear() {
	*n = *New()
}

// Copy return a copy bitmap
func (n *NBitmap) Copy() *NBitmap {
	new := NBitmap{}
	new.len = n.len
	new.words = make([]bitInt, len(n.words))
	copy(new.words, n.words)
	return &new
}

// Union n = n | c
// elements in n or c
func (n *NBitmap) Union(c *NBitmap) {
	length := 0
	for i, cword := range c.words {
		if i >= len(n.words) {
			length = i
			break
		}
		n.words[i] |= cword
	}
	if length != 0 {
		n.words = append(n.words, c.words[length:]...)
	}
}

// Intersect n = n & c
// elements both in n and c
func (n *NBitmap) Intersect(c *NBitmap) {
	for i, cwords := range c.words {
		if i >= len(n.words) {
			break
		}
		n.words[i] &= cwords
	}
	if len(c.words) < len(n.words) {
		n.words = n.words[:len(c.words)]
	}
}

// Except n = n - c
// elements only in n
func (n *NBitmap) Except(c *NBitmap) {
	for i, cwords := range c.words {
		if i >= len(n.words) {
			break
		}
		n.words[i] &^= cwords
	}
}

// SymExcept n = (n - c) | (c - n)
// elements only in n or only in c
func (n *NBitmap) SymExcept(c *NBitmap) {
	length := 0
	for i, cwords := range c.words {
		if i >= len(n.words) {
			length = i
			break
		}
		n.words[i] = (n.words[i] &^ cwords) | (cwords &^ n.words[i])
	}
	if length != 0 {
		n.words = append(n.words, c.words[length:]...)
	}
}

// RBitmap is a bitSet count in [start, end)
type RBitmap struct {
	len   int
	start int
	end   int
	words []bitInt
}

// NewR return a new bitmap, count in [start, end)
func NewR(start int, end int) *RBitmap {
	if start >= end {
		return nil
	}
	return &RBitmap{
		len:   0,
		start: start,
		end:   end,
		words: make([]bitInt, bitmapSize),
	}
}

// String return formated string of bitmap
func (r *RBitmap) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range r.words {
		if word != 0 {
			for j := 0; j < bitSize; j++ {
				if word&(1<<bitInt(j)) != 0 {
					if buf.Len() > len("{") {
						buf.WriteByte(' ')
					}
					fmt.Fprintf(&buf, "%d", r.start+bitSize*i+j)
				}
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Len return numbers in bitmap
func (r *RBitmap) Len() int {
	return r.len
}

// Has return true if x is in the bitmap
func (r *RBitmap) Has(x int) bool {
	if x < r.start || x >= r.end {
		return false
	}
	x -= r.start
	word, bit := x/bitSize, bitInt(x%bitSize)
	return word < len(r.words) && r.words[word]&(1<<bit) != 0
}

// Add add x to the bitmap
func (r *RBitmap) Add(x int) {
	if x < r.start || x >= r.end {
		return
	}
	x -= r.start
	word, bit := x/bitSize, bitInt(x%bitSize)
	if word >= len(r.words) {
		r.words = append(r.words, make([]bitInt, word+1-len(r.words))...)
	}
	num := bitInt(1 << bit)
	if r.words[word]&num == 0 {
		r.len++
		r.words[word] |= num
	}
}

// Remove remove x in bitmap
func (r *RBitmap) Remove(x int) {
	if x < r.start || x >= r.end {
		return
	}
	x -= r.start
	word, bit := x/bitSize, bitInt(x%bitSize)
	if word < len(r.words) {
		num := bitInt(1 << bit)
		if r.words[word]&num != 0 {
			r.len--
			r.words[word] &^= num
		}
	}
}

// Clear make the bitmap empty
func (r *RBitmap) Clear() {
	*r = *NewR(r.start, r.end)
}

// Copy return a copy bitmap
func (r *RBitmap) Copy() *RBitmap {
	new := RBitmap{}
	new.len = r.len
	new.start = r.start
	new.end = r.end
	new.words = make([]bitInt, len(r.words))
	copy(new.words, r.words)
	return &new
}

// Union r = r | c
// elements in r or c
// r must have the same range of c
func (r *RBitmap) Union(c *RBitmap) {
	if r.start != c.start || r.end != c.end {
		return
	}
	length := 0
	for i, cword := range c.words {
		if i >= len(r.words) {
			length = i
			break
		}
		r.words[i] |= cword
	}
	if length != 0 {
		r.words = append(r.words, c.words[length:]...)
	}
}

// Intersect r = r & c
// elements both in r and c
// r must have the same range of c
func (r *RBitmap) Intersect(c *RBitmap) {
	if r.start != c.start || r.end != c.end {
		return
	}
	for i, cwords := range c.words {
		if i >= len(r.words) {
			break
		}
		r.words[i] &= cwords
	}
	if len(c.words) < len(r.words) {
		r.words = r.words[:len(c.words)]
	}
}

// Except r = r - c
// elements only in r
// r must have the same range of c
func (r *RBitmap) Except(c *RBitmap) {
	if r.start != c.start || r.end != c.end {
		return
	}
	for i, cwords := range c.words {
		if i >= len(r.words) {
			break
		}
		r.words[i] &^= cwords
	}
}

// SymExcept r = (r - c) | (c - r)
// elements only in r or only in c
// r must have the same range of c
func (r *RBitmap) SymExcept(c *RBitmap) {
	if r.start != c.start || r.end != c.end {
		return
	}
	length := 0
	for i, cwords := range c.words {
		if i >= len(r.words) {
			length = i
			break
		}
		r.words[i] = (r.words[i] &^ cwords) | (cwords &^ r.words[i])
	}
	if length != 0 {
		r.words = append(r.words, c.words[length:]...)
	}
}

// CBitmap is a bitSet
type CBitmap struct {
	len     int
	n       bitInt
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
	c := CBitmap{}
	c.len = 0
	c.n = bitInt(n)
	c.numSize = numSize + 1
	c.bitSize = bitSize / c.numSize
	c.mask = (1 << bitInt(c.numSize)) - 1
	c.words = make([]bitInt, bitmapSize)
	return &c
}

// String return formated string of bitmap
func (c *CBitmap) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range c.words {
		if word != 0 {
			for j := 0; j < c.bitSize; j += c.numSize {
				if word&(bitInt(c.mask<<bitInt(j))) != 0 {
					if buf.Len() > len("{") {
						buf.WriteByte(' ')
					}
					fmt.Fprintf(&buf, "%d", c.bitSize*i+j/c.numSize)
				}
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Len return numbers in bitmap
func (c *CBitmap) Len() int {
	return c.len
}

// Has return true if x is in the bitmap
func (c *CBitmap) Has(x int) bool {
	if x < 0 {
		return false
	}
	word, bit := x/c.bitSize, bitInt(x%c.bitSize*c.numSize)
	return word < len(c.words) && c.words[word]&(1<<bit) != 0
}

// Add add x to the bitmap
func (c *CBitmap) Add(x int) {
	if x < 0 {
		return
	}
	word, bit := x/c.bitSize, bitInt(x%c.bitSize*c.numSize)
	if word >= len(c.words) {
		c.words = append(c.words, make([]bitInt, word+1-len(c.words))...)
	}
	numSize := bitInt(c.mask << bit)
	if c.words[word]&numSize == 0 {
		c.len++
	}
	if ((c.words[word] & numSize) >> bit) < c.n {
		c.words[word] += 1 << bit
	}
}

// Remove remove x in bitmap
func (c *CBitmap) Remove(x int) {
	if x < 0 {
		return
	}
	word, bit := x/c.bitSize, bitInt(x%c.bitSize*c.numSize)
	if word < len(c.words) {
		numSize := bitInt(c.mask << bit)
		if c.words[word]&numSize != 0 {
			c.words[word] -= 1 << bit
			if c.words[word]&numSize == 0 {
				c.len--
			}
		}
	}
}

// Count return the numSize of x elements
func (c *CBitmap) Count(x int) int {
	if x < 0 {
		return 0
	}
	word, bit := x/c.bitSize, bitInt(x%c.bitSize*c.numSize)
	if word >= len(c.words) {
		return 0
	}
	numSize := bitInt(c.mask << bit)
	return int((numSize & c.words[word]) >> bit)
}

// RemoveAll remove x in bitmap
func (c *CBitmap) RemoveAll(x int) {
	if x < 0 {
		return
	}
	word, bit := x/c.bitSize, bitInt(x%c.bitSize*c.numSize)
	if word < len(c.words) {
		numSize := bitInt(c.mask << bit)
		if c.words[word]&numSize != 0 {
			c.len--
			c.words[word] &^= numSize
		}
	}
}

// Clear make the bitmap empty
func (c *CBitmap) Clear() {
	*c = *NewC(int(c.n))
}

// Copy return a copy bitmap
func (c *CBitmap) Copy() *CBitmap {
	new := CBitmap{}
	new.len = c.len
	new.n = c.n
	new.mask = c.mask
	new.bitSize = c.bitSize
	new.numSize = c.numSize
	new.words = make([]bitInt, len(c.words))
	copy(new.words, c.words)
	return &new
}

// RCBitmap is a bitSet count in [start, end)
type RCBitmap struct {
	len        int
	n          bitInt
	start, end int
	bitSize    int
	mask       int
	numSize    int
	words      []bitInt
}

// NewRC return a new bitmap count [start, end)
func NewRC(start int, end int, n int) *RCBitmap {
	if n <= 0 || n > (1<<(bitSize-1)-1) || start >= end {
		return nil
	}
	numSize := int(math.Log2(float64(n)))
	rc := RCBitmap{}
	rc.len = 0
	rc.n = bitInt(n)
	rc.start, rc.end = start, end
	rc.numSize = numSize + 1
	rc.bitSize = bitSize / rc.numSize
	rc.mask = (1 << bitInt(rc.numSize)) - 1
	rc.words = make([]bitInt, bitmapSize)
	return &rc
}

// String return formated string of bitmap
func (rc *RCBitmap) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range rc.words {
		if word != 0 {
			for j := 0; j < rc.bitSize; j += rc.numSize {
				if word&(bitInt(rc.mask<<bitInt(j))) != 0 {
					if buf.Len() > len("{") {
						buf.WriteByte(' ')
					}
					fmt.Fprintf(&buf, "%d", rc.start+rc.bitSize*i+j/rc.numSize)
				}
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

// Len return numbers in bitmap
func (rc *RCBitmap) Len() int {
	return rc.len
}

// Has return true if x is in the bitmap
func (rc *RCBitmap) Has(x int) bool {
	if x < rc.start || x >= rc.end {
		return false
	}
	x -= rc.start
	word, bit := x/rc.bitSize, bitInt(x%rc.bitSize*rc.numSize)
	return word < len(rc.words) && rc.words[word]&(1<<bit) != 0
}

// Add add x to the bitmap
func (rc *RCBitmap) Add(x int) {
	if x < rc.start || x >= rc.end {
		return
	}
	x -= rc.start
	word, bit := x/rc.bitSize, bitInt(x%rc.bitSize*rc.numSize)
	if word >= len(rc.words) {
		rc.words = append(rc.words, make([]bitInt, word+1-len(rc.words))...)
	}
	numSize := bitInt(rc.mask << bit)
	if rc.words[word]&numSize == 0 {
		rc.len++
	}
	if ((rc.words[word] & numSize) >> bit) < rc.n {
		rc.words[word] += 1 << bit
	}
}

// Remove remove x in bitmap
func (rc *RCBitmap) Remove(x int) {
	if x < rc.start || x >= rc.end {
		return
	}
	x -= rc.start
	word, bit := x/rc.bitSize, bitInt(x%rc.bitSize*rc.numSize)
	if word < len(rc.words) {
		numSize := bitInt(rc.mask << bit)
		if rc.words[word]&numSize != 0 {
			rc.words[word] -= 1 << bit
			if rc.words[word]&numSize == 0 {
				rc.len--
			}
		}
	}
}

// Count return the numSize of x elements
func (rc *RCBitmap) Count(x int) int {
	if x < rc.start || x >= rc.end {
		return 0
	}
	x -= rc.start
	word, bit := x/rc.bitSize, bitInt(x%rc.bitSize*rc.numSize)
	if word >= len(rc.words) {
		return 0
	}
	numSize := bitInt(rc.mask << bit)
	return int((numSize & rc.words[word]) >> bit)
}

// RemoveAll remove x in bitmap
func (rc *RCBitmap) RemoveAll(x int) {
	if x < rc.start || x >= rc.end {
		return
	}
	x -= rc.start
	word, bit := x/rc.bitSize, bitInt(x%rc.bitSize*rc.numSize)
	if word < len(rc.words) {
		numSize := bitInt(rc.mask << bit)
		if rc.words[word]&numSize != 0 {
			rc.len--
			rc.words[word] &^= numSize
		}
	}
}

// Clear make the bitmap empty
func (rc *RCBitmap) Clear() {
	*rc = *NewRC(rc.start, rc.end, int(rc.n))
}

// Copy return a copy bitmap
func (rc *RCBitmap) Copy() *RCBitmap {
	new := RCBitmap{}
	new.len = rc.len
	new.n = rc.n
	new.start, new.end = rc.start, rc.end
	new.mask = rc.mask
	new.bitSize = rc.bitSize
	new.numSize = rc.numSize
	new.words = make([]bitInt, len(rc.words))
	copy(new.words, rc.words)
	return &new
}
