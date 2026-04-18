package test

import (
	"os"
	"testing"

	"github.com/cheryl-chun/confgen/runtime"
	"github.com/spf13/viper"
)

func getConfigPath() string {
	if _, err := os.Stat("test/config.yaml"); err == nil {
		return "test/config.yaml"
	}
	if _, err := os.Stat("config.yaml"); err == nil {
		return "config.yaml"
	}
	return "test/config.yaml"
}

// BenchmarkGenerated_Confgen_Load: Load Config
func BenchmarkGenerated_Confgen_Load(b *testing.B) {
	configPath := getConfigPath()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cfg := &Config{}
		loader := runtime.NewLoader().AddFile(configPath)
		if err := loader.Fill(cfg); err != nil {
			b.Fatalf("failed to load config: %v", err)
		}
	}
}

// BenchmarkGenerated_Viper_Load: Load Config by Viper
func BenchmarkGenerated_Viper_Load(b *testing.B) {
	configPath := getConfigPath()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v := viper.New()
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			b.Fatalf("failed to read config: %v", err)
		}

		var cfg Config
		if err := v.Unmarshal(&cfg); err != nil {
			b.Fatalf("failed to unmarshal: %v", err)
		}
	}
}

// ============= Runtime bench =============

func BenchmarkGenerated_Confgen_FieldAccess(b *testing.B) {
	cfg := &Config{}
	loader := runtime.NewLoader().AddFile(getConfigPath())
	if err := loader.Fill(cfg); err != nil {
		b.Fatalf("failed to load config: %v", err)
	}

	for i := 0; i < b.N; i++ {
		_ = cfg.App.Name
		_ = cfg.Server.Host
		_ = cfg.Server.Port
		_ = cfg.Database.Host
		_ = cfg.Cache.RedisPort
		_ = cfg.Features.EnableMetrics
	}
}

