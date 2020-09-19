package tool

import (
	"fmt"
	"sync"
)

// BitMap ... 0 means val not use,1 means val has be used
type BitMap interface {
	Push(val int64) error
	Remove(val int64) error
	IsValueUsed(val int64) (bool, error)
	GetAvailVal() (val int64, err error)
	AvailNum() int64
}

type bitMap struct {
	bm           []byte
	start        int64
	end          int64
	length       int64
	availableNum int64
	bmLen        int64
	index        int64 // 当前遍历的游标
	lock         sync.RWMutex
}

const (
	bitLength = 8
	bitNum    = 1<<8 - 1 //  0<=bit<=127
)

var (
	// ErrValExist ...
	ErrValExist = fmt.Errorf("val exist")
	// ErrNoAvailNum ...
	ErrNoAvailNum = fmt.Errorf("no avail num")
	// ErrBmNil ...
	ErrBmNil = fmt.Errorf("bitMap can't be nil")
)

// NewBitMap ... bm contains start and end
func NewBitMap(start, end int64) (BitMap, error) {
	if end < start {
		return nil, fmt.Errorf("bit map need end >= start")
	}
	length := end - start + 1
	bmLen := (length + 7) / bitLength
	b := &bitMap{
		bm:           make([]byte, bmLen),
		start:        start,
		end:          end,
		length:       length,
		availableNum: length,
		bmLen:        bmLen,
		index:        0,
		lock:         sync.RWMutex{},
	}
	return b, nil
}

func (b *bitMap) outOfRangeCheck(v int64) (bmIndex, bIndex int64, err error) {
	if v > b.end || v < b.start {
		return 0, 0, fmt.Errorf("val %d out of bitmap range", v)
	}
	bmIndex = (v - b.start) / bitLength
	bIndex = (v - b.start) % bitLength
	return
}

// 使用前先需要判断越界问题
// 本函数内部不在做越界判断，默认查询的都是合法
func (b *bitMap) isValUsed(bmIndex, bIndex int64) bool {
	b.lock.RLock()
	defer b.lock.RUnlock()

	bit := b.bm[bmIndex]
	return (bit & (1 << bIndex)) > 0
}

// means val used
func (b *bitMap) pushVal(bmIndex, bIndex int64) {
	b.lock.Lock()
	defer b.lock.Unlock()

	bit := b.bm[bmIndex]
	if (bit & (1 << bIndex)) > 0 {
		return
	}
	b.availableNum--
	b.bm[bmIndex] = bit | (1 << bIndex)
	return
}

// pop val,means release val
func (b *bitMap) removeVal(bmIndex, bIndex int64) {
	b.lock.Lock()
	defer b.lock.Unlock()

	bit := b.bm[bmIndex]
	if (bit & (1 << bIndex)) == 0 {
		return
	}
	b.availableNum++

	b.bm[bmIndex] = bit & (bitNum - 1<<bIndex)
	return
}

// Push ...
func (b *bitMap) Push(val int64) error {
	if b == nil || b.bm == nil {
		return ErrBmNil
	}
	bmIndex, bIndex, err := b.outOfRangeCheck(val)
	if err != nil {
		return err
	}
	if b.isValUsed(bmIndex, bIndex) {
		return ErrValExist
	}
	b.pushVal(bmIndex, bIndex)
	return nil
}

// Remove ...
func (b *bitMap) Remove(val int64) error {
	if b == nil || b.bm == nil {
		return ErrBmNil
	}
	bmIndex, bIndex, err := b.outOfRangeCheck(val)
	if err != nil {
		return err
	}
	b.removeVal(bmIndex, bIndex)
	return nil
}

// IsValueUsed ...
func (b *bitMap) IsValueUsed(val int64) (bool, error) {
	if b == nil || b.bm == nil {
		return false, ErrBmNil
	}
	bmIndex, bIndex, err := b.outOfRangeCheck(val)
	if err != nil {
		return false, err
	}

	return b.isValUsed(bmIndex, bIndex), nil
}

// GetAvailVal ... 从头遍历寻找一个未使用的数
func (b *bitMap) GetAvailVal() (val int64, err error) {
	if b == nil || b.bm == nil {
		err = ErrBmNil
		return
	}

	if b.availableNum <= 0 {
		err = ErrNoAvailNum
		return
	}

	find := false
	bmIndex := b.index / bitLength
	bIndex := b.index % bitLength
	for i := bmIndex; i < b.bmLen; i++ {
		for d := bIndex; d < bitLength; d++ {
			if b.isValUsed(i, d) {
				continue
			}

			val = i*bitLength + d + b.start
			b.index = i*bitLength + d + 1
			return val, nil
		}
		bIndex = 0
	}

	if !find {
		err = ErrNoAvailNum
	}

	return
}

func (b *bitMap) AvailNum() int64 {
	if b == nil || b.bm == nil {
		return 0
	}
	return b.availableNum
}
