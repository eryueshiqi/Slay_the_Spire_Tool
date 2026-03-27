# AGENTS.md (frontend)

## 作用域
仅针对 `frontend/` 前端代码生效，避免引入后端无关上下文。

## 技术与约束
- UI 框架：Wails 前端（HTML/CSS/JS）
- 组件建议：React + TailwindCSS（如需要组件化）
- 目标：简洁、模块化、可维护

## 目录职责
- `src/pages/`：页面级组件
- `src/components/`：可复用组件
- `src/modules/`：按功能拆分的业务模块
- `src/services/`：接口调用与数据适配
- `src/store/`：状态管理
- `src/styles/`：全局样式与主题
- `src/types/`：类型定义
- `src/assets/`：构建期静态资源
- `public/`：运行期静态资源与 mock

## 编码规则
- 前端只消费后端暴露的 Wails 方法，不实现后端逻辑
- 组件保持单一职责，避免页面内堆叠大块逻辑
- 统一在 `services/` 做数据转换，页面层只处理展示

## 非目标
- 不在 `frontend/` 写 Go 逻辑
- 不直接改 `internal/`，除非明确是全栈联动任务
