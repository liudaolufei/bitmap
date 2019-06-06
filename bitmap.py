#!//usr/bin/env python
#-*- coding: utf-8 -*-

"""
    File Name: bitmap.py
    Author: luoheng
    Email: 1301089462@qq.com
    Version: 1.0
"""

from array import array


class BitMap:
 
    BITMASK = [0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80]
    BIT_CNT = [bin(i).count('1') for i in range(256)]

    def __init__(self, maxnum=1):
        self.nbytes = (maxnum + 7) // 8
        self.bitmap = array('B', [0 for i in range(self.nbytes)])

    def __or__(self, pos):
        return self.bitmap[pos // 8] | self.BITMASK[pos % 8]

    def __ior__(self, pos):
        self.bitmap[pos // 8] |= self.BITMASK[pos % 8]
        return self

    def __and__(self, pos):
        return self.bitmap[pos // 8] & ~self.BITMASK[pos % 8]

    def __iand__(self, pos):
        self.bitmap[pos // 8] &= ~self.BITMASK[pos % 8]
        return self

    def __xor__(self, pos):
        return 1 if self.bitmap[pos // 8] ^ self.BITMASK[pos % 8] else 0

    def __ixor__(self, pos):
        self.bitmap[pos // 8] ^= self.BITMASK[pos % 8]
        return self

    def __len__(self):
        return self.nbytes * 8

    def count(self):
        return sum(self.BIT_CNT[x] for x in self.bitmap)

    def __iter__(self):
        return self
    
    def next(self):
        for i in range(self.nbytes, -1, -1):
            for j in bin(self.bitmap[i])[2:][::-1]:
                yield int(j)


    def any(self):
        """
        Test if any bit is set
        """
        return self.count() > 0

    def none(self):
        """
        Test if no bit is set
        """
        return not self.any()

    def all(self):
        """
        Test if all bits are set
        """
        return (self.count() + 7) // 8 * 8 == len(self)

    def nonzero(self):
        """
        Get all non-zero bits
        """
        return (i for i in range(len(self)) if self & i != 0)

    def tostring(self):
        """
        Convert BitMap to string
        """
        return "".join([("%s" % bin(x)[2:]).zfill(8)
                        for x in self.bitmap[::-1]])

    def __str__(self):
        """
        Overloads string operator
        """
        return self.tostring()

    def __getitem__(self, item):
        """
        Return a bit when indexing like a array
        """
        return self & item

    def __setitem__(self, key, value):
        """
        Sets a bit when indexing like a array
        """
        if value is True:
            self |= key
        elif value is False:
            self &= key
        else:
            raise Exception("Use a boolean value to assign to a bitfield")

    def tohexstring(self):
        """
        Returns a hexadecimal string
        """
        val = self.tostring()
        st = "{0:0x}".format(int(val, 2))
        return st.zfill(len(self.bitmap)*2)

    @classmethod
    def fromhexstring(cls, hexstring):
        """
        Construct BitMap from hex string
        """
        bitstring = format(int(hexstring, 16), "0" + str(len(hexstring)//4) + "b")
        return cls.fromstring(bitstring)

    @classmethod
    def fromstring(cls, bitstring):
        """
        Construct BitMap from string
        """
        nbits = len(bitstring)
        bm = cls(nbits)
        for i in range(nbits):
            if bitstring[-i-1] == '1':
                bm |= i
            elif bitstring[-i-1] != '0':
                raise Exception("Invalid bit string!")
        return bm