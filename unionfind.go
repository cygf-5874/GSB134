// Package unionfind 实现带权并查集：路径压缩、按秩合并、带权 Find，
// 以及基于快照栈的回滚。
//
// 行为契约见 README 的「对外保证」一节。
package unionfind

// DisjointSet 是带权并查集。
type DisjointSet struct {
	parent []int
	weight []int // weight[i]：i 到 parent[i] 的权值
	rank   []int

	// Find 结果缓存
	cacheRoot map[int]int
	cacheW    map[int]int

	// 快照栈
	stack []snap
}

// New 创建含 n 个元素（编号 0..n-1）的并查集，初始各自成集合、权值为 0。
func New(n int) *DisjointSet {
	d := &DisjointSet{
		parent:    make([]int, n),
		weight:    make([]int, n),
		rank:      make([]int, n),
		cacheRoot: map[int]int{},
		cacheW:    map[int]int{},
	}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

// Find 返回元素 x 所在分量的代表元，以及 x 到代表元的权值（pot(x)-pot(root)）。
func (d *DisjointSet) Find(x int) (int, int) {
	return d.find(x, true)
}

// find 是 Find 的内部实现；useCache 为 false 时读写都不经过缓存。
func (d *DisjointSet) find(x int, useCache bool) (int, int) {
	if useCache {
		if r, ok := d.cacheRoot[x]; ok {
			return r, d.cacheW[x]
		}
	}

	root, w := x, 0
	for d.parent[root] != root {
		w += d.weight[root]
		root = d.parent[root]
	}

	cur := x
	for d.parent[cur] != cur {
		next := d.parent[cur]
		d.parent[cur] = root
		d.weight[cur] = w
		cur = next
	}

	if useCache {
		d.cacheRoot[x] = root
		d.cacheW[x] = w
	}
	return root, w
}

// Union 合并 a、b 所在的两个分量，并使 pot(b)-pot(a) == w。
// 返回是否发生了合并（a、b 已同属一个分量或非法输入时为 false）。
func (d *DisjointSet) Union(a, b, w int) bool {
	if a == b {
		d.parent[a] = b
		d.weight[a] = w
		return false
	}

	ra, wa := d.find(a, false)
	rb, wb := d.find(b, false)
	if ra == rb {
		return false
	}

	if d.rank[ra] < d.rank[rb] {
		d.rank[rb]++
	} else {
		d.rank[ra]++
	}
	d.parent[ra] = rb
	d.weight[ra] = wb - wa - w
	return true
}

// Connected 报告 a、b 是否属于同一分量。
func (d *DisjointSet) Connected(a, b int) bool {
	ra, _ := d.find(a, false)
	rb, _ := d.find(b, false)
	return ra == rb
}
