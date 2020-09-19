package tool

import (
	"fmt"
	"sync"
	"testing"
)

func TestNewBitMap(t *testing.T) {
	{
		start := int64(1)
		end := int64(2)
		NewBitMap(start, end)
	}

	{
		start := int64(2)
		end := int64(1)
		NewBitMap(start, end)
	}

	{
		start := int64(1)
		end := int64(300)
		NewBitMap(start, end)
	}
}

func TestPushVal(t *testing.T) {
	b := &bitMap{
		bm:           make([]byte, 5),
		start:        1,
		end:          38,
		length:       39,
		availableNum: 39,
		bmLen:        5,
		lock:         sync.RWMutex{},
	}
	{
		b.pushVal(2, 1)
		b.pushVal(2, 1)
		b.pushVal(2, 0)
		b.pushVal(2, 7)
	}
}

func TestRemoveVal(t *testing.T) {
	b := &bitMap{
		bm:           make([]byte, 5),
		start:        1,
		end:          38,
		length:       39,
		availableNum: 39,
		bmLen:        5,
		lock:         sync.RWMutex{},
	}
	b.pushVal(2, 1)
	b.pushVal(2, 0)
	b.pushVal(2, 7)

	{
		b.removeVal(2, 0)
		b.removeVal(2, 0)
	}
}

func TestIsValUsed(t *testing.T) {
	b := &bitMap{
		bm:           make([]byte, 5),
		start:        1,
		end:          38,
		length:       39,
		availableNum: 39,
		bmLen:        5,
		lock:         sync.RWMutex{},
	}
	b.pushVal(2, 1)
	b.pushVal(2, 0)
	b.pushVal(2, 7)

	{
		isUsed := b.isValUsed(2, 1)
		fmt.Printf("isUsed:%t\n", isUsed)
		isUsed = b.isValUsed(2, 2)
		fmt.Printf("isUsed:%t\n", isUsed)
	}
}

func TestGetAvailVal(t *testing.T) {
	b := &bitMap{
		bm:           make([]byte, 5),
		start:        1,
		end:          38,
		length:       39,
		availableNum: 39,
		bmLen:        5,
		lock:         sync.RWMutex{},
	}
	b.pushVal(2, 1)
	b.pushVal(2, 0)
	b.pushVal(0, 1)
	b.pushVal(2, 7)

	{
		avail, err := b.GetAvailVal()
		if err != nil {
			fmt.Printf("err:%s\n", err.Error())
		}
		fmt.Printf("availNum:%d\n", avail)
		avail, err = b.GetAvailVal()
		if err != nil {
			fmt.Printf("err:%s\n", err.Error())
		}
		fmt.Printf("availNum:%d\n", avail)
	}

	{
		be := &bitMap{}
		avail, err := be.GetAvailVal()
		if err != nil {
			fmt.Printf("err:%s\n", err.Error())
		}
		fmt.Printf("availNum:%d\n", avail)
	}

	{
		b.availableNum = 0
		avail, err := b.GetAvailVal()
		if err != nil {
			fmt.Printf("err:%s\n", err.Error())
		}
		fmt.Printf("availNum:%d\n", avail)
	}
}

func TestAvailNum(t *testing.T) {
	bm, _ := NewBitMap(1, 200)

	availNum := bm.AvailNum()
	fmt.Printf("availNum:%d\n", availNum)
}

func TestPush(t *testing.T) {
	bm, _ := NewBitMap(1, 200)

	err := bm.Push(20)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}

	err = bm.Push(200)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}

	err = bm.Push(200)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}

	err = bm.Push(201)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}

	{
		bm.(*bitMap).bm = nil
		err := bm.Push(21)
		if err != nil {
			fmt.Printf("err:%s\n", err.Error())
		}
	}

}

func TestRemove(t *testing.T) {
	bm, _ := NewBitMap(1, 200)
	bm.Push(200)
	bm.Push(20)
	bm.Push(12)
	bm.Push(11)

	err := bm.Remove(20)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}

	err = bm.Remove(200)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}

	err = bm.Remove(200)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}

	err = bm.Remove(201)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}

	{
		bm.(*bitMap).bm = nil
		err := bm.Remove(21)
		if err != nil {
			fmt.Printf("err:%s\n", err.Error())
		}
	}

}

func TestIsValueUsed(t *testing.T) {
	bm, _ := NewBitMap(1, 200)
	bm.Push(200)
	bm.Push(20)
	bm.Push(12)
	bm.Push(11)

	isUsed, err := bm.IsValueUsed(20)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	fmt.Printf("isUsed:%t\n", isUsed)

	isUsed, err = bm.IsValueUsed(200)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	fmt.Printf("isUsed:%t\n", isUsed)

	isUsed, err = bm.IsValueUsed(202)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	fmt.Printf("isUsed:%t\n", isUsed)

	isUsed, err = bm.IsValueUsed(3)
	if err != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	fmt.Printf("isUsed:%t\n", isUsed)
}
