package database

import "time"

type ConnectionConfig struct {
	ConnectionString      string        `yaml:"conn_string"`
	ConnectionTimeout     time.Duration `yaml:"conn_timeout"`
	UsePGBouncer          bool          `yaml:"use_pgbouncer"`
	PoolMinConns          int32         `yaml:"pool_min_conns"`
	PoolMaxConns          int32         `yaml:"pool_max_conns"`
	PoolHealthCheckPeriod time.Duration `yaml:"pool_health_check_period"`
}
