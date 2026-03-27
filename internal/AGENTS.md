# AGENTS.md (internal)

## 作用域
仅针对 `internal/` 后端代码生效，目标是减少无关上下文读取。

## 技术与约束
- 语言：Go
- 角色：业务逻辑与数据处理层
- 前端桥接：通过 Wails 绑定暴露方法，不在这里写 UI 代码

## 目录职责
- `app/`：应用装配与可绑定入口对象
- `model/`：领域模型与前端展示模型
- `repository/`：数据访问（JSON/SQLite 等）
- `service/`：业务服务编排
- `usecase/`：用例层（跨服务流程）
- `transport/wails/`：Wails 绑定与启动装配
- `config/`：后端运行期配置结构

## 编码规则
- 逻辑代码放在 `internal/`，公共能力放在 `pkg/`
- 所有数据访问走 `repository` 抽象
- 返回给前端的数据在 `service` 或 `usecase` 层组装
- 提交前执行 `gofmt`；测试优先执行 `go test ./...`

## 非目标
- 不修改 `frontend/` 的 UI 组件和样式
- 不引入 Electron 或其他大型桌面框架