func BenchmarkGenerated_Viper_FieldAccess(b *testing.B) {
	v := viper.New()
	v.SetConfigFile(getConfigPath())
	if err := v.ReadInConfig(); err != nil {
		b.Fatalf("failed to read config: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = v.GetString("app.name")
		_ = v.GetString("server.host")
		_ = v.GetInt("server.port")
		_ = v.GetString("database.host")
		_ = v.GetInt("cache.redis_port")
		_ = v.GetBool("features.enable_metrics")
	}
}

func BenchmarkGenerated_Confgen_BatchAccess(b *testing.B) {
	cfg := &Config{}
	loader := runtime.NewLoader().AddFile(getConfigPath())
	if err := loader.Fill(cfg); err != nil {
		b.Fatalf("failed to load config: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 30; j++ {
			_ = cfg.App.Name
			_ = cfg.Server.Host
			_ = cfg.Server.Port
			_ = cfg.Database.Host
			_ = cfg.Cache.RedisPort
		}
	}
}

func BenchmarkGenerated_Viper_BatchAccess(b *testing.B) {
	v := viper.New()
	v.SetConfigFile(getConfigPath())
	if err := v.ReadInConfig(); err != nil {
		b.Fatalf("failed to read config: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 30; j++ {
			_ = v.GetString("app.name")
			_ = v.GetString("server.host")
			_ = v.GetInt("server.port")
			_ = v.GetString("database.host")
			_ = v.GetInt("cache.redis_port")
		}
	}
}

func BenchmarkGenerated_Confgen_AllFields(b *testing.B) {
	cfg := &Config{}
	loader := runtime.NewLoader().AddFile(getConfigPath())
	if err := loader.Fill(cfg); err != nil {
		b.Fatalf("failed to load config: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		walkGeneratedConfigConfgen(cfg)
	}
}

func BenchmarkGenerated_Viper_AllFields(b *testing.B) {
	v := viper.New()
	v.SetConfigFile(getConfigPath())
	if err := v.ReadInConfig(); err != nil {
		b.Fatalf("failed to read config: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		walkGeneratedConfigViper(v)
	}
}


func walkGeneratedConfigConfgen(cfg *Config) {
	_ = cfg.App.Name
	_ = cfg.App.Version
	_ = cfg.App.LogLevel
	_ = cfg.Server.Host
	_ = cfg.Server.Port
	_ = cfg.Server.TimeoutMs
	_ = cfg.Server.SSLEnabled
	_ = cfg.Database.Driver
	_ = cfg.Database.Host
	_ = cfg.Database.Port
	_ = cfg.Database.User
	_ = cfg.Database.Password
	_ = cfg.Database.Dbname
	_ = cfg.Database.PoolSize
	_ = cfg.Cache.RedisHost
	_ = cfg.Cache.RedisPort
	_ = cfg.Cache.TtlSeconds
	_ = cfg.Features.EnableMetrics
	_ = cfg.Features.EnableTracing
	_ = cfg.Features.RateLimitRps
}

func walkGeneratedConfigViper(v *viper.Viper) {
	_ = v.GetString("app.name")
	_ = v.GetString("app.version")
	_ = v.GetString("app.log_level")
	_ = v.GetString("server.host")
	_ = v.GetInt("server.port")
	_ = v.GetInt("server.timeout_ms")
	_ = v.GetBool("server.ssl_enabled")
	_ = v.GetString("database.driver")
	_ = v.GetString("database.host")
	_ = v.GetInt("database.port")
	_ = v.GetString("database.user")
	_ = v.GetString("database.password")
	_ = v.GetString("database.dbname")
	_ = v.GetInt("database.pool_size")
	_ = v.GetString("cache.redis_host")
	_ = v.GetInt("cache.redis_port")
	_ = v.GetInt("cache.ttl_seconds")
	_ = v.GetBool("features.enable_metrics")
	_ = v.GetBool("features.enable_tracing")
	_ = v.GetInt("features.rate_limit_rps")
}

func TestCorrectness_ConfgenVsViper(t *testing.T) {
	configPath := getConfigPath()

	// 使用 confgen 加载配置
	confgenCfg := &Config{}
	loader := runtime.NewLoader().AddFile(configPath)
	if err := loader.Fill(confgenCfg); err != nil {
		t.Fatalf("confgen failed to load config: %v", err)
	}

	// 使用 viper 加载配置
	v := viper.New()
	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		t.Fatalf("viper failed to read config: %v", err)
	}
	var viperCfg Config
	if err := v.Unmarshal(&viperCfg); err != nil {
		t.Fatalf("viper failed to unmarshal: %v", err)
	}

	tests := []struct {
		name        string
		confgenVal  interface{}
		viperVal    interface{}
		description string
	}{
		{"app.name", confgenCfg.App.Name, viperCfg.App.Name, "应用名称"},
		{"app.version", confgenCfg.App.Version, viperCfg.App.Version, "应用版本"},
		{"app.log_level", confgenCfg.App.LogLevel, viperCfg.App.LogLevel, "日志级别"},

		{"server.host", confgenCfg.Server.Host, viperCfg.Server.Host, "服务器主机"},
		{"server.port", confgenCfg.Server.Port, viperCfg.Server.Port, "服务器端口"},
		{"server.timeout_ms", confgenCfg.Server.TimeoutMs, viperCfg.Server.TimeoutMs, "超时时间"},
		{"server.ssl_enabled", confgenCfg.Server.SSLEnabled, viperCfg.Server.SSLEnabled, "SSL 启用"},

		{"database.driver", confgenCfg.Database.Driver, viperCfg.Database.Driver, "数据库驱动"},
		{"database.host", confgenCfg.Database.Host, viperCfg.Database.Host, "数据库主机"},
		{"database.port", confgenCfg.Database.Port, viperCfg.Database.Port, "数据库端口"},
		{"database.user", confgenCfg.Database.User, viperCfg.Database.User, "数据库用户"},
		{"database.password", confgenCfg.Database.Password, viperCfg.Database.Password, "数据库密码"},
		{"database.dbname", confgenCfg.Database.Dbname, viperCfg.Database.Dbname, "数据库名"},
		{"database.pool_size", confgenCfg.Database.PoolSize, viperCfg.Database.PoolSize, "连接池大小"},

		{"cache.redis_host", confgenCfg.Cache.RedisHost, viperCfg.Cache.RedisHost, "Redis 主机"},
		{"cache.redis_port", confgenCfg.Cache.RedisPort, viperCfg.Cache.RedisPort, "Redis 端口"},
		{"cache.ttl_seconds", confgenCfg.Cache.TtlSeconds, viperCfg.Cache.TtlSeconds, "TTL 秒数"},

		{"features.enable_metrics", confgenCfg.Features.EnableMetrics, viperCfg.Features.EnableMetrics, "度量启用"},
		{"features.enable_tracing", confgenCfg.Features.EnableTracing, viperCfg.Features.EnableTracing, "追踪启用"},
		{"features.rate_limit_rps", confgenCfg.Features.RateLimitRps, viperCfg.Features.RateLimitRps, "速率限制"},
	}

	failures := 0
	for _, tt := range tests {
		if tt.confgenVal != tt.viperVal {
			t.Errorf("❌ %s (%s): confgen=%v, viper=%v", tt.name, tt.description, tt.confgenVal, tt.viperVal)
			failures++
		} else {
			t.Logf("✓ %s: %v", tt.name, tt.confgenVal)
		}
	}

	if failures > 0 {
		t.Fatalf("\n❌ test fialed, %d filed error\n", failures)
	}
}
