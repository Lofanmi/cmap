# 贡献指南

感谢您对 CMap 项目的关注！我们欢迎所有形式的贡献，包括但不限于：

- 🐛 Bug 报告
- 💡 功能建议
- 📝 文档改进
- 🔧 代码贡献
- 🧪 测试用例

## 开发环境设置

### 前置要求

- Go 1.21 或更高版本
- Git

### 本地开发

1. Fork 本仓库
2. 克隆您的 fork：
   ```bash
   git clone https://github.com/YOUR_USERNAME/cmap.git
   cd cmap
   ```

3. 添加上游仓库：
   ```bash
   git remote add upstream https://github.com/Lofanmi/cmap.git
   ```

4. 安装依赖：
   ```bash
   go mod download
   ```

## 代码规范

### Go 代码风格

- 遵循 [Effective Go](https://golang.org/doc/effective_go.html) 规范
- 使用 `gofmt` 格式化代码
- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

类型说明：
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

示例：
```
feat: add new concurrent test cases

- Add TestConcurrentStress for high-load testing
- Add TestConcurrentReadWrite for read-write mixed scenarios
- Improve test coverage to 50%

Closes #123
```

## 测试

### 运行测试

```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestNewFunction

# 运行并发测试
go test -v -run TestConcurrent

# 运行基准测试
go test -bench=. -benchmem
```

### 测试覆盖率

```bash
# 生成覆盖率报告
go test -coverprofile=coverage.out

# 查看覆盖率报告
go tool cover -html=coverage.out
```

### 代码质量检查

```bash
# 运行 golangci-lint
golangci-lint run

# 运行 govulncheck
govulncheck ./...
```

## 提交 Pull Request

1. 创建功能分支：
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. 进行更改并提交：
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

3. 推送到您的 fork：
   ```bash
   git push origin feature/your-feature-name
   ```

4. 创建 Pull Request

### Pull Request 检查清单

- [ ] 代码通过所有测试
- [ ] 添加了相应的测试用例
- [ ] 更新了相关文档
- [ ] 遵循了代码规范
- [ ] 提交信息符合规范
- [ ] 没有引入新的警告或错误

## 报告 Bug

### Bug 报告模板

```markdown
**描述**
简要描述 Bug

**重现步骤**
1. 
2. 
3. 

**预期行为**
描述您期望看到的行为

**实际行为**
描述实际发生的行为

**环境信息**
- 操作系统：
- Go 版本：
- CMap 版本：

**附加信息**
任何其他相关信息，如错误日志、截图等
```

## 功能建议

### 功能建议模板

```markdown
**功能描述**
简要描述您希望添加的功能

**使用场景**
描述该功能的使用场景和好处

**实现建议**
如果有的话，提供实现建议

**替代方案**
如果有的话，描述替代方案
```

## 发布流程

### 版本发布检查清单

- [ ] 所有测试通过
- [ ] 文档已更新
- [ ] CHANGELOG.md 已更新
- [ ] 版本号已更新
- [ ] 标签已创建

### 创建发布

1. 更新版本号
2. 更新 CHANGELOG.md
3. 创建 Git 标签
4. 推送到 GitHub
5. 创建 GitHub Release

## 联系方式

如果您有任何问题或建议，请通过以下方式联系我们：

- [GitHub Issues](https://github.com/Lofanmi/cmap/issues)
- [GitHub Discussions](https://github.com/Lofanmi/cmap/discussions)

## 行为准则

我们致力于为每个人提供友好、安全和欢迎的环境。请参阅我们的 [行为准则](CODE_OF_CONDUCT.md)。

## 许可证

通过贡献代码，您同意您的贡献将在 MIT 许可证下发布。 