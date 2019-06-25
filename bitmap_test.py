
import unittest
from bitmap import BitMap

class BitMapTest(unittest.TestCase):
    def testOperator(self):
        bp = BitMap(128)
        print(bp.count())
        print(len(bp))
        bp |= 12
        bp |= 13
        bp |= 14
        print(bp | 12)
        print(bp | 13)
        print(bp | 14)
        print(bp | 15)
        print(bp.count())
        bp &= 13
        print(bp)
        print(bp.count())
        bp ^= 14
        print(bp, bp ^ 14)
        print(bp)
        print(bp.count())

    def testMethod(self):
        pass

if __name__ == "__main__":
    unittest.main()