# CLAUDE.md (frontend)

@AGENTS.md

## Claude/Codex 前端补充
- 默认只读取 `frontend/` 目录内容，按需最小化读取其它目录。
- 优先改动 `src/modules`、`src/components`、`src/services`，减少跨层耦合。
- 若依赖后端新接口，先在 `src/types` 与 `src/services` 定义调用契约。
- 保持样式与交互一致性，避免临时内联样式散落。
