# CLAUDE.md

@AGENTS.md

## Claude 专用补充说明

### AI 辅助上下文
- 优先遵循 AGENTS.md 中的技术栈和约束。
- 不要使用 Electron 或大型浏览器框架。
- 优先生成 Go + Wails 结构代码。
- 遇到 UI 代码需要遵循简洁、模块化、可维护的原则。

### 项目规则（Claude）
- 所有逻辑代码应在 `internal/` 下实现
- 所有前端 UI 代码放在 `frontend/`
- 不要引入额外前端框架除非必要
- 所有数据库操作应通过 `pkg/db` 封装

### 特殊指令
- 如果问到“如何交互 UI 与 Go”，请使用 Wails 绑定
- 如果需生成组件结构，优先 React + TailwindCSS