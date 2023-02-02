package consistent

import (
	"errors"
	"fmt"
	"hash/crc32"
	"sort"
)

type Consistent struct {
	nodesReplicas int            //添加虚拟节点数量
	hashSortNodes []uint32       //所有节点的排序数组
	circle        map[uint32]int //所有节点对应的node
	nodes         map[int]bool   //node是否存在
}

func NewConsistent(nodesReplicas int) *Consistent {
	return &Consistent{
		nodesReplicas: nodesReplicas,
		circle:        make(map[uint32]int),
		nodes:         make(map[int]bool),
	}
}

func (c *Consistent) add(node int) error {
	if _, ok := c.nodes[node]; ok { //判断新加节点是否存在
		return fmt.Errorf("%d already existed", node)
	}
	c.nodes[node] = true
	for i := 0; i < c.nodesReplicas; i++ { //添加虚拟节点
		replicasKey := getReplicasKey(i, node) //虚拟节点
		c.circle[replicasKey] = node
		c.hashSortNodes = append(c.hashSortNodes, replicasKey) //所有节点ID
	}
	sort.Slice(c.hashSortNodes, func(i, j int) bool { //排序
		return c.hashSortNodes[i] < c.hashSortNodes[j]
	})
	return nil
}

func (c *Consistent) remove(node int) error {
	if _, ok := c.nodes[node]; !ok { //判断新加节点是否存在
		return fmt.Errorf("%d not existed", node)
	}
	delete(c.nodes, node)
	for i := 0; i < c.nodesReplicas; i++ {
		replicasKey := getReplicasKey(i, node)
		delete(c.circle, replicasKey) //删除虚拟节点
	}
	c.refreshHashSortNodes()
	return nil
}
func (c *Consistent) GetNode() (node []int) {
	for v := range c.nodes {
		node = append(node, v)
	}
	return node
}

func (c *Consistent) Get(key string, partitions int) (int, error) {
	c.reBlanceNode(partitions)
	if len(c.nodes) == 0 {
		return 0, errors.New("not add node")
	}
	index := c.searchNearbyIndex(key)
	host := c.circle[c.hashSortNodes[index]]
	return host, nil
}

func (c *Consistent) reBlanceNode(partitions int) {
	nodeCount := len(c.nodes)

	switch {
	case partitions > nodeCount:
		for i := partitions - 1; i >= nodeCount; i-- {
			c.add(i)
		}
	case partitions < nodeCount:
		for i := partitions; i < nodeCount; i++ {
			c.remove(i)
		}
	default:
	}

}

func (c *Consistent) refreshHashSortNodes() {
	c.hashSortNodes = nil
	for v := range c.circle {
		c.hashSortNodes = append(c.hashSortNodes, v)
	}
	sort.Slice(c.hashSortNodes, func(i, j int) bool { //排序
		return c.hashSortNodes[i] < c.hashSortNodes[j]
	})
}

func (c *Consistent) searchNearbyIndex(key string) int {
	hashKey := hashKey(key)
	index := sort.Search(len(c.hashSortNodes), func(i int) bool { //key算出的节点，距离最近的node节点下标  sort.Search数组需要升序排列好
		return c.hashSortNodes[i] >= hashKey
	})
	if index >= len(c.hashSortNodes) {
		index = 0
	}
	return index
}

func getReplicasKey(i int, node int) uint32 {
	return hashKey(fmt.Sprintf("%d#%d", node, i))
}

func hashKey(host string) uint32 {
	return crc32.ChecksumIEEE([]byte(host))
}
