package runtime

import (
	"strings"

	"github.com/cheryl-chun/confgen/internal/tree"
)

// RemoteConfigSource 是所有远程配置中心的基础类型，负责
// 处理统一的 key 前缀和 key->path 转换逻辑。
type RemoteConfigSource struct {
	Prefix string
}

// normalizePrefix 返回标准化的前缀形式，确保可用于前缀匹配。
func (s *RemoteConfigSource) normalizePrefix() string {
	return strings.Trim(s.Prefix, "/.")
}

// KeyToPath 将远程配置中心的 key 转换为配置树路径。
// 支持 / 分隔符和 . 分隔符，并且会自动移除指定前缀。
func (s *RemoteConfigSource) KeyToPath(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}

	prefix := s.normalizePrefix()
	if prefix != "" {
		trimmedKey := strings.TrimLeft(key, "/.")
		trimmedPrefix := strings.TrimLeft(prefix, "/.")
		if strings.HasPrefix(trimmedKey, trimmedPrefix) {
			key = strings.TrimPrefix(trimmedKey, trimmedPrefix)
		} else {
			key = trimmedKey
		}
	}

	key = strings.TrimLeft(key, "/.")
	key = strings.TrimRight(key, "/.")
	if key == "" {
		return ""
	}

	key = strings.ReplaceAll(key, "/", ".")
	key = strings.ReplaceAll(key, "_", ".")
	key = strings.Trim(key, ".")
	return strings.ToLower(key)
}

// PathToKey 将配置树路径转换为远程配置中心的 key。
func (s *RemoteConfigSource) PathToKey(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return s.normalizePrefix()
	}

	key := strings.ReplaceAll(path, ".", "/")
	key = strings.Trim(key, "/")
	prefix := s.normalizePrefix()
	if prefix == "" {
		return key
	}
	return prefix + "/" + key
}

// ConfigurePrefix 设置前缀。
func (s *RemoteConfigSource) ConfigurePrefix(prefix string) {
	s.Prefix = prefix
}

// Priority 返回远程配置中心默认优先级。
func (s *RemoteConfigSource) Priority() tree.SourceType {
	return tree.SourceRemote
}
