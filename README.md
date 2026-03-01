# cache-redis

`cache-redis` 是 `cache` 模块的 `redis` 驱动。

## 安装

```bash
go get github.com/infrago/cache@latest
go get github.com/infrago/cache-redis@latest
```

## 接入

```go
import (
    _ "github.com/infrago/cache"
    _ "github.com/infrago/cache-redis"
    "github.com/infrago/infra"
)

func main() {
    infra.Run()
}
```

## 配置示例

```toml
[cache]
driver = "redis"
```

## 公开 API（摘自源码）

- `func (d *redisDriver) Connect(inst *cache.Instance) (cache.Connect, error)`
- `func (c *redisConnection) Open() error  { return nil }`
- `func (c *redisConnection) Close() error { return c.client.Close() }`
- `func (c *redisConnection) Read(key string) ([]byte, error)`
- `func (c *redisConnection) Write(key string, val []byte, expire time.Duration) error`
- `func (c *redisConnection) Exists(key string) (bool, error)`
- `func (c *redisConnection) Delete(key string) error`
- `func (c *redisConnection) Sequence(key string, start, step int64, expire time.Duration) (int64, error)`
- `func (c *redisConnection) Keys(prefix string) ([]string, error)`
- `func (c *redisConnection) Clear(prefix string) error`

## 排错

- driver 未生效：确认模块段 `driver` 值与驱动名一致
- 连接失败：检查 endpoint/host/port/鉴权配置
