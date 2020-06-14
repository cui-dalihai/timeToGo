package bloom


import (
	"github.com/willf/bitset"
	"github.com/spaolacci/murmur3"
)


func max(x, y uint) uint {
	if x > y {
		return x
	}
	return y
}

type BloomFilter struct {
	m uint   // 64位
	k uint
	b *bitset.BitSet
}

func New(m uint, k uint) *BloomFilter {
	return &BloomFilter{max(1,m), max(1,k),bitset.New(m)}
}

func From(data []uint64, k uint) *BloomFilter {
	m := uint(len(data) * 64)
	return &BloomFilter{m, k, bitset.From(data)}
}

func baseHashes(data []byte) [4]uint64 {
	a1 := []byte{1}
	hasher := murmur3.New128()
	hasher.Write(data)  // #nosec
	v1, v2 := hasher.Sum128()
	hasher.Write(a1)
	v3, v4 := hasher.Sum128()
	return [4]uint64 {v1, v2, v3, v4}
}

func location(h [4]uint64, i uint) uint64 {
	ii := uint64(i)
	return h[ii%2] + ii*h[2+(((ii+(ii%2))%4)/2)]
}







