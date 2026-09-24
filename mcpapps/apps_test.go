package mcpapps

import (
	"strings"
	"testing"
)

func TestHTMLRendersSharedMCPAppTemplate(t *testing.T) {
	html := HTML("recall", "Recall")
	for _, marker := range []string{
		`<title>Recall</title>`,
		`expectedView="recall"`,
		`rpcRequest("ui/initialize"`,
		`protocolVersion:"2026-01-26"`,
		`rpcNotify("ui/notifications/initialized"`,
		`message.method==="ui/notifications/tool-result"`,
	} {
		if !strings.Contains(html, marker) {
			t.Fatalf("shared MCP App missing marker %q", marker)
		}
	}
	for _, forbidden := range []string{"{{VIEW}}", "{{TITLE}}", "window.openai", "openai/widget", "openai/outputTemplate"} {
		if strings.Contains(html, forbidden) {
			t.Fatalf("shared MCP App contains forbidden marker %q", forbidden)
		}
	}
}

func TestHTMLUsesStandardHostContextThemeAndSemanticTokens(t *testing.T) {
	html := HTML("workflow", "Workflow")
	for _, marker := range []string{
		`:root{color-scheme:light dark;`,
		`:root[data-theme="dark"]{color-scheme:dark;`,
		`@media(prefers-color-scheme:dark)`,
		`--ad-text-primary:`,
		`var(--ad-text-primary)`,
		`const initialized=await rpcRequest("ui/initialize"`,
		`applyHostContext(initialized&&initialized.hostContext)`,
		`message.method==="ui/notifications/host-context-changed"`,
		`applyHostContext(message.params)`,
		`context.styles.variables`,
		`window.matchMedia("(prefers-color-scheme: dark)")`,
	} {
		if !strings.Contains(html, marker) {
			t.Fatalf("shared MCP App missing theme marker %q", marker)
		}
	}
	if strings.Contains(html, `:root{color-scheme:light;font:`) {
		t.Fatal("shared MCP App still forces a light-only color scheme")
	}
}

func TestHTMLLocalizesFromBrowserAndHostLocale(t *testing.T) {
	html := HTML("agentdock_context", "Context")
	for _, marker := range []string{
		`"zh-CN":{`,
		`resolveLocale([...(navigator.languages||[]),navigator.language])`,
		`candidate==="zh-cn"`,
		`candidate==="zh-hans"`,
		`const applyLocale=value=>{const next=resolveLocale([value])`,
		`context.locale||context.language`,
		`document.documentElement.lang=locale`,
		`message.method==="ui/notifications/host-context-changed"`,
		`applyHostContext(message.params)`,
		`t("nodeUnavailable"`,
		`t("contextUnavailable")`,
		`t("nodeOffline")`,
		`t("moreItems"`,
		`compactShell({action:"context",title:t("capabilities")}`,
	} {
		if !strings.Contains(html, marker) {
			t.Fatalf("shared MCP App missing locale marker %q", marker)
		}
	}

	// 这些文案曾直接写在渲染逻辑里，必须只存在于语言字典，避免再次出现中英混排。
	if strings.Contains(html, `candidate.startsWith("zh-")`) {
		t.Fatal("shared MCP App must not coerce every Chinese locale, including Traditional Chinese, to zh-CN")
	}

	for _, direct := range []string{
		`+" 当前不可用"`,
		`?"不可用":`,
		`?"在线":"离线"`,
		`"还有 "+(items.length-limit)+" 项"`,
		`title:"Capabilities"`,
		`label:"State"`,
	} {
		if strings.Contains(html, direct) {
			t.Fatalf("shared MCP App still contains hard-coded locale text %q", direct)
		}
	}
}

func TestAgentDockContextRendersPluginsAndMultiACP(t *testing.T) {
	html := HTML("agentdock_context", "Context")
	for _, marker := range []string{
		`plugins:"Plugins"`,
		`defaultProfile:"Default"`,
		`function pluginContextItems(items)`,
		`function acpContextItems(acp)`,
		`const plugins=pluginContextItems(data.plugins)`,
		`const acps=acpContextItems(data.acp)`,
		`appendContextOverview(overview,plugins.length,t("plugins"))`,
		`appendContextOverview(overview,acps.length,"ACP")`,
		`appendContextSection(groups,t("plugins"),plugins,8)`,
		`appendContextSection(groups,"ACP",acps,8)`,
		`plugins:node.context.plugins`,
		`id===defaultProfile?t("defaultProfile")`,
		`Renderer-only fallback for older nodes`,
		`skills:"Skills"`,
		`contextPill(summary.skills.length,t("skills"))`,
		`contextPill(summary.plugins.length,t("plugins"))`,
		`contextPill(summary.mcps.length,"MCP")`,
		`contextPill(summary.workflows.length,t("workflow"),"workflow")`,
		`contextPill(summary.recall?t("on"):t("off"),t("recall"))`,
	} {
		if !strings.Contains(html, marker) {
			t.Fatalf("AgentDock context MCP App missing Plugin/ACP marker %q", marker)
		}
	}
	if strings.Contains(html, `const acp=isObject(data.acp)&&data.acp.enabled?[{name:String(data.acp.agent||"ACP")`) {
		t.Fatal("AgentDock context renderer still uses single-ACP primary path")
	}
	if strings.Contains(html, `contextPill(summary.acps.length,"ACP")`) {
		t.Fatal("compact AgentDock context must not show ACP")
	}
}

func TestHTMLZhCNPreservesHistoricalMixedTerminology(t *testing.T) {
	html := HTML("agentdock_context", "Context")
	start := strings.Index(html, `"zh-CN":{`)
	end := strings.Index(html, `const matchLocale=value=>`)
	if start < 0 || end <= start {
		t.Fatal("shared MCP App zh-CN message dictionary is missing")
	}
	zhCN := html[start:end]

	// zh-CN 保持 i18n 前已经成熟的中英混排；协议、产品和结构字段不做机械中文化。
	for _, marker := range []string{
		`status:"STATUS"`,
		`workflow:"Workflow"`,
		`capabilities:"Capabilities"`,
		`devices:"Devices"`,
		`artifact:"Artifact"`,
		`type:"Type"`,
		`action_task:"TASK"`,
		`state_pending:"pending"`,
		`moreItems:"还有 {count} 项"`,
		`online:"在线"`,
		`unavailable:"不可用"`,
		`contextUnavailable:"Context 暂不可用"`,
	} {
		if !strings.Contains(zhCN, marker) {
			t.Fatalf("shared MCP App zh-CN copy drifted from historical UI: missing %q", marker)
		}
	}

	for _, forbidden := range []string{
		`status:"状态"`,
		`workflow:"工作流"`,
		`capabilities:"能力"`,
		`artifact:"产物"`,
		`action_task:"任务"`,
	} {
		if strings.Contains(zhCN, forbidden) {
			t.Fatalf("shared MCP App zh-CN copy is over-translated: %q", forbidden)
		}
	}
}

// ── 自动展开回归测试 ──────────────────────────────────────────────────────────

// TestAutoExpandThresholdsPresent 锁定高度/宽度阈值常量名与数值，防止误改。
func TestAutoExpandThresholdsPresent(t *testing.T) {
	html := HTML("recall", "Recall")
	for _, marker := range []string{
		`AUTO_EXPAND_MIN_HEIGHT=500`,
		`AUTO_EXPAND_MIN_WIDTH=400`,
		`window.innerHeight>=AUTO_EXPAND_MIN_HEIGHT&&window.innerWidth>=AUTO_EXPAND_MIN_WIDTH`,
	} {
		if !strings.Contains(html, marker) {
			t.Fatalf("auto-expand: missing threshold marker %q", marker)
		}
	}
}

