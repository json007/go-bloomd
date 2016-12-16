package bloomfilter

import (
	"hash"
	"hash/fnv"
	"math"
)

type BitSet struct {
	sets []uint64 //本想用[]byte类型，但是位运算不方便，于是用
}

func (bs *BitSet) set(position uint32){
	offset := uint64(1) >> position%8
	key := position/8
	bs.sets[key] = bs.sets[key] | offset
	// bs.sets[position/8] = bs.sets[position/8] | offset
}

func (bs *BitSet) check(position uint32) bool{
	offset := uint64(1) >> position%8
	return (bs.sets[position/8] & offset) == offset
}

type BloomFilter struct {
	bitset *BitSet
	hashfn hash.Hash64
	//有 hash.Hash  hash.Hash32 hash.Hash64 。
	//应该选 hash.Hash，int是最小32bits，相当于要用 2^(32 - log_2 8) = 2^29个byte ~~ (2^10)^3  byte->KB->MB->GB 。
	//差不多一个G了 ，已经足够大了。
	num int//已存储多少个
	size int
	fnNum int
}

func (bf *BloomFilter) getHash(b []byte) [2]uint32 {
	bf.hashfn.Reset()
	bf.hashfn.Write(b)
	hash64 := bf.hashfn.Sum64()
	h1 := uint32(hash64 & ((1 << 32) - 1))
	h2 := uint32(hash64 >> 32)
	return [2]uint32{h1, h2}
}

func (bf *BloomFilter) Add(b []byte){
	// for i,f := range bf.hashfun {
	// 	position = bf.hash(b []byte)
	// 	bf.bitset.set(f(position))
	// }
	positions := bf.getHash(b)
	for _,position := range positions {
		bf.bitset.set(position)
	}
	bf.num++  //得先看看是否已经添加过了
}

func (bf *BloomFilter) IsContain(b []byte) bool {
	// for i,f := range bf.hashfun {
	// 	position = bf.hash(b []byte)
	// 	bf.bitset.set(f(position))
	// }
	positions := bf.getHash(b)
	for _,position := range positions {
		if !bf.bitset.check(position) {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) FalsePositiveRate() float64 {
	return math.Pow((1 - math.Exp(-float64(bf.fnNum*bf.num)/float64(bf.size))), float64(bf.fnNum))
}

func New(size int) *BloomFilter {
	bf := new(BloomFilter)
	// bf.bit = make([]int, size)
	bf.num = 0
	bf.bitset = new(BitSet)
	bf.bitset.sets = make([]uint64,size)
	bf.hashfn = fnv.New64()
	return bf
}
