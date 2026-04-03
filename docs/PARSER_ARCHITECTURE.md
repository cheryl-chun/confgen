# Parser 包架构图解

## 🏗️ 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                       Parser Package                         │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Parser Interface (策略模式)              │  │
│  │  - Parse(reader) → ParseResult                       │  │
│  │  - ParseFile(path) → ParseResult                     │  │
│  │  - SupportedExtensions() → []string                  │  │
│  │  - Name() → string                                   │  │
│  └────────────┬─────────────────────────────────────────┘  │
│               │ implements                                  │
│       ┌───────┴────────┬──────────────┬──────────────┐     │
│       │                │              │              │     │
│  ┌────▼────┐     ┌────▼────┐   ┌────▼────┐   ┌─────▼──┐  │
│  │  Base   │     │  YAML   │   │  JSON   │   │ Future │  │
│  │ Parser  │◄────┤ Parser  │   │ Parser  │   │ Parsers│  │
│  └─────────┘     └─────────┘   └─────────┘   └────────┘  │
│  (基类复用)      (已实现)       (已实现)      (TOML等)    │
│                       │              │                      │
│                       └──────┬───────┘                      │
│                              │ 复用                         │
│                    ┌─────────▼─────────┐                    │
│                    │  buildConfigTree  │                    │
│                    │  (树构建函数)      │                    │
│                    └─────────┬─────────┘                    │
│                              │ 生成                         │
│                    ┌─────────▼─────────┐                    │
│                    │   ConfigNode 树    │                    │
│                    │  (中间表示 IR)     │                    │
│                    └───────────────────┘                    │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         ParserFactory (工厂模式 + 单例)              │  │
│  │                                                       │  │
│  │  parsers map[string]Parser                           │  │
│  │  ├─ "yaml"  → YAMLParser                             │  │
│  │  ├─ "yml"   → YAMLParser                             │  │
│  │  └─ "json"  → JSONParser                             │  │
│  │                                                       │  │
│  │  Methods:                                             │  │
│  │  - Register(parser)              [注册解析器]        │  │
│  │  - GetParser(ext)                [根据扩展名获取]    │  │
│  │  - GetParserByFilePath(path)     [根据路径获取]      │  │
│  │  - ParseFile(path)               [一键解析]          │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## 📊 数据流图

### 1. 解析流程

```
用户输入                解析器选择              文件解析              树构建
   │                       │                      │                     │
   ├─ config.yaml ──────►  Factory  ──────►  YAMLParser  ──────►  ConfigNode
   │                       │                      │                     │
   │                  GetParser(".yaml")     Parse(file)         buildTree()
   │                       │                      │                     │
   ▼                       ▼                      ▼                     ▼
返回结果           选择 YAMLParser         读取 YAML           构建树结构
ParseResult              实例               解析为 map          递归转换
```

### 2. 树构建过程

```
YAML 输入:                         ConfigNode 树:
──────────                         ───────────────

server:                            Root (TypeObject)
  host: "localhost"                 │
  port: 8080                        └─ server (TypeObject)
  features:                            ├─ host (TypeString)
    - ssl                              │   Value: "localhost"
    - cache                            │
                                       ├─ port (TypeInt)
                                       │   Value: 8080
                                       │
                                       └─ features (TypeArray)
                                           ├─ [0] (TypeString)
                                           │   Value: "ssl"
                                           └─ [1] (TypeString)
                                               Value: "cache"

递归过程:
1. buildConfigTree("server", map[...])
   → 创建 TypeObject 节点
   → 遍历 map 的每个 key-value
2. buildConfigTree("host", "localhost")
   → 创建 TypeString 节点
3. buildConfigTree("port", 8080)
   → 创建 TypeInt 节点
4. buildConfigTree("features", []interface{...})
   → 创建 TypeArray 节点
   → 遍历数组的每个元素
5. buildConfigTree("[0]", "ssl")
   → 创建 TypeString 节点
```

## 🎭 设计模式详解

### 1. 策略模式 (Strategy Pattern)

```
┌───────────────────────────────────────┐
│          Context (使用者)              │
│                                        │
│  parser := factory.GetParser("yaml")  │
│  result := parser.Parse(file)         │
└────────────────┬──────────────────────┘
                 │
                 │ 调用统一接口
                 ▼
┌────────────────────────────────────────┐
│        Strategy (Parser 接口)          │
│  + Parse()                             │
│  + ParseFile()                         │
└────────────────┬───────────────────────┘
                 │
        ┌────────┼────────┐
        │                 │
   ┌────▼─────┐     ┌────▼─────┐
   │  YAML    │     │  JSON    │
   │ Strategy │     │ Strategy │
   └──────────┘     └──────────┘
   
优势:
✅ 添加新格式无需修改现有代码
✅ 运行时可切换解析策略
✅ 每个策略独立测试
```

### 2. 工厂模式 (Factory Pattern)

```
┌────────────────────────────────────────┐
│          Client (调用者)                │
│                                         │
│  ParseFile("config.yaml")              │
└────────────────┬───────────────────────┘
                 │
                 │ 1. 请求解析
                 ▼
┌────────────────────────────────────────┐
│       Factory (ParserFactory)          │
│                                         │
│  1. 提取扩展名: ".yaml"                 │
│  2. 查找: parsers["yaml"]               │
│  3. 返回: YAMLParser 实例               │
└────────────────┬───────────────────────┘
                 │
                 │ 2. 返回合适的解析器
                 ▼
┌────────────────────────────────────────┐
│       Product (具体解析器)              │
│                                         │
│  YAMLParser.Parse(file)                │
└─────────────────────────────────────────┘

优势:
✅ 客户端无需知道具体类
✅ 集中管理对象创建
✅ 易于扩展新产品
```

