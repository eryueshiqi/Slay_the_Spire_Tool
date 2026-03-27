# CLAUDE.md (internal)

@AGENTS.md

## Claude/Codex 内部补充
- 默认只读取 `internal/` 与必要的 `pkg/`、`config/` 文件。
- 优先做小范围改动，避免扫描 `frontend/` 目录。
- 若后端接口变更，给出明确的前端调用签名（方法名、参数、返回结构）。
- 需要示例数据时，优先使用 `config/data/*.json`。
