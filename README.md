# agentdock-protocol

Shared contracts used by AgentDock and NexusDock.

The repository intentionally has two logical layers:

- the root `protocol` package owns only the AgentDock ↔ NexusDock Bridge wire protocol: envelopes, Hello capabilities, operations, and MCP App resource identities/contracts;
- `mcpcontract` owns only the canonical model-facing MCP contracts shared by the two entrypoints: input/output schemas, annotations, and bounded behavior vectors.

Neither package owns AgentDock runtime behavior, NexusDock stores, renderer HTML, Recall persistence, Workflow persistence, or HTTP handlers. Those remain in their application repositories.

## 共享 MCP Apps

独立的 `mcpapps` 包维护两个入口复用的 UI 文档。`HTML("view_image", ...)` 返回图片组件，资源身份与兼容契约由根包的 `ImageUIResourceURI` / `ImageUIContract` 定义。图片内容仍使用标准 MCP image 块，各入口自行注册资源和绑定工具，不将执行或 HTTP 逻辑放进共享包。

图片组件只预览本次工具结果，用户点击后通过 ChatGPT 宿主的 `uploadFile` 和 `widgetState.imageIds` 附加到后续对话。未提供宿主接口的客户端只显示预览；不保证同轮自动视觉，不读取额外文件、不使用 OCR、不默认保存到文件库，也不扩大 CSP 网络域名。

验证命令：`go test ./...`、`go test -race ./...`、`go vet ./...`、`node --test mcpapps/image.test.mjs`（Node 20+）。组件上传、宿主状态更新、错误、重复点击、换图竞争及非父 frame 消息均有行为测试。真实模型识图仍须在宿主中独立验收，不能仅凭单元测试推定。
