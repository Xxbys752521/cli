# Calendar 模块测试记录

测试时间：2026-03-31
测试环境：i 讯飞私有化飞书
OpenAPI：`https://open.xfchat.iflytek.com`
OAuth：`https://accounts.xfchat.iflytek.com`
应用：`cli_a9084c0706b8d379` / `星火派-test`

## 1. 环境与鉴权

### 1.1 CLI 与配置

初始化命令：

```bash
./xfchat_cli config init --app-id cli_a9084c0706b8d379 --app-secret-stdin --brand feishu
```

补充配置：

```json
{
  "endpoints": {
    "open": "https://open.xfchat.iflytek.com",
    "accounts": "https://accounts.xfchat.iflytek.com"
  }
}
```

### 1.2 用户授权

授权命令：

```bash
./xfchat_cli auth login --web --scope 'calendar:calendar calendar:calendar:readonly contact:contact.base:readonly contact:user.base:readonly contact:user:search'
```

实际授权后 `auth status` 摘要：

```json
{
  "appId": "cli_a9084c0706b8d379",
  "brand": "feishu",
  "defaultAs": "auto",
  "identity": "user",
  "scope": "auth:user.id:read calendar:calendar calendar:calendar:readonly contact:contact.base:readonly contact:contact:readonly_as_app contact:user.base:readonly contact:user:search",
  "tokenStatus": "valid",
  "userName": "付涛 taofu5",
  "userOpenId": "ou_84c79a9284092a69283066a89e549251"
}
```

结论：

- 登录成功
- User 身份同时拥有 `calendar:calendar`（读写）和 `calendar:calendar:readonly`（只读）
- Bot 身份通过 tenant_access_token 鉴权，拥有 `calendar:calendar`

### 1.3 环境检查

检查命令：

```bash
./xfchat_cli doctor --offline
./xfchat_cli api GET /open-apis/bot/v3/info --as bot
./xfchat_cli contact +get-user --format json
```

关键输出：

```json
{
  "app_resolved": "cli_a9084c0706b8d379 (feishu)",
  "bot.app_name": "星火派-test (HTTP 404, 私有化不支持 bot info 接口)",
  "user.open_id": "ou_84c79a9284092a69283066a89e549251"
}
```

结论：环境可用，Bot 与 User 身份都已就绪。

## 2. 测试对象与素材

### 2.1 日历对象

- User 主日历 ID：`feishu.cn_FnVWAps00pOJSzu2pnkCAc@group.calendar.feishu.cn`（付涛 taofu5）
- Bot 主日历 ID：`feishu.cn_cn2LpP9WiCZuwxhMzlczUd@group.calendar.feishu.cn`（星火派-test）

### 2.2 日程对象

已有日程（用于读取测试）：

- `a9173edb-e446-4346-8f7c-ed348fac2cc0_0`：上线功能讨论会（2026-03-24 15:00-17:00）
- `4d5bb05f-cbb1-8f29-751c-1233dbbd11f6_0`：RDG-AI提效案例分享·第二期（2026-03-26 19:00-21:00）
- `82b71bae-8cc7-483f-9196-227ec491f5b1_0`：Updated Test Event（Bot 历史测试日程）

Bot 创建的测试日程（已清理）：

- `0b436b05-c6e5-40b0-825e-75bcbf67da02_0`：xfchat-calendar-test-20260331（已创建→更新→添加参会人→删除参会人→删除）
- `e1e8049e-8c54-4a40-acc2-e5bd9936ed1c_0`：xfchat-calendar-shortcut-test（+create 创建→删除）

User 创建的测试日程（已清理）：

- `db21020d-c8d8-4123-83b3-73d14b0c80f6_0`：xfchat-calendar-user-create-test（+create 创建→删除）
- `34c36879-95dd-4efc-956a-0552e14eb7e4_0`：xfchat-calendar-user-attendee-test（+create 带参会人创建→删除）
- `55d46513-765c-45ba-baa0-2b3e954002a2_0`：xfchat-calendar-user-atomic-test（原子 API 创建→更新→添加参会人→删除参会人→删除）
- `87bc9b5e-7467-4825-b0e7-76abd307c486_0`：xfchat-primary-resolve-test（默认 primary 解析验证→删除）

