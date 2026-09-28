# unionfind

带权并查集 + 快照回滚（Go 1.24，仅标准库，`go.mod` 无 `require`）。

`Find(x)` 返回 `(代表元, x 到代表元的权值)`；`Union(a, b, w)` 合并两个分量并满足
`pot(b) - pot(a) == w`；`Snapshot` / `Rollback` 提供快照栈式的回滚。

```go
d := unionfind.New(5)
d.Union(0, 1, 7)      // pot(1) - pot(0) == 7
r, w := d.Find(0)     // r 是代表元，w 是 0 到代表元的权值
id := d.Snapshot()
d.Rollback(id)
```

## 怎么跑

```bash
go build ./...            # 编译
go test ./...             # 既有用例
go run ./check            # 固定验收程序（8 个场景）
go run ./check -list      # 列出全部场景
go run ./check --only finding
bash scripts/check.sh     # 薄封装：go run ./check "$@"
```

前提：Go 1.24。`check/`、`scripts/` 是固定验收的一部分，**请勿修改**。

## 本次交付物

这次**不写代码**：请通读 `unionfind.go` 与 `rollback.go`，把 `REVIEW.md` 里的那张表
**填满** —— 每个问题一行，写清它**违反哪一条保证**、出在**哪个函数**、以及**一句可判定的证据**。

- 表头已给定：`| 编号 | 违反保证 | 函数 | 可判定证据 |`，编号固定为 `D1`..`D6`。
- `违反保证` 填 1~8 的**编号**；`函数` 填导出的函数名（`New` / `Find` / `Union` /
  `Connected` / `Snapshot` / `Rollback`）。
- `可判定证据` 必须**可复现**：给出一句话，含具体的输入形状、计数或可观察差异，
  不许写「看起来不对」「可能有风险」这类无法核对的话。
- **源码一个字节都不许改**（`unionfind.go`、`rollback.go`），`check/` 会校验 SHA-256。

## 对外保证

下面 8 条是这份实现的契约。**它们是契约，不是「当前行为」的转述。**

1. **路径压缩 + 按秩合并**：`Union` / `Find` 必须做路径压缩，并按**秩**决定合并方向，
   使树高保持接近对数级；不得让秩的记录与合并方向脱节而退化成链。
2. **代表元与秩语义**：`Find(x)` 返回的**代表元**必须是所在分量的根；同一分量的任意元素
   返回同一代表元。`rank` 仅用于合并决策，其取值必须与子树规模一致。
3. **带权语义**：`Find(x)` 返回 `(root, w)`，其中 `w == pot(x) - pot(root)`；
   `Union(a, b, w)` 建立关系 `pot(b) - pot(a) == w`。合并与压缩之后，任意元素再次 `Find`
   得到的权值必须与上述关系**逐元素一致**。
4. **回滚一致**：`Rollback(id)` 之后，**所有元素**的父指针、权值、秩必须与编号 `id` 的快照
   时刻**逐元素相同**（包括被路径压缩改动过的节点）。
5. **快照栈语义**：`Snapshot()` 依次压栈并返回递增编号；`Rollback(id)` 回滚到该快照并丢弃
   其后的快照；非法编号不得改变状态。
6. **确定性**：结果不得依赖 `map` 迭代顺序；同一操作序列重复执行结果相同。
7. **边界**：空集、单元素、自环。`Union(x, x, w)`（自环）必须被拒绝（返回 `false`）且
   **不改变任何状态**；空集上调用 `Find` 之外的非法索引应可预期，不得 panic。
8. **纯函数性质**：除显式的 `Union` 与 `Rollback` 外，任何方法都不得改变集合的内部状态
   （例如 `Find` 可以压缩本分量，但不得留下跨调用的全局副作用）。

## 固定件场景（8 个）

| 分组 | 场景 | 覆盖 |
| --- | --- | --- |
| `unmodified` | `unmodified_sources` | 2 个不可修改源码文件的 SHA-256 未变 |
| `format` | `table_shaped` | `REVIEW.md` 的表格式合法、恰好 6 行且四列齐全 |
| `finding` | `D1`..`D6` | 每个缺陷的「违反保证」编号与「函数」命中允许集合，证据可判定 |

`--only <组名>`（`unmodified` / `format` / `finding`）可以只跑一组。失败不早退。

## 目录

```
.
├── .gitattributes
├── .gitignore
├── go.mod               module unionfind
├── unionfind.go         待评审的实现（只审不改）
├── rollback.go          待评审的实现（只审不改）
├── unionfind_test.go    既有用例
├── PROMPT.md            本题的 User Prompt（逐字）
├── README.md
├── REVIEW.md            本次交付物（把表填满）
├── check/main.go        固定验收程序（勿改）
└── scripts/
    ├── build.sh         go build ./...
    └── check.sh         薄封装：go run ./check "$@"
```