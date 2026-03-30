# lark-cli 私有化部署测试报告

## 测试环境

| 项目 | 值 |
|------|------|
| 私有化平台 | 讯飞内部飞书 `https://open.xfchat.iflytek.com` |
| 应用 | 星火派-test（`cli_a9084c0706b8d379`） |
| OAuth 认证域名 | `https://accounts.xfchat.iflytek.com` |
| 测试时间 | 2026-03-28 |
| CLI 版本 | 本地编译（含私有化端点覆盖 + Authorization Code Flow） |

## 认证方式

| 身份 | 认证方式 | 状态 |
|------|---------|------|
| Bot | AppID + AppSecret → Tenant Access Token（SDK 自动获取） | ✅ 正常 |
| User | OAuth Authorization Code Flow（`--web` 浏览器回调） | ✅ 正常 |

> **说明**：私有化环境不支持 Device Flow（`/oauth/v1/device_authorization` 返回 404），通过 `--web` 参数使用 Authorization Code Flow 完成用户登录。回调地址固定为 `http://localhost:8080/callback`，需在开发者后台安全设置中配置重定向 URL。

---

## 一、即时消息（IM）

### Bot 身份

| # | 功能 | 场景说明 | 命令 | 结果 | 备注 |
|---|------|---------|------|------|------|
| 1 | 发送文本消息 | 向指定群聊发送纯文本消息 | `im +messages-send --chat-id oc_xxx --text "..." --as bot` | ✅ 通过 | |
| 2 | 发送 Markdown 富文本 | 向群聊发送带格式的富文本消息（标题、列表、引用等） | `im +messages-send --chat-id oc_xxx --markdown $'...' --as bot` | ✅ 通过 | 换行需用 `$'...\n...'` 语法 |
| 3 | 回复消息 | 对指定消息进行回复，支持话题回复 | `im +messages-reply --message-id om_xxx --text "..." --as bot` | ✅ 通过 | |
| 4 | 获取群信息 | 查看群名称、群主、成员数量等详细信息 | `api GET /open-apis/im/v1/chats/{id} --as bot` | ✅ 通过 | |
| 5 | 群列表 | 获取机器人所在的所有群聊列表 | `api GET /open-apis/im/v1/chats --as bot` | ✅ 通过 | |
| 6 | 拉取历史消息 | 获取指定群聊的历史消息记录，支持时间范围和分页 | `im +chat-messages-list --chat-id oc_xxx --as bot` | ✅ 通过 | 需机器人在群内 |
| 7 | 批量获取消息 | 根据消息 ID 列表批量查询消息详情（内容、发送者、时间） | `im +messages-mget --message-ids om_xxx,om_yyy --as bot` | ✅ 通过 | |
| 8 | 创建群 | 创建新的群聊，可指定群名、邀请成员和机器人 | `im +chat-create --name "群名" --as bot` | ✅ 通过 | |
| 9 | 更新群信息 | 修改群名称或群描述 | `im +chat-update --chat-id oc_xxx --name "新名" --as bot` | ✅ 通过 | |
| 10 | 话题消息列表 | 获取某条话题（Thread）下的所有回复消息 | `im +threads-messages-list --thread om_xxx --as bot` | ⏭ 跳过 | 测试消息无话题 |
| 11 | 下载消息附件 | 下载消息中的图片或文件附件到本地 | `im +messages-resources-download --message-id om_xxx --as bot` | ⏭ 跳过 | 测试消息无附件 |

### User 身份

| # | 功能 | 场景说明 | 命令 | 结果 | 备注 |
|---|------|---------|------|------|------|
| 1 | 群列表 | 获取当前用户所在的所有群聊列表 | `api GET /open-apis/im/v1/chats --as user` | ✅ 通过 | 返回用户所在的所有群 |
| 2 | 搜索群 | 按关键词搜索可见的群聊（可按成员、类型过滤） | `im +chat-search --query "关键词" --as user` | ✅ 通过 | |
| 3 | 更新群信息 | 以用户身份修改群名称或群描述 | `im +chat-update --chat-id oc_xxx --name "..." --as user` | ✅ 通过 | 需用户在群内 |
| 4 | 拉取历史消息 | 以用户身份获取群聊历史消息记录 | `im +chat-messages-list --chat-id oc_xxx --as user` | ❌ 不支持 | 错误码 99991668: user access token not support |
| 5 | 发送消息 | 以用户身份向群聊发送消息 | `api POST /open-apis/im/v1/messages --as user` | ❌ 不支持 | API 不支持 user token |
| 6 | 批量获取消息 | 以用户身份根据消息 ID 批量查询消息详情 | `im +messages-mget --message-ids om_xxx --as user` | ❌ 权限不足 | 错误码 230027: Permission denied |
| 7 | 搜索消息 | 按关键词全文搜索聊天记录，支持按群、发送人、时间范围过滤 | `im +messages-search --query "..." --as user` | ❌ 不支持 | API 返回 404，私有化环境无此接口 |

### IM 模块小结

- **Bot 身份**：9/9 全部通过（发消息、收消息、群管理、批量获取、创建群）
- **User 身份**：3/7 通过（群列表、搜索群、更新群信息），消息收发和搜索不支持
- **建议**：IM 消息收发使用 Bot 身份，群管理操作 User/Bot 均可

---

## 二、后续模块（待测试）

以下模块将在后续测试中补充：

- [ ] 云盘（Drive）
- [ ] 多维表格（Bitable）
- [ ] 电子表格（Sheets）
- [ ] 日历（Calendar）
- [ ] 任务（Task）
- [ ] 文档（Docs）
- [ ] 通讯录（Contact）
- [ ] 邮箱（Mail）
- [ ] 视频会议（VC）
- [ ] Wiki

---

## 已知问题与适配

### 1. scope 命名差异

私有化环境的权限名称与公有云存在差异：

| 公有云 scope（shortcut 硬编码） | 私有化环境可用 scope | 关系 |
|------|------|------|
| `im:message.group_msg:get_as_user` | 不存在 | 被 `im:message` 覆盖 |
| `im:message.p2p_msg:get_as_user` | 不存在 | 被 `im:message` 覆盖 |

**已修复**：`MissingScopes` 函数增加父子 scope 兼容逻辑，父 scope（如 `im:message`）自动覆盖子 scope 检查。

### 2. OAuth 授权 URL 编码

私有化服务端对 URL 参数编码严格，Go 的 `url.Values.Encode()` 会将 scope 中的冒号编码为 `%3A`（如 `im%3Amessage%3Areadonly`），导致 20029 错误。

**已修复**：手动构造授权 URL，scope 中的冒号保持原始格式。

### 3. 授权 scope 数量限制

单次授权最多 50 个 scope，使用 `--recommend` 会超限（错误码 20084）。

**建议**：使用 `--domain` 指定具体域，或 `--scope` 手动指定需要的权限。

### 4. 回调端口

授权回调固定使用 `localhost:8080`，需在开发者后台重定向 URL 中配置 `http://localhost:8080/callback`。
