# agentdock-protocol

Shared contracts used by AgentDock and NexusDock.

The repository intentionally has two logical layers:

- the root `protocol` package owns only the AgentDock ↔ NexusDock Bridge wire protocol: node pairing and connect HTTP shapes, envelopes, Hello capabilities, operations, and MCP App resource identities/contracts;
- `mcpcontract` owns only the canonical model-facing MCP contracts shared by the two entrypoints: input/output schemas, annotations, and bounded behavior vectors.

Neither package owns AgentDock runtime behavior, NexusDock stores, renderer HTML, Recall persistence, Workflow persistence, or HTTP handlers. Those remain in their application repositories.

## 共享 MCP Apps

独立的 `mcpapps` 包维护两个入口复用的 UI 文档。`HTML("view_image", ...)` 返回图片组件，资源身份与兼容契约由根包的 `ImageUIResourceURI` / `ImageUIContract` 定义。图片内容仍使用标准 MCP image 块，各入口自行注册资源和绑定工具，不将执行或 HTTP 逻辑放进共享包。

图片组件只处理本次工具结果，并先显示预览；只有用户点击确认按钮后才把图片提供给模型。点击后优先使用标准 MCP Apps `ui/update-model-context` / `ui/message`，宿主未暴露标准图片能力时再回退到 ChatGPT 的 `uploadFile` / `widgetState.imageIds`。未提供任何模型图片通道的客户端只显示预览；不保证同轮视觉，不读取额外文件、不默认保存到文件库，也不扩大 CSP 网络域名。

验证命令：`go test ./...`、`go test -race ./...`、`go vet ./...`、`node --test mcpapps/app.test.mjs mcpapps/image.test.mjs`（Node 20+）。组件上传、宿主状态更新、错误、重复点击、换图竞争、ACP 会话详情与历史消息投影及非父 frame 消息均有行为测试。真实宿主交互仍须独立验收，不能仅凭单元测试推定。

## 节点配对与连接

`pairing.go` 定义节点建立 Bridge 之前仅有的两个 HTTP 入口，路径都相对于节点配置的 Nexus endpoint：

- `POST /v1/nodes/pair`：请求体 `PairRequest`（`code`、`device_id`、`name`），成功时返回 `201 Created` 与 `PairResponse`（`node.id`、`device_token`）。其他状态码均表示失败，错误体由各实现自行定义。服务端可以在响应中增加字段，节点必须忽略未知字段。
- `GET /v1/nodes/connect`：WebSocket 升级请求携带 `Authorization: Bearer <device_token>`；节点先发送 `node.hello`，服务端回复 `node.ready`（`protocol_version`、`heartbeat_ms`、不带路径的 HTTPS `public_url`）。`Hello.device_id` 必须等于配对时的 `device_id`。

这些消息不携带任何账号、租户或工作区选择字段：配对码在服务端签发时已经决定归属，节点只认识自己配置的 endpoint。`testdata/` 中的 golden fixture 是两端实现的对照样例。
