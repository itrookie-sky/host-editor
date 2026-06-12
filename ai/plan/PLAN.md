# Host Editor Framework Implementation Plan

## Summary
- 当前模型：Codex（GPT-5）。已加载全局规则 `core.md`、`workflow.md`、项目 `AGENTS.md`，并按任务使用 `writing-plans`、`codegraph`、`frontend-design`、`goframe-v2`。
- 目标：把当前 Wails 默认模板改成可用的 hosts 多版本编辑器，前端使用 Vue 3 + Monaco，文件统一存储到 `$HOME/.host-editor`。
- 默认范围：不读写系统 `/etc/hosts`，不引入路由、状态管理库或 UI 组件库。

## Key Changes
- Go/Wails 后端：
  - 在 `internal/view` 增加 hosts 文件存储能力，暴露 Wails 方法：
    - `ListHostFiles() ([]HostFileInfo, error)`
    - `ReadHostFile(name string) (string, error)`
    - `SaveHostFile(req SaveHostFileRequest) error`
    - `CreateHostFile(name string) (HostFileInfo, error)`
  - 只接受文件名，不接受前端传入路径；校验空名、路径穿越、非法分隔符。
  - 启动或首次列举时确保 `$HOME/.host-editor` 存在；目录为空时创建默认 `hosts` 文件。
  - 保存使用同目录临时文件 + rename，降低写入中断风险。
- 桌面窗口：
  - 修改 `main.go` 的 Wails 配置，设置合理 `MinWidth/MinHeight`。
  - macOS 使用 `mac.TitleBarHiddenInset()`，让内容铺满窗口并保留系统窗口控制按钮。
  - 前端顶栏设置 Wails draggable 区域，避免隐藏 titlebar 后窗口不可拖动。
- Vue/Monaco 前端：
  - 在 `frontend/package.json` 增加 `monaco-editor`。
  - 删除默认 `HelloWorld` 使用链路，用 `App.vue` 组织左侧文件栏、顶部操作栏、编辑区。
  - 新增组件：
    - `HostSidebar.vue`：文件列表、当前选择、新建版本入口。
    - `HostsEditor.vue`：Monaco 初始化、内容变更、只读/加载状态。
    - `types/host.ts`：前后端共享的 TS 类型。
  - `App.vue` 负责加载文件列表、读取当前文件、保存内容、新建版本、错误提示和 dirty 状态。
- UI 方向：
  - 工具型桌面应用，不做 landing page。
  - 采用紧凑三栏/两栏工作台：左侧版本列表，右侧 Monaco 编辑器，顶部窄操作栏。
  - 色彩保持克制：深灰工作台、浅色文本、绿色 active 状态、琥珀色 dirty 状态，避免默认 Wails logo 和演示文案。

## Test Plan
- Go 单元测试：
  - `validateHostFileName`：合法文件名、空名、`../hosts`、包含路径分隔符。
  - store 初始化：目录不存在时创建目录；空目录时创建默认 `hosts`。
  - read/save：保存后读取内容一致；非法文件名返回错误。
- 构建验证：
  - `go test ./...`
  - `cd frontend && npm run build`
  - `wails build`
- 手工验收：
  - `wails dev` 启动后窗口无 titlebar 白边，内容铺满。
  - 左侧能看到默认 `hosts`。
  - 新建一个版本文件后能切换、编辑、保存。
  - 关闭重开后 `$HOME/.host-editor` 中的内容仍存在。

## Assumptions
- hosts 多版本是应用私有文件版本，不直接同步系统 `/etc/hosts`。
- 默认文件名为 `hosts`；新建文件名必须唯一。
- Monaco 只按纯文本编辑 hosts 内容，不做语法高亮定制。
- 当前 `ai/` 目录为未跟踪目录，计划文件写入时只新增计划文档，不改需求文件。
