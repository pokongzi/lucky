package redis

// RedisConfig 结构体表示 Redis 配置
type RedisConfig struct {
	Host        string `ini:"host"`
	Port        int    `ini:"port"`
	Enabled     bool   `ini:"enabled"`
	PoolSize    int    `ini:"pool_size"`
	RequirePass string `ini:"password"`
	MaxClients  int    `ini:"max_clients"`
	Db          int    `ini:"db"`
	MinIdleConns int    `ini:"min_idle_conns"`
}
