# bitmap
bitmap for go. Bitmap is a suitable data structure to count integer in certain range, it uses one bit to represent one integer.
Here I implement a bitmap in golang, it's fast and space efficient.
Here is the bitmap interface:
```go
type Bitmap interface {
	Add(x int)     // add x to bitmap
	Has(x int)     // return true if x is in bitmap
	Remove(x int)  // remove x in bitmap 
	Len() int      // return length of bitmap
	Clear()        // clear bitmap to free memory
}
```
I implement four type of bitmap:
* `NBitmap`: normal bitmap, including set operation, use `New()` to get it.
* `RBitmap`: range bitmap, including set operation, use `NewR(start, end)` to get it.
* `CBitmap`: bitmap that can count elements, use `NewC(n)` to get it.
* `RCBitmap`: range bitmap that can count elements, use `NewRC(start, end, n)` to get it.
# NBitmap
NBitmap is normal bitmap, including set operation.
```go
// get a new bitmap
b := bitmap.New()
// add elements
b.Add(10)
b.Add(100)
// String
b.String() // {10 100}
// Length
b.Len() // 2
// check if has the elements
if b.Has(100) {/* code */}
// Remove elements
b.Remove(10)
// Clear bitmap
// do this to manually free memory
b.Clear()
// operation for sets
c := bitmap.New()
// Union b |= c
// elements in b or c
b.Union(c)
// Intersect b &= c
// elements both in b and c
b.Intersect(c)
// Except b -= c
// elements only in b
b.Except(c)
// SymExcept b = (b - c) | (c - b)
// elements only in b or only in c
b.SymExcept(c)
```
# RBitmap
RBitmap is a range bitmap, it's similar to NBitmap, all elements should be in the range.
##notice##: set operations can only work on two bitmap with same range.
```go
b := bitmap.NewR(0, 5) // this bitmap is used to count [0, 4]
```
# CBitmap
CBitmap is a bitmap that can count elements.
```go
// get a new bitmap of CBitmap, it can count to 5
b := bitmap.NewC(5)
// add elements
b.Add(10)
b.Add(10)
b.Add(10)
b.Add(100)
// Count
b.Count(10)  // 3
// String
b.String() // {10 100}
// Length
b.Len() // 2
// check if has the elements
if b.Has(100) {/* code */}
// Remove elements
b.Remove(10)
b.RemoveAll(10)
// Clear bitmap
// do this to manually free memory
b.Clear()
```
# RCBitmap
RCBitmap is a range CBitmap
```go
b := bitmap.NewRC(0, 5, 4) // count [0, 4], the max count number is 4.
```