测试共享日历（已清理）：

- `feishu.cn_kImYnVO3wQNryMIBMSMqxf@group.calendar.feishu.cn`：xfchat-test-shared-calendar（Bot 创建→删除）
- `feishu.cn_hI0LDDDbqqUNGsnf6wG12e@group.calendar.feishu.cn`：xfchat-user-shared-calendar（User 创建→更新→删除）

### 2.3 代码修改

**scope.go — 别名映射**

CLI 代码声明的细粒度 scope 在私有化不存在，已添加别名映射：

```go
var scopeAliases = map[string]string{
    "contact:user.basic_profile:readonly": "contact:user.base:readonly",
    "calendar:calendar.event:read":        "calendar:calendar:readonly",
    "calendar:calendar.free_busy:read":    "calendar:calendar:readonly",
}
```

修改文件：`internal/auth/scope.go`

**helpers.go — primary 日历 ID 解析**

`+agenda` 和 `+create` 在未指定 `--calendar-id` 时原本传 `"primary"` 字符串，私有化不识别。已修改为先调用 `POST /open-apis/calendar/v4/calendars/primary` 接口解析为实际日历 ID。

修改文件：`shortcuts/calendar/helpers.go`、`shortcuts/calendar/calendar_agenda.go`、`shortcuts/calendar/calendar_create.go`

## 3. Shortcut 覆盖

