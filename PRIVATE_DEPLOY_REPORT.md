# xfchat_cli（私有化 lark-cli）私有化部署调研与测试报告

## 一、项目概览

**xfchat_cli** 是基于官方 `lark-cli` 的私有化 fork，Go 语言编写，MIT 协议；当前发行二进制名为 `xfchat_cli`，配置目录为 `~/.xfchat_cli`。

- **规模**：200+ 命令，19 个 AI Agent Skills，覆盖消息、文档、多维表格、电子表格、日历、任务、知识库、通讯录、视频会议等核心域
- **架构**：三层命令体系 — Shortcuts（`+` 前缀）→ API Commands → Raw API
- **覆盖域**：日历、即时消息、文档、云盘、多维表格、电子表格、任务、Wiki、通讯录、视频会议
- **运行态说明**：`lark-mail` skill 仍随包保留，但当前私有化运行时默认禁用 `mail` 服务命令，不纳入 prod 能力矩阵

---

## 二、登录/鉴权机制

### 2.1 认证流程：OAuth 2.0 Device Flow / Authorization Code Flow

公有云默认采用 **OAuth 2.0 设备授权流（RFC 8628）**；私有化环境不支持 Device Flow 时，当前 fork 通过 `--web` 回落到 Authorization Code Flow：

```
1. xfchat_cli config init            → 配置 AppID + AppSecret + Brand
2. xfchat_cli auth login --recommend 或 xfchat_cli auth login --web
   ├─ Device Flow：POST {accounts}/oauth/v1/device_authorization
   ├─ Authorization Code Flow：浏览器打开授权页并回调 localhost
   ├─ 用户在浏览器完成扫码/授权
   └─ CLI 获取 User Access Token (UAT) + Refresh Token
3. 获取 User Access Token (UAT) + Refresh Token
4. Token 存入系统 Keychain（macOS Keychain / Linux AES-256-GCM）
```

### 2.2 两种身份模式

| 身份 | Token 类型 | 获取方式 | 适用场景 |
|------|-----------|---------|---------|
| **Bot（`--as bot`）** | Tenant Access Token | SDK 用 AppID+AppSecret 自动获取 | 发消息、读写表格等机器人权限操作 |
| **User（`--as user`）** | User Access Token | OAuth 浏览器授权（私有化默认 `--web`） | 通讯录搜索、用户日历等需要用户身份的操作 |

### 2.3 所需凭证

| 凭证 | 来源 | 存储位置 |
|------|------|---------|
| AppID | 飞书开放平台创建应用 | `~/.xfchat_cli/config.json` |
| AppSecret | 飞书开放平台创建应用 | OS Keychain |
| User Access Token | OAuth 浏览器授权获取 | OS Keychain |
| Refresh Token | 随 UAT 一起返回 | OS Keychain |

### 2.4 Token 刷新机制

- 自动刷新：Token 过期前 **5 分钟** 自动触发
- 跨进程安全：使用 `flock` 文件锁防止并发刷新冲突
- Refresh Token 有效期约 7 天

---

## 三、私有化部署适配

### 3.1 原始硬编码端点

代码位置：`internal/core/types.go`

| Brand | Open API | Accounts (OAuth) | MCP |
|-------|----------|-------------------|-----|
| `feishu` | `https://open.feishu.cn` | `https://accounts.feishu.cn` | `https://mcp.feishu.cn` |
| `lark` | `https://open.larksuite.com` | `https://accounts.larksuite.com` | `https://mcp.larksuite.com` |

所有涉及域名的位置（共 4 处）：

| 文件 | 用途 |
|------|------|
| `internal/core/types.go:24-39` | Open/Accounts/MCP 端点解析 |
| `internal/auth/device_flow.go:54-60` | OAuth 端点派生（从 types.go 调用） |
| `internal/registry/remote.go:67-83` | API 元数据拉取 URL |
| `shortcuts/common/mcp_client.go:24-26` | MCP 调用端点（从 types.go 调用） |

### 3.2 代码改动

总改动量约 **30 行**，涉及 3 个文件：

#### 3.2.1 `internal/core/types.go` — 新增端点覆盖机制

```go
// EndpointOverrides allows overriding resolved endpoints for private deployments.
type EndpointOverrides struct {
    Open     string `json:"open,omitempty"`
    Accounts string `json:"accounts,omitempty"`
    MCP      string `json:"mcp,omitempty"`
}

var endpointOverrides *EndpointOverrides

func SetEndpointOverrides(o *EndpointOverrides) {
    endpointOverrides = o
}
```

`ResolveEndpoints` 在返回前检查并应用覆盖：