### 3. 组合模式 (Composite Pattern)

```
BaseParser (基类)
├─ name: string
├─ extensions: []string
├─ Name() → string
└─ SupportedExtensions() → []string
   │
   │ 组合
   ▼
YAMLParser
├─ *BaseParser          ← 复用基类
├─ Parse()              ← 实现特定逻辑
└─ ParseFile()          ← 实现特定逻辑

优势:
✅ 复用通用代码
✅ 减少重复
✅ 保持接口统一
```

### 4. 单例模式 (Singleton Pattern)

```go
var (
    globalFactory     *ParserFactory
    globalFactoryOnce sync.Once
)

func GetFactory() *ParserFactory {
    globalFactoryOnce.Do(func() {
        globalFactory = NewParserFactory()
        globalFactory.RegisterDefaultParsers()
    })
    return globalFactory
}
```

```
第一次调用                    后续调用
GetFactory()                 GetFactory()
    │                            │
    ├─ 创建实例                  │
    ├─ 注册解析器                │
    └─ 返回 factory              └─ 直接返回已有实例
       │                            │
       └──────────┬─────────────────┘
                  │
                  ▼
            同一个实例
            (线程安全)

优势:
✅ 全局唯一实例
✅ 延迟初始化
✅ 线程安全 (sync.Once)
```

## 🔄 类型系统

### ValueType 类型树

```
ValueType (类型枚举)
│
├─ Primitive (基本类型)
│  ├─ TypeString   → Go: string
│  ├─ TypeInt      → Go: int
│  ├─ TypeFloat    → Go: float64
│  ├─ TypeBool     → Go: bool
│  └─ TypeNull     → Go: nil
│
└─ Complex (复合类型)
   ├─ TypeObject   → Go: struct
   │  └─ Children: map[string]*ConfigNode
   │
   └─ TypeArray    → Go: []T
      └─ Items: []*ConfigNode
```

### ConfigNode 结构

```
ConfigNode
├─ Key: string              ← 字段名
├─ Value: any               ← 原始值
├─ Type: ValueType          ← 类型标记
│
├─ Children: map[string]*ConfigNode    ← 对象的子节点
│  └─ key → child node
│
└─ Items: []*ConfigNode     ← 数组的元素
   └─ [0] → item node
```

## 🔗 扩展接口

### 添加新解析器的步骤

```
1. 定义结构体
   ┌─────────────────────┐
   │ type TOMLParser {   │
   │   *BaseParser       │
   │ }                   │
   └─────────────────────┘
            │
            ▼
2. 实现 Parser 接口
   ┌─────────────────────┐
   │ Parse(reader)       │
   │ ParseFile(path)     │
   └─────────────────────┘
            │
            ▼
3. 注册到工厂
   ┌─────────────────────┐
   │ factory.Register(   │
   │   NewTOMLParser()   │
   │ )                   │
   └─────────────────────┘
            │
            ▼
4. 编写测试
   ┌─────────────────────┐
   │ TestTOMLParser_*    │
   │ 覆盖率 > 90%        │
   └─────────────────────┘
```

## 📦 输出格式

### ParseResult 结构

```
ParseResult
├─ Root: *ConfigNode        ← 供类型推断使用
│  │
│  └─ 完整的配置树
│     (标准化的中间表示)
│
└─ Raw: any                 ← 供调试使用
   │
   └─ 原始解析结果
      (保留原始格式)
```

### 后续流程

```
Parser 输出               Analyzer 输入            Codegen 输入
    │                         │                         │
    ▼                         ▼                         ▼
ParseResult              TypeInfo              Go Code String
  └─ Root                 └─ Structs            └─ type Config {...}
     ConfigNode              └─ Fields             └─ type Server {...}
```

## 🎯 关键路径

### 核心执行路径

```
用户调用
  │
  ▼
parser.ParseFile("config.yaml")
  │
  ├─ 1. GetFactory() → 获取/创建工厂实例
  │     │
  │     └─ globalFactoryOnce.Do(初始化)
  │
  ├─ 2. GetParserByFilePath("config.yaml")
  │     │
  │     ├─ 提取扩展名: ".yaml"
  │     └─ 查找: parsers["yaml"] → YAMLParser
  │
  ├─ 3. yamlParser.ParseFile("config.yaml")
  │     │
  │     ├─ 打开文件
  │     ├─ yaml.NewDecoder(file)
  │     └─ decoder.Decode(&raw)
  │
  ├─ 4. buildConfigTree("root", raw)
  │     │
  │     └─ 递归构建 ConfigNode 树
  │         ├─ 判断类型 (map/slice/primitive)
  │         ├─ 创建节点
  │         └─ 处理子节点
  │
  └─ 5. 返回 ParseResult
        │
        ├─ Root: *ConfigNode
        └─ Raw: interface{}
```

## 📏 代码度量

```
代码行数分布:
├─ parser.go      51 行  (接口定义)
├─ types.go       75 行  (数据结构)
├─ yaml.go       102 行  (YAML 实现)
├─ json.go        49 行  (JSON 实现)
└─ factory.go    136 行  (工厂模式)
   ──────────
   总计: ~413 行

测试代码:
├─ types_test.go     ~180 行
├─ yaml_test.go      ~240 行
├─ json_test.go      ~250 行
└─ factory_test.go   ~330 行
   ──────────
   总计: ~1000 行

测试/代码比: 2.4:1 (优秀)
```

---

**文档版本**: v1.0  
**最后更新**: 2026-04-03  
**架构师**: confgen 团队
