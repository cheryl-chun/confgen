package runtime

import (
	"os"
	"strings"

	"github.com/cheryl-chun/confgen/internal/tree"
)

// EnvSource 环境变量配置源
type EnvSource struct {
	Prefix string // 环境变量前缀，如 "APP_"
}

// Load 实现 Source 接口
func (s *EnvSource) Load(configTree *tree.ConfigTree) error {
	// 遍历所有环境变量
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		// 检查前缀
		if s.Prefix != "" && !strings.HasPrefix(key, s.Prefix) {
			continue
		}

		// 移除前缀
		if s.Prefix != "" {
			key = strings.TrimPrefix(key, s.Prefix)
		}

		// 转换环境变量名为配置路径
		// APP_SERVER_HOST -> server.host
		path := s.envKeyToPath(key)

		// 设置到 tree（环境变量优先级高）
		configTree.Set(path, value, tree.SourceSystemEnv, tree.TypeString)
	}

	return nil
}

// Priority 返回环境变量配置源的优先级
func (s *EnvSource) Priority() tree.SourceType {
	return tree.SourceSystemEnv
}

// envKeyToPath 将环境变量名转换为配置路径
// APP_SERVER_HOST -> server.host
func (s *EnvSource) envKeyToPath(key string) string {
	// 转小写并替换 _ 为 .
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, "_", ".")
	return key
}