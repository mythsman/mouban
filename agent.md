# Agent 工作指南（mouban）

## 1) 项目概览
- 项目：`mouban`（Go 1.21）
- 作用：为 `hexo-douban` 提供豆瓣用户书/影/音/游数据抓取与查询服务
- 入口：`main.go`
- 核心目录：
  - `internal/app/`：启动编排、HTTP 路由与中间件
  - `internal/controller/`：HTTP 接口处理（`/guest/*`、`/admin/*`）
  - `internal/crawl/`：抓取逻辑
  - `internal/dao/`：数据库读写
  - `internal/model/`：数据模型
  - `internal/common/`：配置、日志、数据库初始化
  - `internal/agent/`：后台任务与调度流程

## 2) 配置与运行
- 配置文件：`.env`（本地开发可由 `.env.sample` 复制）
- 环境变量规则：统一使用 `MOUBAN_` 前缀、全大写、单下划线
  - 例如：`datasource.host` -> `MOUBAN_DATASOURCE_HOST`
- 启动由 `app.Bootstrap()` 显式执行（不再依赖包 `init()` 自动连接数据库）

### 本地启动
```bash
go run .
```

## 3) 开发约束（必须遵守）
### Do
- 仅做与需求相关的最小改动，保持现有接口行为兼容
- 修改后执行至少与改动相关的测试/验证
- 保持 Go 代码可通过 `gofmt`
- 新增配置时，同步更新：
  - `.env.sample`
  - 相关 README 说明（如影响使用方式）

### Don’t
- 不要在未确认的情况下修改公开 API 路径和返回结构
- 不要提交真实密钥、账号、cookie、代理凭据
- 不要直接运行 `build.sh`（会执行 docker build/push）
- 不要大规模重构无关模块

## 4) 验证策略
按“从小到大”验证，避免一次性全量测试导致噪音：

1. 格式化
```bash
gofmt -w <changed_files>
```

2. 单包测试（优先）
```bash
go test ./internal/util -run Test
```

3. 改动包测试
```bash
go test ./internal/<changed_pkg>
```

4. 全量测试（仅在具备依赖条件时）
```bash
go test ./...
```

> 说明：仓库内部分测试依赖外部网络或 MySQL；若环境未就绪，优先执行“受影响包 + 无外部依赖”的测试。

## 5) 常见任务操作流程
### A. 新增/修改接口
1. 在 `internal/controller/` 增加处理逻辑
2. 必要时补充 `internal/dao/` 与 `internal/model/`
3. 在 `internal/app/http.go` 路由组注册
4. 补充最小测试或调用示例

### B. 修改抓取逻辑
1. 在 `internal/crawl/` 定位目标站点解析逻辑
2. 保持容错（网络失败、字段缺失、重试）
3. 校验入库链路（`internal/dao/`）是否受影响

### C. 修改数据库相关
1. 更新 `internal/model/` 结构
2. 确认 `AutoMigrate` 兼容性
3. 检查查询/写入逻辑和索引影响

## 6) 提交前检查清单（Definition of Done）
- [ ] 改动范围最小且聚焦
- [ ] 代码已格式化（gofmt）
- [ ] 已执行相关测试并记录结果
- [ ] 未引入敏感信息
- [ ] 需要时更新 `README.md` / `.env.sample`
- [ ] 变更说明包含：改了什么、为什么、风险点

## 7) 建议的提交信息模板
```text
<type>: <summary>

- change 1
- change 2
- test: <commands>
```

`type` 可用：`feat` / `fix` / `refactor` / `docs` / `test` / `chore`
