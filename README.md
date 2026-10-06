# api-key-scope

API key 的签发、权限校验、过期与级联吊销，另有一层按所有者放行的覆盖层。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/apikey --sample=idem
go run ./cmd/apikey --sample=cascade
go run ./cmd/apikey --sample=expiry
go run ./cmd/apikey --sample=override
go run ./cmd/apikey --sample=work
```

## 口径

- **签发幂等**：同一个 key 再次签发不改变已经登记的权限与过期时刻。
- **校验**：key 存在、未吊销、未过期，且动作在 key 的权限里，或在该所有者的覆盖层里，才算允许。
- **过期**：`expires` 大于 0 时，当前时刻达到 `expires` 即视为过期。
- **级联吊销**：按所有者吊销后，该所有者名下所有 key 都失效。

## 不变量

- 一个 key 只有一条登记，权限在签发后不再被改写。
- 被吊销或过期的 key 一律拒绝，覆盖层也不能让它复活。
- `scanned` 不随 key 数平方级放大：两千个 key 共四千次调用，`scanned` 不超过 6000。

## 输出契约

场景打印一行 JSON，键随场景而定。`--sample=work` 打印 `{"keys": N, "scanned": N}`。