| 命令 | 场景 | 权限/身份 | 输入 | 输出摘要 | 结论 |
|---|---|---|---|---|---|
| `+agenda` | 今日日程（显式日历 ID） | User, `calendar:calendar:readonly` | `./xfchat_cli calendar +agenda --as user --calendar-id 'feishu.cn_FnVWAps00pOJSzu2pnkCAc@group.calendar.feishu.cn' --format json` | `data=[], count=0` | 通过 |
| `+agenda` | 历史日程（带日期范围） | User, `calendar:calendar:readonly` | `./xfchat_cli calendar +agenda --as user --calendar-id 'feishu.cn_FnVWAps00pOJSzu2pnkCAc@group.calendar.feishu.cn' --start '2026-03-24' --end '2026-03-28' --format json` | 返回 2 条日程（上线功能讨论会、RDG-AI 分享） | 通过 |
| `+agenda` | 默认 primary（代码修复后） | User, `calendar:calendar` | `./xfchat_cli calendar +agenda --as user --format json` | `data=[]`（自动解析为实际日历 ID，今日无日程） | 通过 |
| `+agenda` | 默认 primary 带日期范围 | User, `calendar:calendar` | `./xfchat_cli calendar +agenda --as user --start '2026-03-24' --end '2026-03-28' --format json` | 返回 2 条日程 | 通过 |
| `+agenda` | Bot 日历 | Bot, `calendar:calendar` | `./xfchat_cli calendar +agenda --as bot --calendar-id 'feishu.cn_cn2LpP9WiCZuwxhMzlczUd@group.calendar.feishu.cn' --format json` | 返回 1 条 Updated Test Event | 通过 |
| `+agenda` | Bot 默认 primary | Bot, `calendar:calendar` | `./xfchat_cli calendar +agenda --as bot --format json` | 返回 1 条日程 | 通过 |
| `+freebusy` | 今日忙闲（默认当前用户） | User, `calendar:calendar:readonly` | `./xfchat_cli calendar +freebusy --as user --format json` | `data=null`（今日无忙碌时段） | 通过 |
| `+freebusy` | 历史忙闲（带日期范围） | User, `calendar:calendar:readonly` | `./xfchat_cli calendar +freebusy --as user --start '2026-03-24' --end '2026-03-28' --format json` | 返回 1 条忙碌时段（03-26 19:00-21:00） | 通过 |
| `+freebusy` | 指定 user-id | User, `calendar:calendar:readonly` | `./xfchat_cli calendar +freebusy --as user --user-id 'ou_84c79a9284092a69283066a89e549251' --start '2026-03-24' --end '2026-03-28' --format json` | 返回 1 条忙碌时段 | 通过 |
| `+freebusy` | Bot 身份查询用户忙闲 | Bot, `calendar:calendar` | `./xfchat_cli calendar +freebusy --as bot --user-id 'ou_84c79a9284092a69283066a89e549251' --start '2026-03-24' --end '2026-03-28' --format json` | 返回 1 条忙碌时段 | 通过 |
| `+create` | User 创建日程 | User, `calendar:calendar` | `./xfchat_cli calendar +create --as user --summary 'xfchat-calendar-user-create-test' --start '2026-04-01T10:00:00+08:00' --end '2026-04-01T11:00:00+08:00' --calendar-id 'feishu.cn_FnVWAps00pOJSzu2pnkCAc@group.calendar.feishu.cn'` | `event_id=db21020d-...` | 通过 |
| `+create` | User 创建日程（带参会人） | User, `calendar:calendar` | `./xfchat_cli calendar +create --as user --summary 'xfchat-calendar-user-attendee-test' --start '2026-04-01T14:00:00+08:00' --end '2026-04-01T15:00:00+08:00' --calendar-id 'feishu.cn_FnVWAps...' --attendee-ids 'ou_84c79a9284092a69283066a89e549251'` | `event_id=34c36879-...` | 通过 |
| `+create` | User 默认 primary（代码修复后） | User, `calendar:calendar` | `./xfchat_cli calendar +create --as user --summary 'xfchat-primary-resolve-test' --start '2026-04-01T16:00:00+08:00' --end '2026-04-01T17:00:00+08:00'` | `event_id=87bc9b5e-...`（自动解析主日历） | 通过 |
| `+create` | Bot 创建日程（带参会人） | Bot, `calendar:calendar` | `./xfchat_cli calendar +create --as bot --summary 'xfchat-calendar-shortcut-test' --start '2026-04-01T14:00:00+08:00' --end '2026-04-01T15:00:00+08:00' --calendar-id 'feishu.cn_cn2LpP9...' --attendee-ids 'ou_84c79a9284092a69283066a89e549251'` | `event_id=e1e8049e-...` | 通过 |
| `+suggestion` | User 推荐空闲时段 | User, `calendar:calendar` | `./xfchat_cli calendar +suggestion --as user --format json` | HTTP 404 | 私有化不支持 |
| `+suggestion` | Bot 推荐空闲时段 | Bot, `calendar:calendar` | `./xfchat_cli calendar +suggestion --as bot --attendee-ids 'ou_84c79a...' --start '2026-04-01T00:00:00+08:00' --end '2026-04-01T23:59:59+08:00' --duration-minutes 30 --format json` | HTTP 404 | 私有化不支持 |

## 4. 原子 API 覆盖

### 4.1 Calendars

