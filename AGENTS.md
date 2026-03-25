# AGENTS.md

## 项目概述
这是一个 Windows 现代化简洁桌面客户端项目。
- 后端使用 Go 语言处理业务逻辑。
- 前端使用 Wails 作为 UI 框架（Wails + HTML/CSS/JS）。
- 用Makefile实现各种make指令实现打包测试编译运行
- 数据存储采用 JSON 文件。
- 目标是简洁 UI、跨平台可扩展、运行效率高。

## 技术栈
- Go 语言（后端逻辑）
- Wails 框架（前端 UI + 桥接 Go）
- HTML / CSS / JS 模板（如 TailwindCSS / React）
- JSON 文件用于轻量配置和缓存

## 项目结构约定
- /cmd/ 主程序
- /internal/ 后端业务逻辑
- /frontend/ 前端 Wails 项目
- /pkg/db/ 数据库访问层
- /config/ 配置/数据存储

## 构建与运行命令
- 安装依赖：`go mod tidy`
- 启动开发：`wails dev`
- 编译发布：`wails build`
- 运行主程序：`./your-app.exe`

## 编码规范
- Go 代码使用 `gofmt` 格式化
- 所有前端 JS 使用现代语法（ES6+）
- UI 组件尽量可复用、模块化

## 数据访问约定
- JSON 用于轻量配置，例如主题、用户设置
- SQLite 用于列表或结构化业务数据
- 所有数据库访问通过统一包 `pkg/db`

## 测试 & 本地调试
- Go 单元测试：`go test ./...`
- 前端 UI 调试：`wails dev`