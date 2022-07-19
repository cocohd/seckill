package common

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

/*一致性hash为服务器节点做负载均衡*/

// 声明新切片类型
type units []uint32

func (u units) Len() int {
	return len(u)
}

func (u units) Less(i, j int) bool {
	return u[i] < u[j]
}

func (u units) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

// ConsistentHash 一致性hash结构体
type ConsistentHash struct {
	circle      map[uint32]string
	sortedKey   units
	virtualNode int
	mutex       *sync.Mutex
}

func NewConsistentHash(replicas int) *ConsistentHash {
	return &ConsistentHash{circle: make(map[uint32]string), virtualNode: replicas}
}

// GenerateKey 对于一个节点，生成其对应的虚拟节点
func (c *ConsistentHash) GenerateKey(element string, virtualIndex int) string {
	return element + strconv.Itoa(virtualIndex)
}

// 计算hash值
func (c *ConsistentHash) hash(key string) uint32 {
	// 统一hash计算长度
	//if len(key) < 64 {
	//	tmp := [64]byte
	//	copy(tmp[:], key)
	//	return crc32.ChecksumIEEE(tmp[])
	//}
	return crc32.ChecksumIEEE([]byte(key))
}

// 更新circle的key排序
func (c *ConsistentHash) updateHashedKeys() {
	keysArr := units{}

	for k := range c.circle {
		keysArr = append(keysArr, k)
	}

	sort.Sort(keysArr)

	c.sortedKey = keysArr
}

// Add 向circle添加节点
func (c *ConsistentHash) Add(element string) {
	// map加锁
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.add(element)
}

func (c *ConsistentHash) add(element string) {
	for i := 0; i < c.virtualNode; i++ {
		// 生成每个虚拟节点并加入map
		c.circle[c.hash(c.GenerateKey(element, i))] = element
	}

	// 对circle中的key排序，方面后续用二分查找
	c.updateHashedKeys()
}

// Remove 删除节点
func (c *ConsistentHash) Remove(element string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.remove(element)
}

func (c *ConsistentHash) remove(element string) {
	for i := 0; i < c.virtualNode; i++ {
		delete(c.circle, c.hash(c.GenerateKey(element+strconv.Itoa(i), i)))
	}

	// 排序，保持keys有序
	c.updateHashedKeys()
}

// 顺时针查找最近的节点
func (c *ConsistentHash) search(key uint32) int {
	idx := sort.Search(len(c.sortedKey), func(i int) bool {
		return c.sortedKey[i] > key
	})

	// 未查找到
	if idx >= len(c.sortedKey) {
		idx = 0
	}
	return idx
}

// Get 得到符合特点的最近节点
func (c *ConsistentHash) Get(name string) (string, error) {
	// 加锁，防止并发读取map
	c.mutex.Lock()
	defer c.mutex.Unlock()

	hashedKey := c.hash(name)
	idx := c.search(hashedKey)

	return c.circle[c.sortedKey[idx]], nil
}
