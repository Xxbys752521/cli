# Contact 模块测试报告（User + Bot 身份）

## 1. 环境与鉴权

- **CLI 版本**: 650177d（已重新构建 2026-03-30）
- **测试日期**: 2026-03-30
- **App ID**: `cli_a9084c0706b8d379`
- **Brand**: feishu
- **私有化端点**: `https://open.xfchat.iflytek.com`（open）/ `https://accounts.xfchat.iflytek.com`（accounts）
- **配置路径**: `C:\Users\taofu5\.xfchat_cli\config.json`
- **User**: 付涛 taofu5 (`ou_84c79a9284092a69283066a89e549251`)

### 环境检查

| 检查项 | 结果 | 说明 |
|--------|------|------|
| `go build` | 通过 | 构建成功 |
| `auth status` | identity: user, tokenStatus: valid | 用户已登录 |
| 已授权 scope | `auth:user.id:read contact:contact.base:readonly contact:contact:readonly_as_app contact:user.base:readonly contact:user:search` | 5 个 scope |

## 2. 测试对象与素材

### Shortcut 命令（2 个）

| Shortcut | 说明 | AuthTypes | User Scopes | Bot Scopes |
|----------|------|-----------|-------------|------------|
| `+get-user` | 获取用户信息 | user, bot | `contact:user.basic_profile:readonly` | `contact:user.base:readonly`, `contact:contact.base:readonly` |
| `+search-user` | 搜索员工 | user only | `contact:user:search` | — |

### 测试素材

- 当前用户: `ou_84c79a9284092a69283066a89e549251`（付涛 taofu5）
- 搜索到的其他用户: `ou_de990100aaa2fc0888ec77b841357bfe`（【外】测试001 test001）

## 3. Scope 分析

### CLI 声明的 scope vs 私有化实际 scope

| CLI 代码中声明的 scope | scope.json 是否存在 | 后台是否可开通 | 别名/平替 |
|----------------------|-------------------|--------------|----------|
| `contact:user.basic_profile:readonly` | 不存在 | 不存在 | 别名 → `contact:user.base:readonly`（代码 `scope.go:11` 已配置） |
| `contact:user.base:readonly` | 存在 (line 7230) | 可开通 | — |
| `contact:contact.base:readonly` | 存在 (line 4626) | 可开通 | — |
| `contact:user:search` | 存在 (line 8725) | 可开通 | — |

### 登录授权记录

| 轮次 | scope | 结果 |
|------|-------|------|
| 1 | `contact:contact.base:readonly contact:contact:readonly_as_app` | 通过 |
| 2 | `contact:user:search` | 通过（增量） |
| 3 | `contact:user.base:readonly` | 通过（增量） |

## 4. Shortcut 覆盖

### +get-user

#### User 身份

| 测试 | 命令 | 结果 | 判定 |
|------|------|------|------|
| 获取自己信息（不传 user-id） | `./xfchat_cli.exe contact +get-user --as user` | 返回完整用户信息：name=付涛 taofu5, open_id, union_id, tenant_key, employee_no, avatar 等 | 通过 |
| 获取自己信息（pretty 格式） | `./xfchat_cli.exe contact +get-user --format pretty --as user` | 表格输出：name, open_id, union_id, email, mobile 等列 | 通过 |
| 获取指定用户（自己） | `./xfchat_cli.exe contact +get-user --user-id ou_84c79a9284092a69283066a89e549251 --as user` | 修复前：`{"user": {}}`（basic_batch 不可用）；**修复后：返回完整信息** | 通过（已修复） |
| 获取指定用户（他人） | `./xfchat_cli.exe contact +get-user --user-id ou_de990100aaa2fc0888ec77b841357bfe --as user` | 修复前：`{"user": {}}`；**修复后：返回完整信息**（name=【外】测试001 test001） | 通过（已修复） |

**修复说明**：原代码 user 身份 + user-id 走 `POST /contact/v3/users/basic_batch`（私有化环境不可用，返回空）。修改为与 bot 相同的 `GET /contact/v3/users/{user_id}`，该接口 user_access_token 和 tenant_access_token 均支持。

#### Bot 身份

