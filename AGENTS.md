# AGENT.md

本项目是一个基于 Wails 的桌面应用，后端使用 Go，前端使用 Vue 3。

## 项目领域

- 应用名称：`host-editor`
- 应用类型：桌面应用
- 框架形态：Wails v2 将 Go 后端能力绑定到前端，并通过内置 WebView 承载 Vue 页面
- 主要入口：
  - Go 启动入口：`main.go`
  - Go 应用绑定：`app.go`
  - 前端入口：`frontend/src/main.ts`
  - 前端根组件：`frontend/src/App.vue`

## 技术栈

| 层级     | 技术       |
| -------- | ---------- |
| 桌面框架 | Wails v2   |
| 后端语言 | Go 1.26.3  |
| 前端框架 | Vue 3      |
| 前端构建 | Vite       |
| 前端语言 | TypeScript |
| 包管理   | npm        |

## 目录约定

- `main.go`：配置 Wails 应用窗口、资源嵌入、生命周期和 Go 绑定对象。
- `app.go`：放置可暴露给前端调用的 Go 应用方法和启动生命周期逻辑。
- `frontend/`：Vue 3 前端工程。
- `frontend/src/`：前端源码目录。
- `frontend/dist/`：前端构建产物，由 Wails 通过 `embed.FS` 嵌入，不应手工维护。
- `build/`：Wails 构建资源和产物目录，按 Wails 约定维护。

## 常用命令

在项目根目录执行：

```bash
wails dev
```

用于启动 Wails 开发模式。

```bash
wails build
```

用于构建桌面应用。

在 `frontend/` 目录执行：

```bash
npm run dev
```

用于单独启动 Vite 前端开发服务。

```bash
npm run build
```

用于执行 `vue-tsc --noEmit` 类型检查并构建前端资源。

## 开发约定

- Go 侧新增可供前端调用的方法时，优先放在 `app.go` 或按职责拆分到新的 Go 文件，并通过 Wails `Bind` 暴露。
- 前端调用 Go 方法时，优先使用 Wails 生成的绑定代码，不手写桥接路径。
- 修改窗口、资源嵌入、生命周期回调时，优先检查 `main.go` 和 `wails.json`。
- 修改前端构建、开发服务或安装命令时，优先检查 `frontend/package.json` 和 `wails.json` 中的 `frontend:*` 配置。
- 保持 Go 后端负责系统能力、文件读写和桌面生命周期，Vue 前端负责界面状态和交互展示。
- 不直接修改 Wails 或前端构建生成文件，必要时通过对应源码或配置重新生成。

## 验证建议

- Go 或 Wails 侧改动后，至少执行 `wails build` 或对应的最小 Go 编译验证。
- 前端源码改动后，至少执行 `npm run build`。
- 只修改文档时，可使用 `git diff --check` 检查空白和格式问题。