// TestAutoExpandWindowResizeListener 确认使用 window resize 事件驱动 tryAutoExpand，
// 而非把自动展开逻辑混入 ResizeObserver（保持 ResizeObserver 职责单一）。
func TestAutoExpandWindowResizeListener(t *testing.T) {
	html := HTML("recall", "Recall")
	if !strings.Contains(html, `window.addEventListener("resize",tryAutoExpand`) {
		t.Fatal("auto-expand: must register tryAutoExpand on window resize event")
	}
	// ResizeObserver 回调必须只传 reportSize，不应直接调用 tryAutoExpand
	if strings.Contains(html, `new ResizeObserver(()=>{`) {
		t.Fatal("auto-expand: ResizeObserver callback must not be an arrow function calling tryAutoExpand; keep it as ResizeObserver(reportSize)")
	}
}

// TestAutoExpandUserToggleBlocksAutoExpand 确认 userToggledExpand=true 后
// tryAutoExpand 立即返回，不会覆盖用户选择。
func TestAutoExpandUserToggleBlocksAutoExpand(t *testing.T) {
	html := HTML("recall", "Recall")
	// tryAutoExpand 的第一条守卫必须是 userToggledExpand 检查
	if !strings.Contains(html, `function tryAutoExpand(){`+"\n"+`    if(userToggledExpand)return;`) &&
		!strings.Contains(html, "function tryAutoExpand(){if(userToggledExpand)return;") {
		// 允许两种风格（多行 / 单行），只要 userToggledExpand 是第一个 guard
		if !strings.Contains(html, "if(userToggledExpand)return;") {
			t.Fatal("auto-expand: tryAutoExpand must guard on userToggledExpand first")
		}
	}
	// 点击事件必须设置 userToggledExpand=true
	if !strings.Contains(html, "userToggledExpand=true;") {
		t.Fatal("auto-expand: click handler must set userToggledExpand=true")
	}
}

// TestAutoExpandReRenderRestoresUserExpandedState 确认 re-render（新 compactShell）
// 时读取 userExpandedState 恢复用户已明确的展开状态，不丢失。
func TestAutoExpandReRenderRestoresUserExpandedState(t *testing.T) {
	html := HTML("recall", "Recall")
	for _, marker := range []string{
		`let userExpandedState=false;`,
		`userExpandedState=!expanded;`,
		`if(userToggledExpand){`,
		`if(userExpandedState){`,
	} {
		if !strings.Contains(html, marker) {
			t.Fatalf("auto-expand: re-render state restore missing marker %q", marker)
		}
	}
}

// TestAutoExpandIsUnidirectional 确认 tryAutoExpand 只做展开操作，
// 不包含任何自动折叠（setAttribute aria-expanded false / panel.hidden=true）路径。
func TestAutoExpandIsUnidirectional(t *testing.T) {
	html := HTML("recall", "Recall")

	// 提取 tryAutoExpand 函数体（简单定界：从函数声明到下一个顶级 function）
	start := strings.Index(html, "function tryAutoExpand(){")
	if start < 0 {
		t.Fatal("auto-expand: tryAutoExpand function not found")
	}
	// 找到函数结束的 "}\n\n" 模式（函数体后有空行分隔）
	body := html[start:]
	end := strings.Index(body[1:], "\n  function ")
	if end < 0 {
		end = len(body)
	} else {
		end++ // 跳过 body[1:] 的偏移
	}
	fnBody := body[:end]

	// 函数体内不得出现自动折叠的标志
	if strings.Contains(fnBody, `setAttribute("aria-expanded","false"`) {
		t.Fatal("auto-expand: tryAutoExpand must not auto-collapse (no aria-expanded=false)")
	}
	if strings.Contains(fnBody, "panel.hidden=true") {
		t.Fatal("auto-expand: tryAutoExpand must not set panel.hidden=true (no auto-collapse)")
	}
	// 函数体内必须包含展开操作
	if !strings.Contains(fnBody, `setAttribute("aria-expanded","true"`) {
		t.Fatal("auto-expand: tryAutoExpand must set aria-expanded=true")
	}
}