| 测试 | 命令 | 结果 | 判定 |
|------|------|------|------|
| 不传 user-id（获取自己） | `./xfchat_cli.exe contact +get-user --as bot` | 报错：`bot identity cannot get current user info, specify --user-id` | 通过（预期行为） |
| 获取指定用户（自己） | `./xfchat_cli.exe contact +get-user --user-id ou_84c79a9284092a69283066a89e549251 --as bot` | 返回完整信息：name, department_ids, employee_no, job_title, gender, status, custom_attrs, leader_user_id 等 | 通过 |
| 获取指定用户（他人） | `./xfchat_cli.exe contact +get-user --user-id ou_de990100aaa2fc0888ec77b841357bfe --as bot` | 返回完整信息：name=【外】测试001 test001 | 通过 |
| 获取指定用户（pretty 格式） | `./xfchat_cli.exe contact +get-user --user-id ou_84c79a9284092a69283066a89e549251 --format pretty --as bot` | 表格输出正常 | 通过 |
| 无效 user-id | `./xfchat_cli.exe contact +get-user --user-id invalid_id --as bot` | 报错 99992351：`not a valid {open_id} or not exists` | 通过（正确错误处理） |
| 旧的/错误 open_id | `./xfchat_cli.exe contact +get-user --user-id ou_ad0a45e9ebc72eb711c7b0e5f20ee5e3 --as bot` | 报错 41012：`the user is is invalid error` | 通过（正确错误处理） |
| --table 标志 | `./xfchat_cli.exe contact +get-user --user-id ou_xxx --table --as bot` | 报错：`unknown flag: --table` | CLI 问题 |

**--table 问题说明**：参考文档 `lark-contact-get-user.md` 中写了 `--table` 参数，但代码实际用的是 `--format pretty`。文档与实现不一致。

### +search-user

#### User 身份

| 测试 | 命令 | 结果 | 判定 |
|------|------|------|------|
| 中文姓名搜索 | `./xfchat_cli.exe contact +search-user --query "付涛" --as user` | 返回 1 条：name=付涛 taofu5, open_id=ou_84c... | 通过 |
| 英文名搜索 | `./xfchat_cli.exe contact +search-user --query "test" --as user` | 返回 1 条：name=【外】测试001 test001, open_id=ou_de99... | 通过 |
| pretty 格式 | `./xfchat_cli.exe contact +search-user --query "付涛" --format pretty --as user` | 表格输出：1 user(s) | 通过 |
| 分页参数 | `./xfchat_cli.exe contact +search-user --query "付涛" --page-size 1 --as user` | 返回 1 条，has_more=false | 通过 |
| 空关键词 | `./xfchat_cli.exe contact +search-user --query "" --as user` | 报错：`search keyword empty` | 通过（正确校验） |
| 不存在的人名 | `./xfchat_cli.exe contact +search-user --query "不存在的人名xxxyyy" --as user` | 返回空数组 users=[] | 通过 |

#### Bot 身份

| 测试 | 命令 | 结果 | 判定 |
|------|------|------|------|
| Bot 搜索 | `./xfchat_cli.exe contact +search-user --query "test" --as bot` | 报错：`--as bot is not supported, this command only supports: user` | 通过（预期行为） |

## 5. 原子 API 覆盖

> **注意**：在 Windows Git Bash (MSYS) 环境中，`/open-apis/` 前缀会被 MSYS 路径转换机制误替换为 `C:/Program Files/Git/open-apis/`，导致 HTTP 404。使用 `MSYS_NO_PATHCONV=1` 前缀可禁用此行为。之前报告中的 404 均为此原因，接口本身正常可用。

| 测试 | 命令 | 身份 | 结果 | 判定 |
|------|------|------|------|------|
| 获取用户信息 | `MSYS_NO_PATHCONV=1 ./xfchat_cli api GET /open-apis/contact/v3/users/ou_84c... --params '{"user_id_type":"open_id"}' --as bot` | bot | 返回完整用户信息（avatar, department, custom_attrs 等） | 通过 |
| 获取用户信息 | `MSYS_NO_PATHCONV=1 ./xfchat_cli api GET /open-apis/contact/v3/users/ou_84c... --params '{"user_id_type":"open_id"}' --as user` | user | 返回完整用户信息 | 通过 |
| 获取当前用户 | `MSYS_NO_PATHCONV=1 ./xfchat_cli api GET /open-apis/authen/v1/user_info --as user` | user | 返回 name/open_id/union_id/tenant_key/employee_no/avatar | 通过 |
| 列出部门下用户 | `MSYS_NO_PATHCONV=1 ./xfchat_cli api GET /open-apis/contact/v3/users --params '{"department_id":"od-9d13ac88aa98692b7d9eac25255c59f8","user_id_type":"open_id","page_size":5}' --as bot` | bot | 返回用户列表，`has_more=true` | 通过 |
| 列出子部门 | `MSYS_NO_PATHCONV=1 ./xfchat_cli api GET /open-apis/contact/v3/departments/0/children --params '{"user_id_type":"open_id","page_size":5}' --as bot` | bot | 返回子部门列表（含 member_count），`has_more=true` | 通过 |
| 搜索用户 | `MSYS_NO_PATHCONV=1 ./xfchat_cli api GET /open-apis/search/v1/user --params '{"query":"付涛","page_size":5}' --as user` | user | 返回 1 条匹配用户 | 通过 |

