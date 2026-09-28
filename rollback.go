// rollback.go：并查集的快照与回滚。
//
// 行为契约见 README 的「对外保证」一节（第 4、5 条）。
package unionfind

type snap struct {
	parent []int
	rank   []int
}

// Snapshot 记下当前状态并入快照栈，返回快照编号（从 0 起递增）。
func (d *DisjointSet) Snapshot() int {
	s := snap{
		parent: append([]int(nil), d.parent...),
		rank:   append([]int(nil), d.rank...),
	}
	d.stack = append(d.stack, s)
	return len(d.stack) - 1
}

// Rollback 回滚到编号为 id 的快照。
// id 非法（越界）时不改变任何状态。
func (d *DisjointSet) Rollback(id int) {
	if id < 0 || id >= len(d.stack) {
		return
	}
	d.stack = d.stack[:id+1]
	copy(d.rank, d.stack[id].rank)
}

// Snapshots 返回当前快照栈的深度。
func (d *DisjointSet) Snapshots() int {
	return len(d.stack)
}
