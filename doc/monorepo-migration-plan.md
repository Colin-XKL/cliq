# Monorepo 迁移规划（cliQ）

## 概述
- 目标：将现有代码库重构为 Nx + pnpm + go.work 管理的 monorepo，统一前端与 Go 模块的开发与发布流程。
- 原始结构：`cliq/`（Wails 应用，Go + Vue）、`cliq-hub-backend/`（纯 Go 后端）、`doc/` 文档与示例。
- 迁移原则：先保证“可构建、可运行”，再做“抽取复用”；尽量保留现有行为（Wails 前端命令、接口路径等），避免大规模同时改动。

## 当前仓库盘点（关键点）
- Wails 应用：
  - 配置：`cliq/wails.json:1`
  - 前端：`cliq/frontend`（Vue3 + Vite + pnpm）
  - 模板处理：`cliq/services/template_service.go:1`、`cliq/handlers/file_handler.go:125`
  - 设置服务：`cliq/config/settings_service.go:53`，前端交互 `cliq/frontend/src/composables/useSettings.ts:1`
  - 模型：`cliq/models/models.go:1`
- 纯 Go 后端：
  - HTTP 路由与生成接口：`cliq-hub-backend/internal/http/router/router.go`、`cliq-hub-backend/internal/http/handlers/generate_handler.go`
  - LLM 客户端：`cliq-hub-backend/internal/llm/client.go:37`
  - 模板模型：`cliq-hub-backend/internal/template/model.go:1`
  - YAML 编解码与清洗：`cliq-hub-backend/pkg/yaml/codec.go:1`
- 前端与后端耦合点：
  - 智能生成调用：`cliq/frontend/src/pages/TemplateGenerator.vue:99`
  - Base URL 设置与校验：`cliq/frontend/src/pages/SettingsPage.vue:44`、`cliq/frontend/src/composables/useSettings.ts:1`

## 目标 Monorepo 目录结构
```
/your-monorepo
├── apps/
│   ├── cliq-app/
│   │   ├── frontend/
│   │   │   ├── package.json     (pnpm 管理, 依赖 @repo/shared-vue-ui)
│   │   │   ├── vite.config.ts
│   │   │   └── ...
│   │   ├── main.go
│   │   ├── go.mod               (依赖 @repo/shared-go-lib)
│   │   ├── wails.json
│   │   └── project.json         (Nx：包装 Wails CLI)
│   ├── cliq-hub-backend/
│   │   ├── cmd/server/main.go
│   │   ├── internal/...         (保留原结构)
│   │   ├── go.mod               (依赖 @repo/shared-go-lib)
│   │   └── project.json         (Nx：使用 @nx-go/nx-go)
│   └── cliq-frontend/
│       ├── package.json         (依赖 @repo/shared-vue-ui)
│       ├── vite.config.ts
│       └── project.json         (Nx：使用 @nx/vue)
├── packages/
│   ├── shared-vue-ui/
│   │   ├── package.json         (@repo/shared-vue-ui)
│   │   ├── vite.config.ts       (Vite 库模式)
│   │   └── project.json         (Nx 配置)
│   └── shared-go-lib/
│       ├── lib.go               (统一模板模型、YAML工具等)
│       ├── go.mod               (@repo/shared-go-lib)
│       └── project.json         (Nx 配置)
├── go.work                      (关联所有 go.mod)
├── nx.json                      (Nx 全局配置)
├── package.json                 (根：pnpm, nx, @nx-go/nx-go)
├── pnpm-lock.yaml
└── pnpm-workspace.yaml          (工作区定义)
```

## 分期迁移计划

### 阶段 1：Monorepo 骨架与目录迁移
- 在根创建：`pnpm-workspace.yaml`（包含 `apps/*`、`packages/*`）、`package.json`（添加 `nx`、`@nx/vue`、`@nx-go/nx-go`）、`nx.json`、`go.work`。
- 迁移目录：
  - `cliq/` → `apps/cliq-app/`（保留 `frontend/`、`main.go`、`go.mod`、`wails.json`）
  - `cliq-hub-backend/` → `apps/cliq-hub-backend/`（保留原 `cmd/`、`internal/`、`pkg/`）
- 验证：根执行 `pnpm install`，`go work` 能解析两个 app 的 go 模块。

### 阶段 2：Vue 工作区与共享 UI
- `apps/cliq-app/frontend` 使用根安装，去除局部锁文件。
- 仅抽取表单生成逻辑：在 `packages/shared-vue-ui` 中抽取 `DynamicCommandForm.vue`（来源 `cliq/frontend/src/components/DynamicCommandForm.vue`），提供基于模板变量的表单渲染能力，供 `wails-app` 与 `vue-frontend` 复用。
- 暂不包含复杂组件：不迁移 `TemplateEditorModal.vue:227` 等复杂交互组件，后续逐步迁移。
- 在 `apps/cliq-app/frontend/package.json` 与 `apps/cliq-frontend` 引入 `@repo/shared-vue-ui`。
- 创建 `apps/cliq-frontend` 的最小 Hello World 项目（Vite + Vue）用于验证 workspace：
  - `App.vue` 显示 “Hello World” 文本；
  - 不包含 `wailsjs` 绑定，后续通过 `fetch` 指向后端。