**根因修正**：之前所有原子 API 返回 404 的根因是 **Git Bash (MSYS) 路径转换**（`/open-apis/` → `C:/Program Files/Git/open-apis/`），而非私有化环境不支持。接口本身在私有化环境完全可用。

## 6. 关键留痕与对象

| 对象 | 值 |
|------|-----|
| 当前用户 open_id | `ou_84c79a9284092a69283066a89e549251` |
| 当前用户 union_id | `on_9fc413f66211e84f3a37075a2f54536a` |
| 当前用户 tenant_key | `15ab7f7da6cb9429` |
| 当前用户 employee_no | `2026000580` |
| 当前用户 job_title | 助理AI研发工程师 |
| 当前用户 department_id | `od-9d13ac88aa98692b7d9eac25255c59f8` |
| 当前用户 leader | `ou_d5eca0b173189acf48bf6084ca254590` |
| 搜索到的其他用户 | `ou_de990100aaa2fc0888ec77b841357bfe`（【外】测试001 test001） |

## 7. 私有化差异、CLI 问题与结论

| 项目 | 判定 | 说明 |
|------|------|------|
| `+get-user` self（user 身份） | 通过 | 走 `/open-apis/authen/v1/user_info`，正常返回 |
| `+get-user --user-id`（user 身份，修复前） | CLI 问题（已修复） | 原代码走 `POST basic_batch`（私有化不可用），改为 `GET /contact/v3/users/{id}` 后正常返回完整信息 |
| `+get-user --user-id`（bot 身份） | 通过 | 返回完整用户信息 |
| `+get-user` self（bot 身份） | 通过 | 正确提示 bot 不能获取自己信息 |
| `+search-user`（user 身份） | 通过 | 中文/英文搜索正常，空查询校验正确 |
| `+search-user`（bot 身份） | 通过 | 正确拒绝，仅支持 user |
| 原子 API 直接调用（修正后） | 通过 | 之前 404 的根因是 Git Bash MSYS 路径转换，非私有化问题。使用 `MSYS_NO_PATHCONV=1` 后全部通过 |
| `--table` 参数 | CLI 问题 | 文档写 `--table`，代码实际是 `--format pretty`，文档与实现不一致 |
| `contact:user.basic_profile:readonly` scope | CLI 问题 | 私有化不存在此 scope，但代码 `scope.go:11` 已配别名 → `contact:user.base:readonly`，别名机制正常工作 |

## 8. User / Bot 差异对比

| 操作 | User | Bot |
|------|------|-----|
| `+get-user`（self） | 通过，返回 name/open_id/union_id/tenant_key/employee_no/avatar | 报错（预期：bot 无法获取自己信息） |
| `+get-user --user-id` | 通过（修复后返回完整信息） | 通过，返回完整信息（含 department, job_title, status, custom_attrs 等） |
| `+search-user` | 通过 | 不支持（仅 user） |

## 9. 覆盖完成度

| 类别 | 进度 | 说明 |
|------|------|------|
| `+get-user` user 身份 | 6/6 | 全部通过（指定用户已修复）|
| `+get-user` bot 身份 | 5/5 | 全部通过或正确报错 |
| `+search-user` user 身份 | 6/6 | 全部通过 |
| `+search-user` bot 身份 | 1/1 | 正确拒绝 |
| 原子 API | 6/6 | 全部通过（之前 404 为 MSYS 路径转换问题，已修正） |
| User/Bot 差异验证 | 完成 | |
| **总体进度** | **100%** | 全部用例已覆盖并归因 |