```go
if endpointOverrides != nil {
    if endpointOverrides.Open != "" { ep.Open = endpointOverrides.Open }
    if endpointOverrides.Accounts != "" { ep.Accounts = endpointOverrides.Accounts }
    if endpointOverrides.MCP != "" { ep.MCP = endpointOverrides.MCP }
}
```

#### 3.2.2 `internal/core/config.go` — 配置文件支持 endpoints 字段

`AppConfig` 新增：

```go
Endpoints *EndpointOverrides `json:"endpoints,omitempty"`
```

`RequireConfig` 加载配置时自动应用：

```go
if app.Endpoints != nil {
    SetEndpointOverrides(app.Endpoints)
}
```

#### 3.2.3 `internal/registry/remote.go` — 元数据 URL 走统一端点

```go
// 改前：硬编码 open.feishu.cn / open.larksuite.com
// 改后：
base := core.ResolveEndpoints(configuredBrand).Open + "/api/tools/open/api_definition"
```

### 3.3 使用方式

编辑 `~/.xfchat_cli/config.json`：

```json
{
  "apps": [
    {
      "appId": "你的AppID",
      "appSecret": "你的AppSecret",
      "brand": "feishu",
      "endpoints": {
        "open": "https://open.xfchat.iflytek.com",
        "accounts": "https://open.xfchat.iflytek.com"
      },
      "users": []
    }
  ]
}
```

> `accounts` 是 OAuth 认证域名，如果私有化环境认证和 API 在同一个域名下就填一样的。`mcp` 字段可选，私有化环境没有 MCP 服务时不填即可。

---

## 四、测试结果

### 4.1 测试环境

- **私有化平台**：讯飞内部飞书 `https://open.xfchat.iflytek.com`
- **应用**：星火派-test（`cli_a9084c0706b8d379`）
- **测试身份**：Bot
- **测试时间**：2026-03-28

### 4.2 功能测试明细

#### 连通性检查（`xfchat_cli doctor`）

```
config_file    ✅ config.json found
app_resolved   ✅ app: cli_a9084c0706b8d379 (feishu)
endpoint_open  ✅ https://open.xfchat.iflytek.com reachable
endpoint_mcp   ⚠️ mcp.feishu.cn reachable（指向公有云，需配置 mcp 端点）
token_exists   ❌ no user logged in（Bot 模式不需要）
```

#### Bot 信息

```bash
./xfchat_cli api GET /open-apis/bot/v3/info --as bot
```

```json
{"bot": {"app_name": "星火派-test", "activate_status": 2}, "code": 0, "msg": "ok"}
```

**结果：✅ 通过**

#### 即时消息（IM）

| 操作 | 命令 | 结果 |
|------|------|------|
| 发送文本消息 | `im +messages-send --text "..." --as bot` | ✅ 通过 |
| 发送 Markdown 富文本 | `im +messages-send --markdown $'...' --as bot` | ✅ 通过 |
| 回复消息 | `im +messages-reply --message-id om_xxx --text "..." --as bot` | ✅ 通过 |
| 获取群信息 | `api GET /open-apis/im/v1/chats/{id} --as bot` | ✅ 通过（19人群，2机器人） |
| 拉取历史消息 | `im +chat-messages-list --chat-id oc_xxx --as bot` | ✅ 通过 |
| 群列表 | `api GET /open-apis/im/v1/chats --as bot` | ✅ 通过（返回5个群） |

> **注意**：Markdown 内容中的换行需要使用 shell 的 `$'...\n...'` 语法，而不是双引号中的 `\n`（shell 不会将双引号中的 `\n` 转为真正换行符）。

#### 云盘（Drive）

| 操作 | 命令 | 结果 |
|------|------|------|
| 文件列表 | `api GET /open-apis/drive/v1/files --as bot` | ✅ 通过（返回多维表格、文档等文件） |

#### 多维表格（Bitable）

| 操作 | 命令 | 结果 |
|------|------|------|
| 表格列表 | `api GET /open-apis/bitable/v1/apps/{id}/tables --as bot` | ✅ 通过（2张表） |
| 记录读取 | `api GET /open-apis/bitable/v1/apps/{id}/tables/{id}/records --as bot` | ✅ 通过（含图片附件字段） |

#### 任务（Task）

| 操作 | 命令 | 结果 |
|------|------|------|
| 创建任务 | `api POST /open-apis/task/v2/tasks --as bot --data '{...}'` | ✅ 通过（task_id: t398481） |

#### 日历（Calendar）

| 操作 | 命令 | 结果 |
|------|------|------|
| 日历列表 | `api GET /open-apis/calendar/v4/calendars --as bot` | ✅ 通过（含主日历 + 多个共享日历） |

#### 电子表格（Sheets）

| 操作 | 命令 | 结果 |
|------|------|------|
| 创建表格 | `api POST /open-apis/sheets/v3/spreadsheets --as bot --data '{...}'` | ✅ 通过 |

