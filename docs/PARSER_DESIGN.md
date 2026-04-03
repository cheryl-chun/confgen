# Parser 包设计文档

## 📋 概述

Parser 包负责将不同格式的配置文件（YAML/JSON/TOML等）统一解析为标准的树形结构（`ConfigNode` 树），供后续的类型推断和代码生成使用。

## 🏗️ 架构设计

### 设计模式

- **工厂模式** (Factory Pattern): 根据文件扩展名自动选择合适的解析器
- **策略模式** (Strategy Pattern): 不同格式的解析器实现统一的 `Parser` 接口
- **组合模式** (Composite Pattern): `BaseParser` 提供通用功能，具体解析器组合使用
- **单例模式** (Singleton Pattern): 全局解析器工厂

### 核心组件

```
Parser Package
├── Parser Interface          # 统一接口
├── BaseParser               # 基类（复用逻辑）
├── YAMLParser              # YAML 实现
├── JSONParser              # JSON 实现
├── ParserFactory           # 工厂（自动选择）
└── ConfigNode Tree         # 中间表示（IR）
```

## 📄 文件说明

### 1. parser.go - 接口定义

**职责**: 定义解析器的核心接口

```go
type Parser interface {
    Parse(reader io.Reader) (*ParseResult, error)
    ParseFile(path string) (*ParseResult, error)
    SupportedExtensions() []string
    Name() string
}
```

**设计要点**:
- `Parse()` 和 `ParseFile()` 提供两种输入方式
- `BaseParser` 避免重复实现 `Name()` 和 `SupportedExtensions()`

---

### 2. types.go - 数据结构

**职责**: 定义配置文件的中间表示（IR）

#### ConfigNode 树结构

```
输入 (YAML):               转换为 ConfigNode 树:
server:                    Root (TypeObject)
  host: localhost          └── server (TypeObject)
  port: 8080                   ├── host (TypeString)
  features:                    ├── port (TypeInt)
    - ssl                      └── features (TypeArray)
    - cache                        ├── [0] (TypeString)
                                   └── [1] (TypeString)
```

#### 核心类型

```go
// 值类型枚举
type ValueType int
const (
    TypeString, TypeInt, TypeFloat, TypeBool,
    TypeArray, TypeObject, TypeNull
)

// 配置节点
type ConfigNode struct {
    Key      string                 // 字段名
    Value    any                    // 原始值
    Type     ValueType              // 类型
    Children map[string]*ConfigNode // 子对象
    Items    []*ConfigNode          // 数组元素
}

// 解析结果
type ParseResult struct {
    Root *ConfigNode  // 配置树根节点
    Raw  any          // 原始数据（调试用）
}
```

**设计要点**:
- 统一不同格式的数据结构
- 支持嵌套对象和数组
- 提供类型判断方法: `IsObject()`, `IsArray()`, `IsPrimitive()`

---

### 3. yaml.go - YAML 解析器

**职责**: 解析 YAML 文件

```go
type YAMLParser struct {
    *BaseParser
}
```

**核心函数**:

```go
func buildConfigTree(key string, value interface{}) *ConfigNode {
    // 递归构建配置树
    // 1. 判断值的类型（map/slice/基本类型）
    // 2. 创建对应的 ConfigNode
    // 3. 递归处理子节点
}
```

**支持的扩展名**: `.yaml`, `.yml`

---

### 4. json.go - JSON 解析器

**职责**: 解析 JSON 文件

```go
type JSONParser struct {
    *BaseParser
}
```

**设计亮点**:
- JSON 数据结构与 YAML 相同（都是 `map[string]interface{}`）
- **复用** `yaml.go` 中的 `buildConfigTree()` 函数
- 只需实现 JSON 的读取逻辑

**支持的扩展名**: `.json`

**注意事项**:
- JSON 中所有数字都解析为 `float64`（JSON 规范限制）

---

### 5. factory.go - 解析器工厂

**职责**: 根据文件扩展名自动选择解析器

```go
type ParserFactory struct {
    mu      sync.RWMutex
    parsers map[string]Parser
}
```

#### 工作流程

```
ParseFile("config.yaml")
    ↓
1. 提取扩展名 ".yaml"
    ↓
2. 查找工厂: parsers["yaml"] → YAMLParser
    ↓
3. 调用: yamlParser.ParseFile()
    ↓
4. 返回: *ParseResult
```

#### 关键方法

```go
// 注册解析器
func Register(parser Parser)

// 根据扩展名获取
func GetParser(ext string) (Parser, error)

// 根据文件路径获取
func GetParserByFilePath(path string) (Parser, error)

// 便捷方法：直接解析
func ParseFile(path string) (*ParseResult, error)
```

**设计优势**:
- **开闭原则**: 添加新格式只需实现接口并注册
- **线程安全**: `sync.RWMutex` 保护
- **单例模式**: 全局工厂实例
- **扩展名规范化**: 自动处理 `.yaml`/`yaml`/`YAML` 等变体

---

