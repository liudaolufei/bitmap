package bitmap

import (
	"bytes"
	"fmt"
)

// CBitmap is a bitSet
type CBitmap struct {
	len     int
	n       int
	bitSize int
	mask    int
	num     int
	words   []bitInt
}

// NewC return a new bitmap
func NewC(n int) *CBitmap {
	if n <= 0 || n > (1<<(bitSize-1)-1) {
		return nil
	}
	i := 0
	for i = bitSize - 1; i >= 0; i-- {
		if n&(1<<bitInt(i)) != 0 {
			break
		}
	}
	cb := CBitmap{}
	cb.len = 0
	cb.n = n
	cb.num = i + 1
	cb.bitSize = bitSize / cb.num
	cb.mask = (1 << bitInt(cb.num)) - 1
	cb.words = make([]bitInt, bitmapSize)
	return &cb
}

// String return formated string of bitmap
func (cb *CBitmap) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range cb.words {
		if word != 0 {
			for j := 0; j < cb.bitSize; j += cb.num {
				if word&(bitInt(cb.mask<<bitInt(j))) != 0 {
					if buf.Len() > len("{") {
						buf.WriteByte(' ')
					}
					fmt.Fprintf(&buf, "%d", cb.bitSize*i+j/cb.num)
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
	word, bit := x/cb.bitSize, bitInt(x%cb.bitSize*cb.num)
	return word < len(cb.words) && cb.words[word]&(1<<bit) != 0
}

// Add add x to the bitmap
func (cb *CBitmap) Add(x int) {
	if x < 0 {
		return
	}
	word, bit := x/cb.bitSize, bitInt(x%cb.bitSize*cb.num)
	if word >= len(cb.words) {
		cb.words = append(cb.words, make([]bitInt, word+1-len(cb.words))...)
	}
	num := bitInt(cb.mask << bit)
	if cb.words[word]&num == 0 {
		cb.len++
	}
	if ((cb.words[word] & num) >> bit) != bitInt(cb.mask) {
		cb.words[word] += 1 << bit
	}
}

// Remove remove x in bitmap
func (cb *CBitmap) Remove(x int) {
	if x < 0 {
		return
	}
	word, bit := x/cb.bitSize, bitInt(x%cb.bitSize*cb.num)
	if word < len(cb.words) {
		num := bitInt(cb.mask << bit)
		if cb.words[word]&num != 0 {
			cb.words[word] -= 1 << bit
			if cb.words[word]&num == 0 {
				cb.len--
			}
		}
	}
}

// Count return the num of x elements
func (cb *CBitmap) Count(x int) int {
	if x < 0 {
		return 0
	}
	word, bit := x/cb.bitSize, bitInt(x%cb.bitSize*cb.num)
	if word >= len(cb.words) {
		return 0
	}
	num := bitInt(cb.mask << bit)
	return int((num & cb.words[word]) >> bit)
}

// RemoveAll remove x in bitmap
func (cb *CBitmap) RemoveAll(x int) {
	if x < 0 {
		return
	}
	word, bit := x/cb.bitSize, bitInt(x%cb.bitSize*cb.num)
	if word < len(cb.words) {
		num := bitInt(cb.mask << bit)
		if cb.words[word]&num != 0 {
			cb.len--
			cb.words[word] &^= num
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
	new.num = cb.num
	new.words = make([]bitInt, len(cb.words))
	copy(new.words, cb.words)
	return &new
}

// // Union cb = cb | c
// // elements in cb or c
// func (cb *CBitmap) Union(c *CBitmap) {
// 	length := 0
// 	for i, cword := range c.words {
// 		if i >= len(cb.words) {
// 			length = i
// 			break
// 		}
// 		cb.words[i] |= cword
// 	}
// 	if length != 0 {
// 		cb.words = append(cb.words, c.words[length:]...)
// 	}
// }

// // Intersect cb = cb & c
// // elements both in cb and c
// func (cb *CBitmap) Intersect(c *CBitmap) {
// 	for i, cwords := range c.words {
// 		if i >= len(cb.words) {
// 			break
// 		}
// 		cb.words[i] &= cwords
// 	}
// 	if len(c.words) < len(cb.words) {
// 		cb.words = cb.words[:len(c.words)]
// 	}
// }

// // Except cb = cb - c
// // elements only in cb
// func (cb *CBitmap) Except(c *CBitmap) {
// 	for i, cwords := range c.words {
// 		if i >= len(cb.words) {
// 			break
// 		}
// 		cb.words[i] &^= cwords
// 	}
// }

// // SymExcept cb = (cb - c) | (c - cb)
// // elements only in cb or only in c
// func (cb *CBitmap) SymExcept(c *CBitmap) {
// 	length := 0
// 	for i, cwords := range c.words {
// 		if i >= len(cb.words) {
// 			length = i
// 			break
// 		}
// 		cb.words[i] = (cb.words[i] &^ cwords) | (cwords &^ cb.words[i])
// 	}
// 	if length != 0 {
// 		cb.words = append(cb.words, c.words[length:]...)
// 	}
// }
