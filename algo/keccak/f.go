package keccak

// this file must not import any dependencies

var randomizer = [24]uint64{
	0x0000000000000001,
	0x0000000000008082,
	0x800000000000808A,
	0x8000000080008000,
	0x000000000000808B,
	0x0000000080000001,
	0x8000000080008081,
	0x8000000000008009,
	0x000000000000008A,
	0x0000000000000088,
	0x0000000080008009,
	0x000000008000000A,
	0x000000008000808B,
	0x800000000000008B,
	0x8000000000008089,
	0x8000000000008003,
	0x8000000000008002,
	0x8000000000000080,
	0x000000000000800A,
	0x800000008000000A,
	0x8000000080008081,
	0x8000000000008080,
	0x0000000080000001,
	0x8000000080008008,
}

// this function is optimized for performance and only partially for readability
func permute(a *[25]uint64) {
	var x, y [5]uint64
	var m uint64

	for j := 0; j < 6; j++ {
		r := j << 2

		x[0] = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]
		x[1] = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]
		x[2] = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]
		x[3] = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]
		x[4] = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]

		y[0] = x[4] ^ (x[1]<<1 | x[1]>>63)
		y[1] = x[0] ^ (x[2]<<1 | x[2]>>63)
		y[2] = x[1] ^ (x[3]<<1 | x[3]>>63)
		y[3] = x[2] ^ (x[4]<<1 | x[4]>>63)
		y[4] = x[3] ^ (x[0]<<1 | x[0]>>63)

		x[0] = a[0] ^ y[0]
		m = a[6] ^ y[1]
		x[1] = m<<44 | m>>(64-44)
		m = a[12] ^ y[2]
		x[2] = m<<43 | m>>(64-43)
		m = a[18] ^ y[3]
		x[3] = m<<21 | m>>(64-21)
		m = a[24] ^ y[4]
		x[4] = m<<14 | m>>(64-14)

		a[0] = x[0] ^ (x[2] &^ x[1])
		a[6] = x[1] ^ (x[3] &^ x[2])
		a[12] = x[2] ^ (x[4] &^ x[3])
		a[18] = x[3] ^ (x[0] &^ x[4])
		a[24] = x[4] ^ (x[1] &^ x[0])

		m = a[10] ^ y[0]
		x[2] = m<<3 | m>>(64-3)
		m = a[16] ^ y[1]
		x[3] = m<<45 | m>>(64-45)
		m = a[22] ^ y[2]
		x[4] = m<<61 | m>>(64-61)
		m = a[3] ^ y[3]
		x[0] = m<<28 | m>>(64-28)
		m = a[9] ^ y[4]
		x[1] = m<<20 | m>>(64-20)

		a[10] = x[0] ^ (x[2] &^ x[1])
		a[16] = x[1] ^ (x[3] &^ x[2])
		a[22] = x[2] ^ (x[4] &^ x[3])
		a[3] = x[3] ^ (x[0] &^ x[4])
		a[9] = x[4] ^ (x[1] &^ x[0])

		m = a[20] ^ y[0]
		x[4] = m<<18 | m>>(64-18)
		m = a[1] ^ y[1]
		x[0] = m<<1 | m>>(64-1)
		m = a[7] ^ y[2]
		x[1] = m<<6 | m>>(64-6)
		m = a[13] ^ y[3]
		x[2] = m<<25 | m>>(64-25)
		m = a[19] ^ y[4]
		x[3] = m<<8 | m>>(64-8)

		a[20] = x[0] ^ (x[2] &^ x[1])
		a[1] = x[1] ^ (x[3] &^ x[2])
		a[7] = x[2] ^ (x[4] &^ x[3])
		a[13] = x[3] ^ (x[0] &^ x[4])
		a[19] = x[4] ^ (x[1] &^ x[0])

		m = a[5] ^ y[0]
		x[1] = m<<36 | m>>(64-36)
		m = a[11] ^ y[1]
		x[2] = m<<10 | m>>(64-10)
		m = a[17] ^ y[2]
		x[3] = m<<15 | m>>(64-15)
		m = a[23] ^ y[3]
		x[4] = m<<56 | m>>(64-56)
		m = a[4] ^ y[4]
		x[0] = m<<27 | m>>(64-27)

		a[5] = x[0] ^ (x[2] &^ x[1])
		a[11] = x[1] ^ (x[3] &^ x[2])
		a[17] = x[2] ^ (x[4] &^ x[3])
		a[23] = x[3] ^ (x[0] &^ x[4])
		a[4] = x[4] ^ (x[1] &^ x[0])

		m = a[15] ^ y[0]
		x[3] = m<<41 | m>>(64-41)
		m = a[21] ^ y[1]
		x[4] = m<<2 | m>>(64-2)
		m = a[2] ^ y[2]
		x[0] = m<<62 | m>>(64-62)
		m = a[8] ^ y[3]
		x[1] = m<<55 | m>>(64-55)
		m = a[14] ^ y[4]
		x[2] = m<<39 | m>>(64-39)

		a[15] = x[0] ^ (x[2] &^ x[1])
		a[21] = x[1] ^ (x[3] &^ x[2])
		a[2] = x[2] ^ (x[4] &^ x[3])
		a[8] = x[3] ^ (x[0] &^ x[4])
		a[14] = x[4] ^ (x[1] &^ x[0])

		a[0] ^= randomizer[r]

		///////////////////////////////////////////////////////////////////////////

		x[0] = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]
		x[1] = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]
		x[2] = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]
		x[3] = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]
		x[4] = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]

		y[0] = x[4] ^ (x[1]<<1 | x[1]>>63)
		y[1] = x[0] ^ (x[2]<<1 | x[2]>>63)
		y[2] = x[1] ^ (x[3]<<1 | x[3]>>63)
		y[3] = x[2] ^ (x[4]<<1 | x[4]>>63)
		y[4] = x[3] ^ (x[0]<<1 | x[0]>>63)

		x[0] = a[0] ^ y[0]
		m = a[16] ^ y[1]
		x[1] = m<<44 | m>>(64-44)
		m = a[7] ^ y[2]
		x[2] = m<<43 | m>>(64-43)
		m = a[23] ^ y[3]
		x[3] = m<<21 | m>>(64-21)
		m = a[14] ^ y[4]
		x[4] = m<<14 | m>>(64-14)

		a[0] = x[0] ^ (x[2] &^ x[1])
		a[16] = x[1] ^ (x[3] &^ x[2])
		a[7] = x[2] ^ (x[4] &^ x[3])
		a[23] = x[3] ^ (x[0] &^ x[4])
		a[14] = x[4] ^ (x[1] &^ x[0])

		m = a[20] ^ y[0]
		x[2] = m<<3 | m>>(64-3)
		m = a[11] ^ y[1]
		x[3] = m<<45 | m>>(64-45)
		m = a[2] ^ y[2]
		x[4] = m<<61 | m>>(64-61)
		m = a[18] ^ y[3]
		x[0] = m<<28 | m>>(64-28)
		m = a[9] ^ y[4]
		x[1] = m<<20 | m>>(64-20)

		a[20] = x[0] ^ (x[2] &^ x[1])
		a[11] = x[1] ^ (x[3] &^ x[2])
		a[2] = x[2] ^ (x[4] &^ x[3])
		a[18] = x[3] ^ (x[0] &^ x[4])
		a[9] = x[4] ^ (x[1] &^ x[0])

		m = a[15] ^ y[0]
		x[4] = m<<18 | m>>(64-18)
		m = a[6] ^ y[1]
		x[0] = m<<1 | m>>(64-1)
		m = a[22] ^ y[2]
		x[1] = m<<6 | m>>(64-6)
		m = a[13] ^ y[3]
		x[2] = m<<25 | m>>(64-25)
		m = a[4] ^ y[4]
		x[3] = m<<8 | m>>(64-8)

		a[15] = x[0] ^ (x[2] &^ x[1])
		a[6] = x[1] ^ (x[3] &^ x[2])
		a[22] = x[2] ^ (x[4] &^ x[3])
		a[13] = x[3] ^ (x[0] &^ x[4])
		a[4] = x[4] ^ (x[1] &^ x[0])

		m = a[10] ^ y[0]
		x[1] = m<<36 | m>>(64-36)
		m = a[1] ^ y[1]
		x[2] = m<<10 | m>>(64-10)
		m = a[17] ^ y[2]
		x[3] = m<<15 | m>>(64-15)
		m = a[8] ^ y[3]
		x[4] = m<<56 | m>>(64-56)
		m = a[24] ^ y[4]
		x[0] = m<<27 | m>>(64-27)

		a[10] = x[0] ^ (x[2] &^ x[1])
		a[1] = x[1] ^ (x[3] &^ x[2])
		a[17] = x[2] ^ (x[4] &^ x[3])
		a[8] = x[3] ^ (x[0] &^ x[4])
		a[24] = x[4] ^ (x[1] &^ x[0])

		m = a[5] ^ y[0]
		x[3] = m<<41 | m>>(64-41)
		m = a[21] ^ y[1]
		x[4] = m<<2 | m>>(64-2)
		m = a[12] ^ y[2]
		x[0] = m<<62 | m>>(64-62)
		m = a[3] ^ y[3]
		x[1] = m<<55 | m>>(64-55)
		m = a[19] ^ y[4]
		x[2] = m<<39 | m>>(64-39)

		a[5] = x[0] ^ (x[2] &^ x[1])
		a[21] = x[1] ^ (x[3] &^ x[2])
		a[12] = x[2] ^ (x[4] &^ x[3])
		a[3] = x[3] ^ (x[0] &^ x[4])
		a[19] = x[4] ^ (x[1] &^ x[0])

		a[0] ^= randomizer[r+1]

		///////////////////////////////////////////////////////////////////////////

		x[0] = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]
		x[1] = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]
		x[2] = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]
		x[3] = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]
		x[4] = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]

		y[0] = x[4] ^ (x[1]<<1 | x[1]>>63)
		y[1] = x[0] ^ (x[2]<<1 | x[2]>>63)
		y[2] = x[1] ^ (x[3]<<1 | x[3]>>63)
		y[3] = x[2] ^ (x[4]<<1 | x[4]>>63)
		y[4] = x[3] ^ (x[0]<<1 | x[0]>>63)

		x[0] = a[0] ^ y[0]
		m = a[11] ^ y[1]
		x[1] = m<<44 | m>>(64-44)
		m = a[22] ^ y[2]
		x[2] = m<<43 | m>>(64-43)
		m = a[8] ^ y[3]
		x[3] = m<<21 | m>>(64-21)
		m = a[19] ^ y[4]
		x[4] = m<<14 | m>>(64-14)

		a[0] = x[0] ^ (x[2] &^ x[1])
		a[11] = x[1] ^ (x[3] &^ x[2])
		a[22] = x[2] ^ (x[4] &^ x[3])
		a[8] = x[3] ^ (x[0] &^ x[4])
		a[19] = x[4] ^ (x[1] &^ x[0])

		m = a[15] ^ y[0]
		x[2] = m<<3 | m>>(64-3)
		m = a[1] ^ y[1]
		x[3] = m<<45 | m>>(64-45)
		m = a[12] ^ y[2]
		x[4] = m<<61 | m>>(64-61)
		m = a[23] ^ y[3]
		x[0] = m<<28 | m>>(64-28)
		m = a[9] ^ y[4]
		x[1] = m<<20 | m>>(64-20)

		a[15] = x[0] ^ (x[2] &^ x[1])
		a[1] = x[1] ^ (x[3] &^ x[2])
		a[12] = x[2] ^ (x[4] &^ x[3])
		a[23] = x[3] ^ (x[0] &^ x[4])
		a[9] = x[4] ^ (x[1] &^ x[0])

		m = a[5] ^ y[0]
		x[4] = m<<18 | m>>(64-18)
		m = a[16] ^ y[1]
		x[0] = m<<1 | m>>(64-1)
		m = a[2] ^ y[2]
		x[1] = m<<6 | m>>(64-6)
		m = a[13] ^ y[3]
		x[2] = m<<25 | m>>(64-25)
		m = a[24] ^ y[4]
		x[3] = m<<8 | m>>(64-8)

		a[5] = x[0] ^ (x[2] &^ x[1])
		a[16] = x[1] ^ (x[3] &^ x[2])
		a[2] = x[2] ^ (x[4] &^ x[3])
		a[13] = x[3] ^ (x[0] &^ x[4])
		a[24] = x[4] ^ (x[1] &^ x[0])

		m = a[20] ^ y[0]
		x[1] = m<<36 | m>>(64-36)
		m = a[6] ^ y[1]
		x[2] = m<<10 | m>>(64-10)
		m = a[17] ^ y[2]
		x[3] = m<<15 | m>>(64-15)
		m = a[3] ^ y[3]
		x[4] = m<<56 | m>>(64-56)
		m = a[14] ^ y[4]
		x[0] = m<<27 | m>>(64-27)

		a[20] = x[0] ^ (x[2] &^ x[1])
		a[6] = x[1] ^ (x[3] &^ x[2])
		a[17] = x[2] ^ (x[4] &^ x[3])
		a[3] = x[3] ^ (x[0] &^ x[4])
		a[14] = x[4] ^ (x[1] &^ x[0])

		m = a[10] ^ y[0]
		x[3] = m<<41 | m>>(64-41)
		m = a[21] ^ y[1]
		x[4] = m<<2 | m>>(64-2)
		m = a[7] ^ y[2]
		x[0] = m<<62 | m>>(64-62)
		m = a[18] ^ y[3]
		x[1] = m<<55 | m>>(64-55)
		m = a[4] ^ y[4]
		x[2] = m<<39 | m>>(64-39)

		a[10] = x[0] ^ (x[2] &^ x[1])
		a[21] = x[1] ^ (x[3] &^ x[2])
		a[7] = x[2] ^ (x[4] &^ x[3])
		a[18] = x[3] ^ (x[0] &^ x[4])
		a[4] = x[4] ^ (x[1] &^ x[0])

		a[0] ^= randomizer[r+2]

		///////////////////////////////////////////////////////////////////////////

		x[0] = a[0] ^ a[5] ^ a[10] ^ a[15] ^ a[20]
		x[1] = a[1] ^ a[6] ^ a[11] ^ a[16] ^ a[21]
		x[2] = a[2] ^ a[7] ^ a[12] ^ a[17] ^ a[22]
		x[3] = a[3] ^ a[8] ^ a[13] ^ a[18] ^ a[23]
		x[4] = a[4] ^ a[9] ^ a[14] ^ a[19] ^ a[24]

		y[0] = x[4] ^ (x[1]<<1 | x[1]>>63)
		y[1] = x[0] ^ (x[2]<<1 | x[2]>>63)
		y[2] = x[1] ^ (x[3]<<1 | x[3]>>63)
		y[3] = x[2] ^ (x[4]<<1 | x[4]>>63)
		y[4] = x[3] ^ (x[0]<<1 | x[0]>>63)

		x[0] = a[0] ^ y[0]
		m = a[1] ^ y[1]
		x[1] = m<<44 | m>>(64-44)
		m = a[2] ^ y[2]
		x[2] = m<<43 | m>>(64-43)
		m = a[3] ^ y[3]
		x[3] = m<<21 | m>>(64-21)
		m = a[4] ^ y[4]
		x[4] = m<<14 | m>>(64-14)

		a[0] = x[0] ^ (x[2] &^ x[1])
		a[1] = x[1] ^ (x[3] &^ x[2])
		a[2] = x[2] ^ (x[4] &^ x[3])
		a[3] = x[3] ^ (x[0] &^ x[4])
		a[4] = x[4] ^ (x[1] &^ x[0])

		m = a[5] ^ y[0]
		x[2] = m<<3 | m>>(64-3)
		m = a[6] ^ y[1]
		x[3] = m<<45 | m>>(64-45)
		m = a[7] ^ y[2]
		x[4] = m<<61 | m>>(64-61)
		m = a[8] ^ y[3]
		x[0] = m<<28 | m>>(64-28)
		m = a[9] ^ y[4]
		x[1] = m<<20 | m>>(64-20)

		a[5] = x[0] ^ (x[2] &^ x[1])
		a[6] = x[1] ^ (x[3] &^ x[2])
		a[7] = x[2] ^ (x[4] &^ x[3])
		a[8] = x[3] ^ (x[0] &^ x[4])
		a[9] = x[4] ^ (x[1] &^ x[0])

		m = a[10] ^ y[0]
		x[4] = m<<18 | m>>(64-18)
		m = a[11] ^ y[1]
		x[0] = m<<1 | m>>(64-1)
		m = a[12] ^ y[2]
		x[1] = m<<6 | m>>(64-6)
		m = a[13] ^ y[3]
		x[2] = m<<25 | m>>(64-25)
		m = a[14] ^ y[4]
		x[3] = m<<8 | m>>(64-8)

		a[10] = x[0] ^ (x[2] &^ x[1])
		a[11] = x[1] ^ (x[3] &^ x[2])
		a[12] = x[2] ^ (x[4] &^ x[3])
		a[13] = x[3] ^ (x[0] &^ x[4])
		a[14] = x[4] ^ (x[1] &^ x[0])

		m = a[15] ^ y[0]
		x[1] = m<<36 | m>>(64-36)
		m = a[16] ^ y[1]
		x[2] = m<<10 | m>>(64-10)
		m = a[17] ^ y[2]
		x[3] = m<<15 | m>>(64-15)
		m = a[18] ^ y[3]
		x[4] = m<<56 | m>>(64-56)
		m = a[19] ^ y[4]
		x[0] = m<<27 | m>>(64-27)

		a[15] = x[0] ^ (x[2] &^ x[1])
		a[16] = x[1] ^ (x[3] &^ x[2])
		a[17] = x[2] ^ (x[4] &^ x[3])
		a[18] = x[3] ^ (x[0] &^ x[4])
		a[19] = x[4] ^ (x[1] &^ x[0])

		m = a[20] ^ y[0]
		x[3] = m<<41 | m>>(64-41)
		m = a[21] ^ y[1]
		x[4] = m<<2 | m>>(64-2)
		m = a[22] ^ y[2]
		x[0] = m<<62 | m>>(64-62)
		m = a[23] ^ y[3]
		x[1] = m<<55 | m>>(64-55)
		m = a[24] ^ y[4]
		x[2] = m<<39 | m>>(64-39)

		a[20] = x[0] ^ (x[2] &^ x[1])
		a[21] = x[1] ^ (x[3] &^ x[2])
		a[22] = x[2] ^ (x[4] &^ x[3])
		a[23] = x[3] ^ (x[0] &^ x[4])
		a[24] = x[4] ^ (x[1] &^ x[0])

		a[0] ^= randomizer[r+3]
	}
}
