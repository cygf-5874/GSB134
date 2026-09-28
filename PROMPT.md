这份并查集要接手之前得先审一遍。unionfind 是 Go 1.24 的带权并查集库，仅标准库，自检走 scripts/check.sh。
模块名即库名 `unionfind`，`go.mod` 不写 `require`；`go test ./...` 跑既有用例，构建走 `bash scripts/build.sh`。

本次不写代码：通读 `unionfind.go` 与 `rollback.go`，把 `REVIEW.md` 的表填满 ——
每个问题一行，写清它违反 README「对外保证」的哪一条、出在哪个函数、以及一句可判定的证据。
源码一个字节都不许改。

验收：
- go test ./... 全绿；
- bash scripts/check.sh 退出码 0，8 个场景全过（unmodified 1 + format 1 + finding 6）。

约束：
1. 不改 `check/`，不改 `unionfind.go` / `rollback.go`；`check/` 会校验两个源码文件的 SHA-256。
2. `REVIEW.md` 表头与编号 `D1`..`D6` 已给定，每行四列都要填。
3. 「可判定证据」必须含具体输入或计数，不许写「可能 / 看起来」这类无法核对的措辞。
4. 仅标准库，不引入任何第三方依赖。