| 方法 | 场景 | 权限/身份 | 输入 | 输出摘要 | 结论 |
|---|---|---|---|---|---|
| `calendars.primary` | 查询用户主日历 | User, `calendar:calendar` | `./xfchat_cli calendar calendars primary --as user --format json` | `calendar_id=feishu.cn_FnVWAps...`，`type=primary`，`role=owner` | 通过 |
| `calendars.primary` | 查询 Bot 主日历 | Bot, `calendar:calendar` | `./xfchat_cli calendar calendars primary --as bot --format json` | `calendar_id=feishu.cn_cn2LpP9...`，`summary=星火派-test` | 通过 |
| `calendars.get` | 用户日历详情 | User, `calendar:calendar` | `./xfchat_cli calendar calendars get --params '{"calendar_id":"feishu.cn_FnVWAps..."}' --as user --format json` | 返回日历名称、权限、角色 | 通过 |
| `calendars.get` | Bot 查用户日历 | Bot, `calendar:calendar` | `./xfchat_cli calendar calendars get --params '{"calendar_id":"feishu.cn_FnVWAps..."}' --as bot --format json` | `role=free_busy_reader`（Bot 对用户日历仅忙闲可见） | 通过 |
| `calendars.list` | 用户日历列表 | User, `calendar:calendar` | `./xfchat_cli calendar calendars list --as user --format json` | 返回 1 个主日历，`has_more=false` | 通过 |
| `calendars.list` | Bot 日历列表 | Bot, `calendar:calendar` | `./xfchat_cli calendar calendars list --as bot --format json` | 返回 1 个主日历 + 多个 shared 日历 | 通过 |
| `calendars.search` | 搜索日历 | User, `calendar:calendar` | `./xfchat_cli calendar calendars search --data '{"query":"付涛"}' --as user --format json` | 返回 1 个主日历 | 通过 |
| `calendars.search` | 搜索日历 | Bot, `calendar:calendar` | `./xfchat_cli calendar calendars search --data '{"query":"付涛"}' --as bot --format json` | `page_token=""`（空结果，Bot 无法搜索用户日历） | 通过 |
| `calendars.create` | 创建共享日历 | Bot, `calendar:calendar` | `./xfchat_cli calendar calendars create --data '{"summary":"xfchat-test-shared-calendar","description":"test"}' --as bot --format json` | `calendar_id=feishu.cn_kImYnVO3...`，`type=shared` | 通过 |
| `calendars.create` | 创建共享日历 | User, `calendar:calendar` | `./xfchat_cli calendar calendars create --data '{"summary":"xfchat-user-shared-calendar","description":"test shared calendar by user"}' --as user --format json` | `calendar_id=feishu.cn_hI0LDDD...`，`type=shared` | 通过 |
| `calendars.patch` | 更新日历描述 | Bot, `calendar:calendar` | `./xfchat_cli calendar calendars patch --params '{"calendar_id":"feishu.cn_cn2LpP9..."}' --data '{"description":"updated by CLI atomic api test"}' --as bot --format json` | 返回更新后的日历信息 | 通过 |
| `calendars.patch` | 更新日历描述 | User, `calendar:calendar` | `./xfchat_cli calendar calendars patch --params '{"calendar_id":"feishu.cn_hI0LDDD..."}' --data '{"description":"patched by user"}' --as user --format json` | 返回更新后的日历信息 | 通过 |
| `calendars.delete` | 删除共享日历 | Bot, `calendar:calendar` | `./xfchat_cli calendar calendars delete --params '{"calendar_id":"feishu.cn_kImYnVO3..."}' --as bot --format json` | `code=0` | 通过 |
| `calendars.delete` | 删除共享日历 | User, `calendar:calendar` | `./xfchat_cli calendar calendars delete --params '{"calendar_id":"feishu.cn_hI0LDDD..."}' --as user --format json` | `code=0` | 通过 |

### 4.2 Events

