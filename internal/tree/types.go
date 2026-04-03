package tree

type ValueType int

const (
	TypeString ValueType = iota
	TypeInt
	TypeFloat
	TypeBool
	TypeArray
	TypeObject
	TypeNull
)

// SourceType 配置来源类型
type SourceType int

const (
	SourceDefault      SourceType = iota // Default value
	SourceRemote                         // Remote configuration center (etcd/zookeeper/consul, etc.)
	SourceFile                           // Configuration file (yaml/json/toml, etc.)
	SourceRuntimeOverride                // Runtime dynamic override (via API, etc.)
	SourceSessionEnv                     // Session environment variable (current process)
	SourceSystemEnv                      // System environment variable (global persistent)
	SourceCodeOverride                   // Code explicit override (Set method)
)

// SourcePriority 配置源的优先级（数字越大优先级越高）
//
// 优先级设计原则（按你的建议）：
// 1. 系统环境变量 > 会话环境变量：持久配置 > 临时配置
// 2. 环境变量 > 配置文件：符合 12-Factor App，便于运维
// 3. 配置文件 > 远程配置：本地优先，便于开发调试
// 4. 远程配置 > 默认值：动态配置 > 硬编码
//
// 特殊说明：
// - CodeOverride：代码中显式 Set()，优先级最高（程序员明确意图）
// - RuntimeOverride：运行时 API 动态修改，优先级介于配置文件和远程之间
//
// 典型场景：
// - 开发环境：SystemEnv > File（本地配置）
// - 测试环境：SystemEnv > File
// - 生产环境：SystemEnv > SessionEnv > File > Remote
var SourcePriority = map[SourceType]int{
	SourceDefault:         0,   // 优先级最低：硬编码默认值
	SourceRemote:          10,  // 远程配置中心（etcd/zookeeper/consul）
	SourceRuntimeOverride: 15,  // 运行时动态修改（介于远程和配置文件之间）
	SourceFile:            20,  // 配置文件（本地）
	SourceSessionEnv:      30,  // 会话环境变量（临时）
	SourceSystemEnv:       40,  // 系统环境变量（持久）
	SourceCodeOverride:    100, // 优先级最高：代码显式设置
}

func (t ValueType) String() string {
	names := []string{
		"String", "Int", "Float", "Bool", "Array", "Object", "Null",
	}
	if t >= 0 && int(t) < len(names) {
		return names[t]
	}
	return "Unknown"
}

func (s SourceType) String() string {
	names := []string{
		"Default",
		"Remote",
		"File",
		"RuntimeOverride",
		"SessionEnv",
		"SystemEnv",
		"CodeOverride",
	}
	if s >= 0 && int(s) < len(names) {
		return names[s]
	}
	return "Unknown"
}
