package unionfind

import "testing"

func TestSelfFind(t *testing.T) {
	d := New(5)
	r, w := d.Find(2)
	if r != 2 || w != 0 {
		t.Fatalf("Find(2) = (%d,%d)，期望 (2,0)", r, w)
	}
}

func TestUnionWeights(t *testing.T) {
	d := New(4)
	if !d.Union(0, 1, 7) {
		t.Fatalf("Union(0,1,7) = false，期望 true")
	}
	if d.Union(0, 1, 7) {
		t.Fatalf("重复 Union 应为 false")
	}
	r1, w1 := d.Find(1)
	if r1 != 1 || w1 != 0 {
		t.Fatalf("Find(1) = (%d,%d)，期望 (1,0)", r1, w1)
	}
	r0, w0 := d.Find(0)
	if r0 != 1 || w0 != -7 {
		t.Fatalf("Find(0) = (%d,%d)，期望 (1,-7)", r0, w0)
	}
}

func TestConnected(t *testing.T) {
	d := New(3)
	if !d.Union(0, 2, 5) {
		t.Fatalf("Union(0,2,5) = false，期望 true")
	}
	if !d.Connected(0, 2) {
		t.Fatalf("Connected(0,2) = false，期望 true")
	}
	if d.Connected(0, 1) {
		t.Fatalf("Connected(0,1) = true，期望 false")
	}
}

func TestSnapshotRollback(t *testing.T) {
	d := New(3)
	d.Union(0, 1, 4)
	id := d.Snapshot()
	d.Rollback(id)
	if n := d.Snapshots(); n != 1 {
		t.Fatalf("Snapshots() = %d，期望 1", n)
	}
	r, _ := d.Find(0)
	if r != 1 {
		t.Fatalf("回滚后 Find(0) 代表元 = %d，期望 1", r)
	}
}

func TestSnapshotIDsIncrease(t *testing.T) {
	d := New(2)
	a := d.Snapshot()
	b := d.Snapshot()
	if b <= a {
		t.Fatalf("快照编号未递增：%d then %d", a, b)
	}
}

func TestRollbackInvalidIDNoop(t *testing.T) {
	d := New(2)
	d.Union(0, 1, 1)
	d.Rollback(-1)
	d.Rollback(99)
	if !d.Connected(0, 1) {
		t.Fatalf("非法回滚后状态被改变")
	}
}