| 方法 | 场景 | 权限/身份 | 输入 | 输出摘要 | 结论 |
|---|---|---|---|---|---|
| `events.instance_view` | 今日日程实例（User 日历） | User, `calendar:calendar` | `./xfchat_cli calendar events instance_view --params '{"calendar_id":"feishu.cn_FnVWAps...","start_time":"1743350400","end_time":"1743436800"}' --as user --format json` | `data={}` (今日无日程) | 通过 |
| `events.instance_view` | 今日日程实例（Bot 日历） | Bot, `calendar:calendar` | `./xfchat_cli calendar events instance_view --params '{"calendar_id":"feishu.cn_cn2LpP9...","start_time":"1743350400","end_time":"1743436800"}' --as bot --format json` | `data={}` | 通过 |
| `events.get` | 获取日程详情 | User, `calendar:calendar` | `./xfchat_cli calendar events get --params '{"calendar_id":"feishu.cn_FnVWAps...","event_id":"a9173edb-..."}' --as user --format json` | 返回完整日程详情（时间、地点、提醒等） | 通过 |
| `events.get` | Bot 读取创建的日程 | Bot, `calendar:calendar` | `./xfchat_cli calendar events get --params '{"calendar_id":"feishu.cn_cn2LpP9...","event_id":"0b436b05-..."}' --as bot --format json` | 返回完整日程详情 | 通过 |
| `events.search` | 搜索日程 | User, `calendar:calendar` | `./xfchat_cli calendar events search --params '{"calendar_id":"feishu.cn_FnVWAps..."}' --data '{"query":"上线功能讨论会"}' --as user --format json` | 返回 1 条匹配日程 | 通过 |
| `events.search` | 搜索日程 | Bot, `calendar:calendar` | `./xfchat_cli calendar events search --params '{"calendar_id":"feishu.cn_cn2LpP9..."}' --data '{"query":"test"}' --as bot --format json` | 返回 1 条 Updated Test Event | 通过 |
| `events.create` | 创建日程 | Bot, `calendar:calendar` | `./xfchat_cli calendar events create --params '{"calendar_id":"feishu.cn_cn2LpP9..."}' --data '{"summary":"xfchat-calendar-test-20260331",...}' --as bot --format json` | `event_id=0b436b05-...`，`status=confirmed` | 通过 |
| `events.create` | 创建日程 | User, `calendar:calendar` | `./xfchat_cli calendar events create --params '{"calendar_id":"feishu.cn_FnVWAps..."}' --data '{"summary":"xfchat-calendar-user-atomic-test",...}' --as user --format json` | `event_id=55d46513-...`，`event_organizer=付涛 taofu5` | 通过 |
| `events.patch` | 更新日程 | Bot, `calendar:calendar` | `./xfchat_cli calendar events patch --params '{"calendar_id":"feishu.cn_cn2LpP9...","event_id":"0b436b05-..."}' --data '{"summary":"...-updated","description":"updated via atomic api"}' --as bot --format json` | 返回更新后的日程 | 通过 |
| `events.patch` | 更新日程 | User, `calendar:calendar` | `./xfchat_cli calendar events patch --params '{"calendar_id":"feishu.cn_FnVWAps...","event_id":"55d46513-..."}' --data '{"summary":"...-updated","description":"updated by user"}' --as user --format json` | 返回更新后的日程 | 通过 |
| `events.delete` | 删除日程 | Bot, `calendar:calendar` | `./xfchat_cli calendar events delete --params '{"calendar_id":"feishu.cn_cn2LpP9...","event_id":"0b436b05-...","need_notification":"false"}' --as bot --format json` | `code=0` | 通过 |
| `events.delete` | 删除日程 | User, `calendar:calendar` | `./xfchat_cli calendar events delete --params '{"calendar_id":"feishu.cn_FnVWAps...","event_id":"55d46513-...","need_notification":"false"}' --as user --format json` | `code=0` | 通过 |

### 4.3 Event Attendees