## 🧪 测试

### 测试覆盖率

```bash
go test -cover github.com/cheryl-chun/confgen/internal/parser
```

**结果**: **96.0% 覆盖率** ✅

### 测试文件

| 文件 | 职责 | 测试数量 |
|------|------|---------|
| `types_test.go` | 测试 `ConfigNode` 数据结构 | 7 |
| `yaml_test.go` | 测试 YAML 解析器 | 8 |
| `json_test.go` | 测试 JSON 解析器 | 9 |
| `factory_test.go` | 测试工厂模式和集成 | 13 |

### 测试覆盖的场景

#### 1. 基本类型测试
- ✅ String, Int, Float, Bool, Null
- ✅ 类型推断准确性

#### 2. 复杂结构测试
- ✅ 嵌套对象（多层嵌套）
- ✅ 数组（基本类型数组、对象数组）
- ✅ 混合结构（对象中的数组、数组中的对象）

#### 3. 边界情况测试
- ✅ 空对象、空数组
- ✅ 无效格式（语法错误）
- ✅ 文件不存在
- ✅ 不支持的扩展名

#### 4. 工厂模式测试
- ✅ 扩展名规范化（`.yaml` vs `yaml` vs `YAML`）
- ✅ 解析器注册和查找
- ✅ 单例模式验证
- ✅ 线程安全（隐式测试）

#### 5. 集成测试
- ✅ 端到端文件解析
- ✅ 多格式文件处理
- ✅ 错误处理和传播

---

## 🎯 使用示例

### 基本用法

```go
import "github.com/cheryl-chun/confgen/internal/parser"

// 方式 1：直接解析文件（自动选择解析器）
result, err := parser.ParseFile("config.yaml")
if err != nil {
    log.Fatal(err)
}

// 访问配置树
root := result.Root
serverNode := root.Children["server"]
host := serverNode.Children["host"].Value.(string)
```

### 手动选择解析器

```go
// 方式 2：手动选择解析器
yamlParser := parser.NewYAMLParser()
result, err := yamlParser.ParseFile("config.yaml")
```

### 注册自定义解析器

```go
// 方式 3：扩展新格式
type TOMLParser struct {
    *parser.BaseParser
}

func NewTOMLParser() *TOMLParser {
    return &TOMLParser{
        BaseParser: parser.NewBaseParser("TOML", []string{".toml"}),
    }
}

func (p *TOMLParser) Parse(reader io.Reader) (*parser.ParseResult, error) {
    // 实现 TOML 解析逻辑
}

// 注册到全局工厂
parser.Register(NewTOMLParser())
```

---

## 📊 数据流

```
配置文件                    ConfigNode 树                类型推断器
(YAML/JSON)                (中间表示)                   (analyzer)
    ↓                           ↓                            ↓
┌─────────┐               ┌─────────┐                 ┌─────────┐
│ server: │   Parser      │  Root   │   Analyzer      │ Config  │
│   host  │  ========>    │  └─ Server  ============>  │  Server │
│   port  │   Parse()     │     ├─ Host│  Infer()     │    Host │
└─────────┘               │     └─ Port│              └─────────┘
                         └─────────┘                   代码生成
                         (本包输出)                    (下一步)
```

---

## ✅ 设计优势

### 1. **可扩展性**
- 添加新格式（TOML、INI、Properties）无需修改现有代码
- 只需实现 `Parser` 接口并注册

### 2. **代码复用**
- `BaseParser` 减少重复代码
- `buildConfigTree()` 在 YAML 和 JSON 之间复用

### 3. **统一抽象**
- 不同格式统一转换为 `ConfigNode` 树
- 后续处理无需关心原始格式

### 4. **类型安全**
- 明确的类型枚举（`ValueType`）
- 类型判断方法（`IsObject()` 等）

### 5. **易测试性**
- 接口设计利于 mock
- 纯函数便于单元测试
- 96% 测试覆盖率

---

## 🚀 未来扩展

### 计划支持的格式
- [ ] TOML (`.toml`)
- [ ] INI (`.ini`)
- [ ] Properties (`.properties`)
- [ ] HCL (`.hcl`) - Terraform 配置语言
- [ ] XML (`.xml`)

### 潜在改进
- [ ] 更详细的错误信息（行号、列号）
- [ ] 支持注释保留（用于生成带注释的代码）
- [ ] 流式解析（大文件优化）
- [ ] Schema 验证（可选）

---

## 🤝 贡献指南

添加新的解析器时，请确保：

1. ✅ 实现 `Parser` 接口的所有方法
2. ✅ 复用 `BaseParser`（如果适用）
3. ✅ 在 `factory.go` 的 `RegisterDefaultParsers()` 中注册
4. ✅ 编写完整的单元测试
5. ✅ 更新本文档

测试要求：
- 覆盖率 > 90%
- 包含边界情况测试
- 包含错误处理测试

---

**维护者**: confgen 团队  
**最后更新**: 2026-04-03  
**版本**: v0.1.0
