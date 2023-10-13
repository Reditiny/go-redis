package database

import "math/rand"

const (
	maxLevel = 16
)

type sortedSet struct {
	dict     map[string]*element
	skipList *skipList
}

type skipList struct {
	header *node
	tail   *node
	length int
	level  int
}

// score-key pair
type element struct {
	score float64
	key   string
}

// Level next node and span to next node in this level
type Level struct {
	forward *node
	span    int
}

// node in skip list
type node struct {
	element
	backward *node
	level    []*Level
}

type Consumer func(node *node) bool

// makeNode makes a new node
func makeNode(level int, score float64, key string) *node {
	n := &node{
		element: element{
			score: score,
			key:   key,
		},
		level: make([]*Level, level),
	}
	for i := range n.level {
		n.level[i] = new(Level)
	}
	return n
}

// makeSortedSet makes a new sorted set
func makeSortedSet() *sortedSet {
	return &sortedSet{
		dict:     make(map[string]*element),
		skipList: makeSkipList(),
	}
}

// makeSkipList makes a new skip list
func makeSkipList() *skipList {
	return &skipList{
		header: makeNode(maxLevel, 0, ""),
		level:  maxLevel,
	}
}

// found pre node and span to next node in each level
func (skipList *skipList) findPre(key string, score float64) ([]*node, []int) {
	// 寻找插入位置每个 level 的前驱节点  以及跨度
	pre := make([]*node, maxLevel)
	rank := make([]int, maxLevel)
	x := skipList.header
	spanCount := 0
	for i := skipList.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score || (x.level[i].forward.score == score && x.level[i].forward.key < key)) {
			x = x.level[i].forward
			spanCount++
		}
		if i == skipList.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1] + spanCount
		}
		spanCount = 0
		pre[i] = x
	}
	return pre, rank
}

// insert into skip list
func (skipList *skipList) insert(key string, score float64) {
	pre, rank := skipList.findPre(key, score)
	// 生成随机层数 l   0～l 层插入节点
	level := randomLevel()
	x := makeNode(level, score, key)
	for i := 0; i < level; i++ {
		// forward 关系
		x.level[i].forward = pre[i].level[i].forward
		pre[i].level[i].forward = x
		// 节点跨度 span
		span1 := rank[0] - rank[i] + 1
		span2 := pre[i].level[i].span - span1 + 1
		x.level[i].span = span2
		pre[i].level[i].span = span1
	}
	// backward 关系
	if pre[0] == skipList.header {
		x.backward = nil
	} else {
		x.backward = pre[0]
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		skipList.tail = x
	}
	skipList.length++
}

// randomLevel returns a random level
func randomLevel() int {
	return rand.Int()%maxLevel + 1
}

// forEach
func (skipList *skipList) forEach(consumer Consumer) {
	x := skipList.header.level[0].forward
	for x != nil {
		if !consumer(x) {
			break
		}
		x = x.level[0].forward
	}
}

// range
func (skipList *skipList) rangeScope(start int, stop int, desc bool) []*element {
	if stop == -1 {
		stop = skipList.length - 1
	}
	if start < 0 || stop < 0 || start > stop || start > skipList.length || stop > skipList.length {
		return nil
	}
	elements := make([]*element, stop-start+1)
	cur := skipList.header.level[0].forward
	for i := 0; i < start; i++ {
		cur = cur.level[0].forward
	}
	for i := 0; i < stop-start+1; i++ {
		elements[i] = &element{
			score: cur.score,
			key:   cur.key,
		}
		cur = cur.level[0].forward
	}
	if desc { // reverse elements
		for i := 0; i < len(elements)/2; i++ {
			elements[i], elements[len(elements)-1-i] = elements[len(elements)-1-i], elements[i]
		}
	}
	return elements
}

// remove by key and score
func (skipList *skipList) remove(key string, score float64) bool {
	pre, _ := skipList.findPre(key, score)
	for i := 0; i < maxLevel; i++ {
		if pre[i].level[i].forward == nil || pre[i].level[i].forward.score != score || pre[i].level[i].forward.key != key {
			break
		}
		pre[i].level[i].span += pre[i].level[i].forward.level[i].span - 1
		pre[i].level[i].forward = pre[i].level[i].forward.level[i].forward
	}
	if pre[0].level[0].forward != nil {
		pre[0].level[0].forward.backward = pre[0]
	} else {
		skipList.tail = pre[0]
	}
	skipList.length--
	return true
}

// update score by key
func (skipList *skipList) update(key string, odlScore, newScore float64) bool {
	skipList.remove(key, odlScore)
	skipList.insert(key, newScore)
	return true
}
