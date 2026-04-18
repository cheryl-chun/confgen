package runtime

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cheryl-chun/confgen/internal/tree"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdSource 从 etcd 拉取配置，并将键值写入配置树。
//
// key 前缀支持通过 Prefix 字段配置。
type EtcdSource struct {
	RemoteConfigSource
	Endpoints    []string
	DialTimeout  time.Duration
	FetchTimeout time.Duration
}

// NewEtcdSource 创建一个 etcd 配置源。
func NewEtcdSource(endpoints []string, prefix string) *EtcdSource {
	return &EtcdSource{
		RemoteConfigSource: RemoteConfigSource{Prefix: prefix},
		Endpoints:          endpoints,
		DialTimeout:        5 * time.Second,
		FetchTimeout:       10 * time.Second,
	}
}

// Load 实现 Source 接口
func (s *EtcdSource) Load(configTree *tree.ConfigTree) error {
	if len(s.Endpoints) == 0 {
		return fmt.Errorf("etcd endpoints are required")
	}

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   s.Endpoints,
		DialTimeout: s.DialTimeout,
	})
	if err != nil {
		return fmt.Errorf("failed to create etcd client: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), s.FetchTimeout)
	defer cancel()

	prefix := s.normalizePrefix()
	if prefix == "" {
		prefix = "/"
	} else {
		prefix = "/" + prefix
	}

	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("failed to fetch etcd keys: %w", err)
	}

	for _, kv := range resp.Kvs {
		path := s.KeyToPath(string(kv.Key))
		if path == "" {
			continue
		}

		value, valueType := inferRemoteValueType(string(kv.Value))
		if err := configTree.Set(path, value, tree.SourceRemote, valueType); err != nil {
			return err
		}
	}

	return nil
}

// Priority 返回 etcd 配置源的优先级。
func (s *EtcdSource) Priority() tree.SourceType {
	return tree.SourceRemote
}

func inferRemoteValueType(raw string) (any, tree.ValueType) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw, tree.TypeString
	}

	if b, err := strconv.ParseBool(raw); err == nil {
		return b, tree.TypeBool
	}

	if i, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return int(i), tree.TypeInt
	}

	if f, err := strconv.ParseFloat(raw, 64); err == nil {
		return f, tree.TypeFloat
	}

	return raw, tree.TypeString
}
