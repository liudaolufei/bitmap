# bitmap
bitmap for go

# Usage
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
if b.Has(100) {
    // code
}
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