| 方法 | 场景 | 权限/身份 | 输入 | 输出摘要 | 结论 |
|---|---|---|---|---|---|
| `event.attendees.create` | 添加参会人 | Bot, `calendar:calendar` | `./xfchat_cli calendar event.attendees create --params '{"calendar_id":"feishu.cn_cn2LpP9...","event_id":"0b436b05-...","user_id_type":"open_id"}' --data '{"attendees":[{"type":"user","user_id":"ou_84c79a..."}],"need_notification":true}' --as bot --format json` | 返回参会人 `付涛 taofu5`，`rsvp_status=needs_action` | 通过 |
| `event.attendees.create` | 添加参会人 | User, `calendar:calendar` | `./xfchat_cli calendar event.attendees create --params '{"calendar_id":"feishu.cn_FnVWAps...","event_id":"55d46513-...","user_id_type":"open_id"}' --data '{"attendees":[{"type":"user","user_id":"ou_84c79a..."}],"need_notification":false}' --as user --format json` | 返回参会人 `付涛 taofu5`，`is_organizer=true`，`rsvp_status=accept` | 通过 |
| `event.attendees.list` | 查询参会人列表（Bot 日程） | Bot, `calendar:calendar` | `./xfchat_cli calendar event.attendees list --params '{"calendar_id":"feishu.cn_cn2LpP9...","event_id":"0b436b05-...","user_id_type":"open_id"}' --as bot --format json` | 返回 1 个参会人 | 通过 |
| `event.attendees.list` | 查询参会人列表（User 日程，大量参会人） | User, `calendar:calendar` | `./xfchat_cli calendar event.attendees list --params '{"calendar_id":"feishu.cn_FnVWAps...","event_id":"4d5bb05f-...","user_id_type":"open_id","page_size":"50"}' --as user --format json` | 返回 50 个参会人，`has_more=true`，含分页 token | 通过 |
| `event.attendees.list` | 查询参会人（小日程） | User, `calendar:calendar` | `./xfchat_cli calendar event.attendees list --params '{"calendar_id":"feishu.cn_FnVWAps...","event_id":"a9173edb-...","user_id_type":"open_id"}' --as user --format json` | 返回 2 个参会人（靳驰骋、付涛） | 通过 |
| `event.attendees.batch_delete` | 删除参会人 | Bot, `calendar:calendar` | `./xfchat_cli calendar event.attendees batch_delete --params '{"calendar_id":"feishu.cn_cn2LpP9...","event_id":"0b436b05-...","user_id_type":"open_id"}' --data '{"attendee_ids":["user_76021261..."],"need_notification":false}' --as bot --format json` | `code=0` | 通过 |
| `event.attendees.batch_delete` | 删除参会人 | User, `calendar:calendar` | `./xfchat_cli calendar event.attendees batch_delete --params '{"calendar_id":"feishu.cn_FnVWAps...","event_id":"55d46513-...","user_id_type":"open_id"}' --data '{"attendee_ids":["user_76021261..."],"need_notification":false}' --as user --format json` | `code=0` | 通过 |

### 4.4 Freebusys

| 方法 | 场景 | 权限/身份 | 输入 | 输出摘要 | 结论 |
|---|---|---|---|---|---|
| `freebusys.list` | 查询忙闲（User） | User, `calendar:calendar` | `./xfchat_cli calendar freebusys list --data '{"time_min":"2026-03-31T00:00:00+08:00","time_max":"2026-03-31T23:59:59+08:00","user_id":"ou_84c79a..."}' --as user --format json` | `code=0, data={}` (今日无忙碌) | 通过 |
| `freebusys.list` | 查询忙闲（Bot） | Bot, `calendar:calendar` | `./xfchat_cli calendar freebusys list --data '{"time_min":"2026-03-31T00:00:00+08:00","time_max":"2026-03-31T23:59:59+08:00","user_id":"ou_84c79a..."}' --as bot --format json` | `code=0, data={}` | 通过 |

## 5. 私有化差异、CLI 问题与结论

### 5.1 私有化环境已确认差异

- **`+suggestion` 接口不存在**：`POST /open-apis/calendar/v4/freebusy/suggestion` 在 User 和 Bot 身份下均返回 HTTP 404
- **Bot 搜索用户日历**：`calendars.search` Bot 身份搜索用户名返回空结果（权限隔离正常）

### 5.2 Scope 粒度差异（已修复）

CLI 代码声明细粒度 scope（如 `calendar:calendar.event:read`），私有化环境仅注册粗粒度 scope：