#### 文档创建（Docs — MCP 依赖）

| 操作 | 命令 | 结果 |
|------|------|------|
| 创建文档 | `docs +create --title "..." --markdown "..." --as bot` | ❌ 失败 |

失败原因：`docs +create` 通过 MCP 协议调用，MCP 端点默认指向公有云 `mcp.feishu.cn`，私有化环境的 token 在公有云无效。需要私有化环境部署 MCP 服务后在 config 中配置 `"mcp"` 端点。

#### 通讯录搜索（Contact）

| 操作 | 命令 | 结果 |
|------|------|------|
| 搜索用户 | `contact +search-user --query "..." --as bot` | ⏭ 跳过（仅支持 user 身份） |

### 4.3 测试汇总

| 功能域 | 测试项 | 状态 | 备注 |
|--------|--------|------|------|
| Bot 信息 | 获取应用信息 | ✅ | |
| IM | 发送文本 | ✅ | |
| IM | 发送富文本 | ✅ | 需用 `$'...'` 传换行 |
| IM | 回复消息 | ✅ | |
| IM | 群信息/群列表 | ✅ | |
| IM | 历史消息 | ✅ | |
| Drive | 文件列表 | ✅ | |
| Bitable | 表格/记录 | ✅ | |
| Task | 创建任务 | ✅ | |
| Calendar | 日历列表 | ✅ | |
| Sheets | 创建表格 | ✅ | |
| Docs | 创建文档(MCP) | ❌ | 需私有化 MCP 服务 |
| Contact | 搜索用户 | ⏭ | 需 User 身份 |

**通过率：11/13（84.6%），其中 1 项需 MCP 私有化部署，1 项需 User 身份授权。**

---

## 五、Agent 集成评估

### 5.1 Agent 友好度评分：9.5/10

| 维度 | 评分 | 说明 |
|------|------|------|
| 结构化输出 | 10/10 | 支持 JSON / NDJSON / CSV / Table，统一信封格式 |
| 错误处理 | 10/10 | 分类退出码(0-5)，带 hint 和 console_url |
| 非交互模式 | 10/10 | `--no-wait` 非阻塞认证，`--dry-run` 预览 |
| 安全性 | 9/10 | 输入注入防护、ANSI 转义剥离、Keychain 存储 |
| 域覆盖 | 9/10 | 覆盖绝大部分飞书核心能力；当前私有化运行时默认禁用 `mail` |
| 可发现性 | 10/10 | `schema` 自省命令，19 个预写 AI Skills |
| 幂等性 | 9/10 | 消息发送支持 `--idempotency-key` |

### 5.2 Agent 典型使用模式

```bash
# 认证（非阻塞，适合 Agent）
result=$(./xfchat_cli auth login --domain calendar,im --no-wait --format json)
# → 返回 verification_url，Agent 把 URL 给用户

# 发消息
./xfchat_cli im +messages-send --chat-id oc_xxx --text "Hello" --as bot

# 查日程
./xfchat_cli calendar +agenda --format json

# 预览后执行（安全模式）
./xfchat_cli im +messages-send --chat-id oc_xxx --text "Hi" --dry-run
./xfchat_cli im +messages-send --chat-id oc_xxx --text "Hi"
```

### 5.3 已有 19 个 AI Skills

```bash
# 查看随包 skills 目录
xfchat_cli skills path

# 查看已随包携带的 skills
xfchat_cli skills list
```

### 5.4 局限性

| 局限 | 影响 |
|------|------|
| 无内置并行执行 | Agent 需自行用 shell `&` 并行 |
| 无输出缓存 | 重复查询需 Agent 自行缓存 |
| 实时事件需长连接 | `event +subscribe` 需要持续运行的进程 |

---

## 六、结论

### 私有化部署

**可行性高，改动极小**（约 30 行代码，3 个文件）。核心通过 `config.json` 的 `endpoints` 字段覆盖默认域名，对原有代码无侵入。11/13 项功能在讯飞私有化飞书（`open.xfchat.iflytek.com`）上验证通过。

### 待解决

1. **MCP 服务**：部分 shortcut（如 `docs +create`）依赖 MCP，需要私有化环境部署 MCP 服务
2. **User 身份**：通讯录搜索等操作需要 OAuth 浏览器授权；私有化环境下应优先验证 `--web` 回调链路
3. **元数据接口**：`/api/tools/open/api_definition` 在私有化环境可能不存在，可通过 `LARKSUITE_CLI_REMOTE_META=off` 禁用

### Agent 使用

**非常适合**。工具本身是 Agent-Native 设计，结构化 I/O、非交互认证、错误分类、19 个随包 Skills，可直接作为 AI Agent 的飞书操控工具使用。