### 阶段 3：Go 工作区与共享库
- 创建 `packages/shared-go-lib`，抽取模板相关全部核心能力：
  - 模板定义：统一 `Template`/`Command`/`Variable` 的结构体，合并 `cliq/models/models.go:1` 与 `cliq-hub-backend/internal/template/model.go:1`（仅保留一份定义）。
  - 模板解析：迁移并导出命令解析/模板生成/模板反序列化方法，来源 `cliq/services/template_service.go:20`、`:84`、`:127`；同时保留 YAML 编解码与清洗工具，来源 `cliq-hub-backend/pkg/yaml/codec.go:1`。
  - 模板校验：迁移并导出校验规则，来源 `cliq/services/template_service.go:225` 及相关验证入口（`ValidateYAMLTemplate`）。
- 更新依赖：`apps/cliq-app` 与 `apps/cliq-hub-backend` 统一从 `@repo/shared-go-lib` 引用模板模型、解析与校验逻辑。
- 调整 `apps/cliq-app/go.mod` 与 `apps/cliq-hub-backend/go.mod` 的模块路径；通过 `go.work` 绑定本地模块，确保正确引用。

### 阶段 4：Nx 项目配置
- `apps/cliq-app/project.json`
  - `build`：包装 `wails build`（前置 `pnpm -C apps/wails-app/frontend build`）。
  - `serve`：包装 `wails dev`（沿用 `wails.json` 的 dev 配置）。
- `apps/cliq-hub-backend/project.json`
  - `build`：`go build ./cmd/server`。
  - `serve`：`go run ./cmd/server`。
- `apps/cliq-frontend/project.json`
  - `build`/`serve`：按 `@nx/vue` + `vite` 约定。
- `packages/shared-vue-ui/project.json`
  - `build`：Vite 库模式产出（ES + 类型）。
- `packages/shared-go-lib/project.json`
  - `build`：`go build` 或 `go vet`/`go test`。

### 阶段 5：CI 与脚本一致性（GitHub Actions 适配）
- 增加 GitHub Actions 工作流，针对 push 与 pull_request 自动构建/测试所有 apps：
  - 触发：`on: [push, pull_request]`
  - 环境：Node（pnpm）与 Go（go.work），缓存 pnpm 与 Go modules。
  - 步骤：
    - `actions/setup-node` + `pnpm/action-setup` 安装与缓存依赖。
    - `actions/setup-go` 安装 Go，启用缓存。
    - 执行 `pnpm install --frozen-lockfile`。
    - 构建：`pnpm nx run cliq-app:build`、`pnpm nx run cliq-hub-backend:build`、`pnpm nx run cliq-frontend:build`。
    - 测试：如存在测试，执行 `pnpm nx run <project>:test` 或 `go test ./...`。
- 更新或新增工作流：根安装与 `nx run` 构建（Wails、Go 后端、纯前端）。
- 根 `package.json` 增加快捷脚本：
  - `nx build cliq-app`、`nx serve cliq-app`
  - `nx build cliq-hub-backend`、`nx serve cliq-hub-backend`
  - `nx build cliq-frontend`、`nx serve cliq-frontend`

### 阶段 6：验证与回归
- 前端：`pnpm nx graph` 验证依赖；`SettingsPage.vue:44` 的 Base URL 与 `TemplateGenerator.vue:99` 的远程生成在新结构下可达。
- 后端：`go run ./cmd/server` 路由可用，智能生成接口工作；LLM 配置兼容 `client.go:37`。
- Wails：`wails dev/build` 在 Nx 任务下可运行，`frontend:install/build/dev` 行为与旧版一致（源自 `wails.json:1`）。

## 关键设计与注意点
- Go 模块路径：优先使用本地短路径（如 `repo/shared-go-lib`）并通过 `go.work` 绑定，确保能正确被 `apps/wails-app` 与 `apps/go-backend` 引用；后续若对外发布，再切换到 `github.com/<org>/<repo>/shared-go-lib`。
- 领域模型统一：`Template`/`Command`/`Variable` 保持一处定义，减少序列化/校验重复。
- 前端依赖统一：所有 Node 依赖通过根 `pnpm` 管理；应用通过 workspace 互相依赖。
- Wails 集成：Wails 的前端命令与 dev server 配置由 `wails.json` 继续驱动，Nx 负责编排。


## 快速验证（迁移符合预期）
- 安装依赖：`pnpm install`
- 查看依赖图：`pnpm nx graph`
- 运行 Wails 开发：`pnpm nx run cliq-app:serve`（应打开 Wails dev，前端由 `wails.json` 驱动）
- 构建 Wails 应用：`pnpm nx run cliq-app:build`
- 运行 Go 后端：`pnpm nx run cliq-hub-backend:serve` 或 `go run ./apps/cliq-hub-backend/cmd/server`
- 运行纯 Web 前端：`pnpm nx run cliq-frontend:serve`（应显示 Hello World）
- 验证共享 UI：在 `cliq-app` 与 `cliq-frontend` 中均可 `import DynamicCommandForm from '@repo/shared-vue-ui'` 并渲染基础表单
- 验证共享 Go 库：
  - `apps/cliq-app` 能调用共享库提供的解析/校验方法（如 `ValidateYAMLTemplate`）；
  - `apps/cliq-hub-backend` 的生成/校验接口使用共享模型与工具构建并通过。

## 风险与回滚
- 风险：模块路径切换导致编译/导入失败；Wails 前端锁文件与根安装冲突；抽取组件后样式或自动导入异常。
- 回滚策略：保留迁移前分支；分期实施，每阶段通过 `nx build`/`go build` 验证再推进；抽取失败时先还原为应用内私有组件。