| CLI 声明的细粒度 scope | scope.json 中实际 scope | 映射方式 |
|---|---|---|
| `calendar:calendar.event:read` | `calendar:calendar:readonly` | 别名映射（scope.go） |
| `calendar:calendar.free_busy:read` | `calendar:calendar:readonly` | 别名映射（scope.go） |
| `calendar:calendar.event:create` | `calendar:calendar` | coveredByParent 前缀匹配 |
| `calendar:calendar.event:update` | `calendar:calendar` | coveredByParent 前缀匹配 |
| `calendar:calendar.event:delete` | `calendar:calendar` | coveredByParent 前缀匹配 |

修改文件：`internal/auth/scope.go`

### 5.3 CLI 兼容问题（已修复）

- **`+agenda`/`+create` 默认 `primary` 关键字**：当 `--calendar-id` 未指定时，代码原本将 `calendarId` 设为 `"primary"` 字符串直接传入 API，私有化返回 `190002 invalid parameters`。已修改为先调用 `POST /open-apis/calendar/v4/calendars/primary` 接口解析为实际日历 ID。修改文件：`shortcuts/calendar/helpers.go`、`shortcuts/calendar/calendar_agenda.go`、`shortcuts/calendar/calendar_create.go`

### 5.4 accounts 端点配置问题（历史，已修复）

`~/.xfchat_cli/config.json` 中 `accounts` 端点原配置为 `https://open.xfchat.iflytek.com`，导致 OAuth 授权 URL 域名错误，报错 20029。

- **修复**：将 `accounts` 端点改为 `https://accounts.xfchat.iflytek.com`
- **判定**：CLI 问题（配置问题）

## 6. 覆盖完成度

### 6.1 已覆盖情况

**Shortcuts（4 个）：**

| Shortcut | User 身份 | Bot 身份 | 备注 |
|---|---|---|---|
| `+agenda` | 通过 | 通过 | 默认 primary 已修复 |
| `+freebusy` | 通过 | 通过 | |
| `+create` | 通过 | 通过 | 含参会人、默认 primary 场景 |
| `+suggestion` | 私有化不支持 | 私有化不支持 | HTTP 404 |

**原子 API（4 资源 × 多方法）：**

| 资源.方法 | User 身份 | Bot 身份 | 备注 |
|---|---|---|---|
| `calendars.primary` | 通过 | 通过 | |
| `calendars.get` | 通过 | 通过 | |
| `calendars.list` | 通过 | 通过 | |
| `calendars.search` | 通过 | 通过 | Bot 搜索用户名返回空（权限隔离正常） |
| `calendars.create` | 通过 | 通过 | |
| `calendars.delete` | 通过 | 通过 | |
| `calendars.patch` | 通过 | 通过 | |
| `events.instance_view` | 通过 | 通过 | |
| `events.get` | 通过 | 通过 | |
| `events.search` | 通过 | 通过 | |
| `events.create` | 通过 | 通过 | |
| `events.patch` | 通过 | 通过 | |
| `events.delete` | 通过 | 通过 | |
| `event.attendees.create` | 通过 | 通过 | |
| `event.attendees.list` | 通过 | 通过 | 含分页验证 |
| `event.attendees.batch_delete` | 通过 | 通过 | |
| `freebusys.list` | 通过 | 通过 | |

### 6.2 无法打通的项

- **`+suggestion`**
  `POST /open-apis/calendar/v4/freebusy/suggestion` 在私有化不存在（User/Bot 均 HTTP 404），无法通过权限或代码修改解决

### 6.3 总结

- **Bot 全链路已打通**：查日历、查日程、创建日程、更新日程、删除日程、管理参会人、查忙闲、搜索日程
- **User 全链路已打通**：查日历、查日程、创建日程、更新日程、删除日程、管理参会人、查忙闲、搜索日程、创建/删除共享日历
- **已识别并修复的问题**：scope 粒度别名映射（scope.go）、accounts 端点配置、`primary` 关键字自动解析（helpers.go）
- **私有化不支持**：`+suggestion`（/freebusy/suggestion 接口 404）
- Calendar 模块当前可以按"User 与 Bot 均可完整操作日历全链路"的模式在 i 讯飞私有化环境下落地
