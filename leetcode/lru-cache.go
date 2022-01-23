package main

func main()  {
	lRUCache := LRUCache{capacity:2};
	lRUCache.Put(1, 1); // 缓存是 {1=1}
	lRUCache.Put(2, 2); // 缓存是 {1=1, 2=2}
	lRUCache.Get(1);    // 返回 1
	lRUCache.Put(3, 3); // 该操作会使得关键字 2 作废，缓存是 {1=1, 3=3}
	lRUCache.Get(2);    // 返回 -1 (未找到)
	lRUCache.Put(4, 4); // 该操作会使得关键字 1 作废，缓存是 {4=4, 3=3}
	lRUCache.Get(1);    // 返回 -1 (未找到)
	lRUCache.Get(3);    // 返回 3
	lRUCache.Get(4);    // 返回 4
}

type Node struct {
	key int
	val int
	pre *Node
	next *Node
}

type LRUCache struct {
	capacity int
	l int
	root *Node
	last *Node
	cache bool
	cacheMap map[int]*Node
}

func Constructor(capacity int) LRUCache {
	return LRUCache{capacity:capacity, cache:true, cacheMap:make(map[int]*Node)}
}


func (this *LRUCache) Get(key int) int {
	p := this.search(key)
	if p != nil {
		println(p.val)
		return p.val
	}

	println(-1)
	return -1
}

func (this *LRUCache) moveToHead(p *Node) {
	if p != this.root {
		pre := p.pre
		next := p.next

		pre.next = next

		if next == nil {
			this.last = pre
		} else {
			next.pre = pre
		}

		p.next = this.root
		this.root.pre = p
		this.root = p
	}
}

func (this *LRUCache) search(key int) *Node {
	if this.cache {
		v, ok := this.cacheMap[key]
		if !ok {
			return nil
		} else {
			this.moveToHead(v)
			return v
		}
	}

	p := this.root
	for p != nil {
		if p.key == key {
			this.moveToHead(p)
			return p
		}
		p = p.next
	}

	return nil
}

func (this *LRUCache) Put(key int, value int)  {
	p := this.search(key)
	if p != nil {
		p.val = value
		return
	}

	n := &Node{
		key:key,
		val:value,
		next:this.root,
	}

	if this.root != nil {
		this.root.pre = n
	} else {
		this.last = n
	}
	this.root = n
	this.l++
	if this.cache {
		this.cacheMap[key] = n
	}

	if this.l > this.capacity {
		if this.cache {
			delete(this.cacheMap, this.last.key)
		}
		this.last.pre.next = nil
		this.last = this.last.pre
		this.l--
	}
}


/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */
