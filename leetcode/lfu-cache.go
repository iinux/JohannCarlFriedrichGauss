package main

import "fmt"

// 怎么查找
// 怎么升级
// 怎么插入
// 怎么移除

func main() {
	/*
	testLFU(
		[]string{"LFUCache", "put", "put", "get", "put", "get", "get", "put", "get", "get", "get"},
		[][]int{{2}, {1, 1}, {2, 2}, {1}, {3, 3}, {2}, {3}, {4, 4}, {1}, {3}, {4}},
		[]interface{}{nil, nil, nil, 1, nil, -1, 3, nil, -1, 3, 4},
	)

	 */
	testLFU(
		[]string{"LFUCache","put","put","put","put","get","get","get","get","put","get","get","get","get","get"},
		[][]int{{3},{1,1},{2,2},{3,3},{4,4},{4},{3},{2},{1},{5,5},{1},{2},{3},{4},{5}},
		[]interface{}{nil,nil,nil,nil,nil,4,3,2,-1,nil,-1,2,3,-1,5},
		)
}

func testLFU(cmdList []string, argList [][]int, answer []interface{}) {
	var lfuCache LFUCache
	for i, cmd := range cmdList {
		switch cmd {
		case "LFUCache":
			lfuCache = Constructor(argList[i][0])
			fmt.Println("     cap", argList[i][0])
		case "put":
			lfuCache.Put(argList[i][0], argList[i][1])
			fmt.Println("     put", argList[i])
		case "get":
			result := lfuCache.Get(argList[i][0])

			if result != answer[i] {
				fmt.Println("fail", cmdList[i], argList[i], answer[i], result)
				return
			} else {
				fmt.Println("pass", cmdList[i], argList[i], answer[i], result)
			}
		}
	}

	fmt.Println("success")
}

type Node struct {
	key   int
	value int
	count int

	pre  *Node
	next *Node
}

type Entry struct {
	head *Node
	tail *Node
}

type LFUCache struct {
	entryMap map[int]*Entry
	keyMap   map[int]*Node
	cap      int
	len      int
	maxCount int
}

func Constructor(capacity int) LFUCache {
	return LFUCache{
		entryMap: make(map[int]*Entry),
		keyMap:   make(map[int]*Node),
		cap:      capacity,
		maxCount:1,
	}
}

func (this *LFUCache) getEntry(count int) *Entry {
	entry, ok := this.entryMap[count]
	if !ok {
		entry = &Entry{
			head: nil,
			tail: nil,
		}
		this.entryMap[count] = entry
	}
	return entry
}

func (this *LFUCache) addToEntry(entry *Entry, p *Node) {
	if entry.head == nil {
		entry.head = p
		entry.tail = p
	} else {
		entry.head.pre = p
		p.next = entry.head

		entry.head = p
	}
}

func (this *LFUCache) removeFromEntry(entry *Entry, p *Node) {
	if entry.head == p && entry.tail == p {
		entry.head = nil
		entry.tail = nil
	} else if entry.head == p {
		entry.head = p.next
	} else if entry.tail == p {
		entry.tail = p.pre
	} else {
		pre := p.pre
		next := p.next
		pre.next = next
		next.pre = pre
	}
}

func (this *LFUCache) removeEntryLast(entry *Entry, p *Node) bool {
	if entry.tail == nil {
		return false
	}

	if p.key == entry.tail.key {
		return false
	}

	delete(this.keyMap, entry.tail.key)

	if entry.head == entry.tail {
		entry.head = nil
		entry.tail = nil
	} else {
		entry.tail = entry.tail.pre
	}
	return true
}

func (this *LFUCache) incCount(p *Node) {
	entry := this.getEntry(p.count)
	nextEntry := this.getEntry(p.count + 1)

	this.removeFromEntry(entry, p)
	this.addToEntry(nextEntry, p)

	p.count++
	if p.count > this.maxCount {
		this.maxCount = p.count
	}
}

func (this *LFUCache) Get(key int) int {
	v, ok := this.keyMap[key]
	if ok {
		this.incCount(v)
		return v.value
	}

	return -1
}

func (this *LFUCache) Put(key int, value int) {
	if this.cap == 0 {
		return
	}
	v, ok := this.keyMap[key]
	if ok {
		v.value = value
		this.incCount(v)
		return
	}

	p := &Node{
		key:   key,
		value: value,
		count: 1,
	}

	entry := this.getEntry(1)
	this.addToEntry(entry, p)
	this.keyMap[p.key] = p
	this.len++

	if this.len > this.cap {
		for i := 1; i <= this.maxCount; i++ {
			entry = this.getEntry(i)
			result := this.removeEntryLast(entry, p)
			if result {
				break
			}
		}

		this.len--
	}
}

/**
 * Your LFUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */
