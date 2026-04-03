# confgen 项目文档

本目录包含 confgen 项目的详细设计文档和测试报告。

## 📚 文档列表

### Parser 包文档

1. **[PARSER_ARCHITECTURE.md](PARSER_ARCHITECTURE.md)** 📐
   - Parser 包的整体架构图解
   - 设计模式详解（工厂、策略、组合、单例）
   - 数据流图和类型系统
   - 扩展指南

2. **[PARSER_DESIGN.md](PARSER_DESIGN.md)** 📝
   - 每个文件的详细说明
   - 核心结构体和接口设计
   - 使用示例和最佳实践
   - 未来扩展计划

3. **[PARSER_TEST_REPORT.md](PARSER_TEST_REPORT.md)** ✅
   - 测试覆盖率报告（96.0%）
   - 37 个单元测试详情
   - 测试质量评估
   - 运行命令和性能指标

---

## 🎯 快速导航

### 我想了解...

#### 架构设计
→ 阅读 [PARSER_ARCHITECTURE.md](PARSER_ARCHITECTURE.md)  
查看完整的架构图、设计模式和数据流

#### 具体实现
→ 阅读 [PARSER_DESIGN.md](PARSER_DESIGN.md)  
了解每个文件的职责和实现细节

#### 测试质量
→ 阅读 [PARSER_TEST_REPORT.md](PARSER_TEST_REPORT.md)  
查看测试覆盖率和质量指标

#### 如何扩展
→ 参考 [PARSER_DESIGN.md#扩展指南](PARSER_DESIGN.md#🤝-贡献指南)  
学习如何添加新的解析器

---

## 📊 项目概览

### 当前状态

| 模块 | 状态 | 覆盖率 | 文档 |
|------|------|--------|------|
| Parser | ✅ 完成 | 96.0% | ✅ |
| Analyzer | ⚠️ 待开发 | - | - |
| CodeGen | ⚠️ 待开发 | - | - |
| Runtime | ⚠️ 待开发 | - | - |

### Parser 包统计

```
源代码:     ~413 行
测试代码:   ~1000 行
测试用例:   37 个
测试覆盖率: 96.0%
设计模式:   4 种
支持格式:   2 种 (YAML, JSON)
```

---

## 🏗️ Parser 包架构速览

```
Parser Package (已完成 ✅)
│
├── 核心接口
│   ├── Parser Interface        (统一解析接口)
│   └── BaseParser             (基类复用)
│
├── 解析器实现
│   ├── YAMLParser             (支持 .yaml, .yml)
│   └── JSONParser             (支持 .json)
│
├── 工厂模式
│   └── ParserFactory          (自动选择解析器)
│
└── 数据结构
    ├── ConfigNode             (配置树节点)
    ├── ValueType              (类型枚举)
    └── ParseResult            (解析结果)
```

---

## 🔄 数据流概览

```
输入文件 → Parser → ConfigNode 树 → Analyzer → TypeInfo → CodeGen → Go 代码
(YAML/JSON)  (已完成)  (中间表示)     (待开发)   (待开发)   (待开发)  (输出)
```

---

## 🧪 如何运行测试

```bash
# 进入项目根目录
cd d:\Personal\generateConfig

# 运行 Parser 包测试
go test ./internal/parser -v

# 查看覆盖率
go test ./internal/parser -cover

# 生成详细覆盖率报告
go test ./internal/parser -coverprofile=coverage.out
go tool cover -func=coverage.out

# 生成 HTML 覆盖率报告
go tool cover -html=coverage.out -o docs/coverage.html
```

---

## 📖 阅读顺序建议

### 对于新加入的开发者

1. 先读 **README.md**（项目根目录）- 了解项目目标
2. 再读 **PARSER_ARCHITECTURE.md** - 理解整体架构
3. 然后读 **PARSER_DESIGN.md** - 深入实现细节
4. 最后读 **PARSER_TEST_REPORT.md** - 学习测试方法

### 对于贡献者

1. 先读 **PARSER_DESIGN.md#贡献指南** - 了解规范
2. 再读 **PARSER_ARCHITECTURE.md#扩展接口** - 学习扩展方法
3. 参考现有测试文件编写测试
4. 提交前确保覆盖率 > 90%

---

## 🎓 设计模式学习

Parser 包是学习以下设计模式的优秀示例：

| 模式 | 位置 | 作用 |
|------|------|------|
| **策略模式** | Parser 接口 | 可切换的解析策略 |
| **工厂模式** | ParserFactory | 根据扩展名创建解析器 |
| **组合模式** | BaseParser | 复用通用功能 |
| **单例模式** | GetFactory() | 全局唯一工厂实例 |

详细说明见 [PARSER_ARCHITECTURE.md#设计模式详解](PARSER_ARCHITECTURE.md#🎭-设计模式详解)

---

## 🚀 下一步开发

根据项目路线图，接下来需要开发：

### 1. Analyzer 包（类型推断）
- 从 ConfigNode 树推断 Go 类型
- 处理数组统一类型
- 处理嵌套结构体命名

### 2. CodeGen 包（代码生成）
- 生成 Go 结构体定义
- 生成字段标签（json, yaml, mapstructure）
- 代码格式化

### 3. Runtime 包（运行时库）
- Trie 数据结构
- 多源配置加载
- 热重载支持

---

## 📝 文档维护

### 更新频率
- 重大功能变更：立即更新
- 小改进：每周汇总更新
- 测试报告：每次测试后更新

### 文档规范
- 使用 Markdown 格式
- 包含图表和示例代码
- 保持简洁清晰
- 及时更新版本号

---

## 🤝 贡献指南

提交代码时，请确保：

1. ✅ 代码通过所有测试
2. ✅ 新功能包含单元测试（覆盖率 > 90%）
3. ✅ 更新相关文档
4. ✅ 遵循代码风格（gofmt）
5. ✅ 提交信息清晰

---

## 📞 联系方式

- 项目主页: https://github.com/cheryl-chun/confgen
- Issue 跟踪: https://github.com/cheryl-chun/confgen/issues
- 讨论区: https://github.com/cheryl-chun/confgen/discussions

---

**文档版本**: v1.0  
**最后更新**: 2026-04-03  
**维护团队**: confgen 开发团队
