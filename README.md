# agentdock-protocol

Shared contracts used by AgentDock and NexusDock.

The repository intentionally keeps shared contracts in three narrow layers:

- the root `protocol` package owns only the AgentDock ↔ NexusDock Bridge wire protocol: envelopes, Hello capabilities, operations, and MCP App resource identities/contracts;
- `mcpcontract` owns only the canonical model-facing MCP contracts shared by the two entrypoints: input/output schemas, annotations, and bounded behavior vectors;
- `runtimecontract/vN` owns versioned OpenAPI wire shapes for `/internal/runtime/*`. Published versions remain available so AgentDock and NexusDock can evolve independently; breaking changes create a new version instead of rewriting an existing one.

These packages define wire contracts, not application behavior. AgentDock runtime behavior, NexusDock stores, renderer HTML, Recall persistence, Workflow persistence, and HTTP handlers remain in their application repositories.

## 共享 MCP Apps

独立的 `mcpapps` 包维护两个入口复用的 UI 文档。`HTML("view_image", ...)` 返回图片组件，资源身份与兼容契约由根包的 `ImageUIResourceURI` / `ImageUIContract` 定义。图片内容仍使用标准 MCP image 块，各入口自行注册资源和绑定工具，不将执行或 HTTP 逻辑放进共享包。

图片组件只处理本次工具结果，并先显示预览；只有用户点击确认按钮后才把图片提供给模型。点击后优先使用标准 MCP Apps `ui/update-model-context` / `ui/message`，宿主未暴露标准图片能力时再回退到 ChatGPT 的 `uploadFile` / `widgetState.imageIds`。未提供任何模型图片通道的客户端只显示预览；不保证同轮视觉，不读取额外文件、不默认保存到文件库，也不扩大 CSP 网络域名。

验证命令：`go test ./...`、`go test -race ./...`、`go vet ./...`、`node --test mcpapps/image.test.mjs`（Node 20+）。组件上传、宿主状态更新、错误、重复点击、换图竞争及非父 frame 消息均有行为测试。真实模型识图仍须在宿主中独立验收，不能仅凭单元测试推定。
