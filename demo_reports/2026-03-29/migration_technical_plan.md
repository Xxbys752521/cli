# 基于 `xfchat_cli` 的迁移技术方案

日期：2026-03-29

## 一、方案定位

这份方案不是讨论“要不要把 `spark` 推倒重来”，而是讨论一件更具体的事：

**是否把飞书能力层从项目内分散封装，逐步收敛到 `xfchat_cli` 这一统一能力层。**

如果做这件事，目标不是换技术名词，而是解决几个长期问题：

- 飞书能力重复封装
- 私有化适配重复做
- Agent 需要背过多接口细节
- 多机器人之间很难共享同一套能力

## 二、迁移目标

迁移后的目标形态是：

1. 飞书原子能力统一由 `xfchat_cli` 承接。
2. `spark` 继续保留业务编排、业务规则、内部工具和产品逻辑。
3. Agent 主要负责“任务拆解和能力组合”，不再承担过多飞书接口细节判断。
4. sandbox 只保留在确实需要文件处理、代码执行和动态 Skill 解包的场景。

## 三、迁移范围

### 1. 优先迁移的能力

优先迁移到 `xfchat_cli` 的是飞书开放平台原子能力：

- IM
- Docs / Drive / Wiki
- Calendar
- Task
- Base / Sheets
- Contact
- VC / Minutes
- Mail

这类能力的共同特点是：

- 与飞书 API 一一对应
- 具有稳定输入输出
- 能被多个场景复用

### 2. 暂不迁移的能力

以下能力继续保留在 `spark` 或其他内部工具层：

- 内部数据库能力
- 搜索与检索能力
- 周报、会议提案、群知识等强业务语义模块
- 调度系统本身
- 业务规则和产品态对话逻辑

## 四、生产级鉴权与权限方案

这是迁移能否落地的关键。

### 1. 原则

生产环境里必须做到：

1. 同一套运行时能区分不同用户身份。
2. `bot` 和 `user` 的调用边界清楚。
3. token 不进入 prompt，不在业务代码中散传。
4. 权限不足时能快速定位是环境问题、scope 问题还是身份问题。

### 2. 建议的三层模型

#### A. 应用层

- 每个飞书应用有自己的 `appId/appSecret`
- 这层用于获取 tenant token
- 建议存放在统一密钥系统，不散落在代码和脚本里

#### B. 用户授权层

- 每个用户独立持有自己的 user token / refresh token
- 用户授权关系与应用绑定
- 这层是生产级多用户调用的基础

#### C. 运行时身份解析层

每次执行命令时，先决定本次调用走哪种身份：

- 如果能力只支持 `bot`，直接走应用身份
- 如果能力只支持 `user`，按当前用户上下文取 token
- 如果两种身份都支持，按能力策略或默认规则选择

运行时只向 `xfchat_cli` 注入最终解析后的身份上下文，不把 token 暴露给 Agent。

### 3. 建议的数据模型

#### `app_registry`

- `app_id`
- `brand`
- `endpoint_open`
- `endpoint_accounts`
- `endpoint_mcp`
- `secret_ref`
- `status`

#### `user_auth_binding`

- `user_open_id`
- `app_id`
- `access_token_ref`
- `refresh_token_ref`
- `granted_scope`
- `expires_at`
- `refresh_expires_at`
- `status`

#### `capability_policy`

- `capability_name`
- `allowed_identity`
- `required_scope`
- `fallback_identity`
- `fallback_behavior`

### 4. 生产实现建议

开发机上可以继续用本地 `config.json + keychain`。  
生产环境不建议直接照搬本地模式，建议做两层包装：

1. 统一配置与 token 存储服务
2. `xfchat_cli` 运行前由 wrapper 注入对应配置视图

也就是说，保留 `xfchat_cli` 的身份语义和权限模型，替换其底层存储实现。

### 5. 权限处理策略

权限判断不要留给 Agent 猜。

建议做法：

1. 每个能力声明自己的 required scopes。
2. 执行前做 preflight 检查。
3. 若失败，返回结构化错误：
   - 缺少哪些 scope
   - 当前走的是 `bot` 还是 `user`
   - 是否支持切换身份
   - 是否需要重新授权

这样 Agent 看到的是“结构化结论”，不是一堆 OAuth 细节。

## 五、现有 tool / skill 的迁移方式

### 1. 工具分三类处理

#### 第一类：飞书原子能力

迁到 `xfchat_cli`

例如：

- 发消息
- 建文档
- 上传文件
- 建日程
- 建任务
- 建 Base / 写 Sheet

#### 第二类：业务组合能力

继续保留在 `spark`，但底层改为调用 `xfchat_cli`

例如：

- 建会 -> 写纪要 -> 建任务
- 群摘要 -> 发消息 -> 落文档
- 会议信息提取 -> 回填日程/任务

这类能力仍属于业务层，不建议直接下沉到 CLI。

#### 第三类：非飞书能力

保持独立工具形态

例如：

- 搜索
- 数据库访问
- 内部平台 API
- 周报和知识处理逻辑

### 2. Skill 的迁移建议

如果一个 Skill 的本质是“教 Agent 如何使用飞书能力”，建议迁到 `xfchat_cli skills` 体系。  
如果一个 Skill 的本质是“封装某个业务流程”，建议保留在业务层。

### 3. ToolFactory 的收敛方向

迁移后，`spark-pai` 的 tool 面建议收敛成三组：

- `FeishuCapabilityToolset`：统一代理到 `xfchat_cli`
- `BusinessToolset`：保留业务工具
- `SandboxToolset`：保留少量受控执行能力

这样 Agent 看到的工具面会更小，更稳定。

## 六、Sandbox 是否保留

结论：**保留，但缩边界。**

### 1. 不建议继续由 sandbox 承载的能力

以下场景不建议继续走 sandbox：

- 纯飞书 API 调用
- 文档、任务、日程、表格等标准开放能力
- 仅仅为了执行 HTTP 请求而引入 shell 或脚本

原因很直接：

- 增加冷启动
- 增加失败面
- 权限边界更模糊
- Agent 更容易做多余推理

### 2. 建议保留 sandbox 的场景

保留在这些场景：

- 文件加工
- 代码执行
- 动态 Skill 解包和隔离运行
- 需要临时工作目录的复杂处理

简单说：

- 飞书能力交给 `xfchat_cli`
- 文件和代码类能力交给 sandbox

## 七、迁移路径

### Phase 0：能力盘点

先把现有 `spark` 里的飞书相关 service / tool / skill 做一张映射表：

- 名称
- 当前用途
- 是否属于飞书原子能力
- 是否已有 `xfchat_cli` 对应命令
- 当前身份要求
- 是否依赖 sandbox

这一步只盘点，不改代码。

### Phase 1：加一层 `spark -> xfchat_cli` 适配

先不要上来就重写 Agent。  
先做 adapter：

- 对外仍保留现有工具名
- 对内改为执行 `xfchat_cli` 对应命令

这样可以先切走底层飞书调用面，不影响上层业务逻辑。

### Phase 2：把身份和权限治理独立出来

建设统一的：

- app registry
- user auth binding
- capability policy
- token refresh service
- permission preflight

这一步完成后，Agent 不再碰 token 细节。

### Phase 3：逐步下线旧飞书封装

按模块下线旧 service / tool：

- 先 IM、Docs、Drive
- 再 Calendar、Task
- 再 Base、Sheets、Wiki

每下线一类，就把 ownership 明确迁到 `xfchat_cli`。

### Phase 4：收敛 Skill 和 Prompt

当底层调用面稳定后，再同步收敛：

- Skill 说明
- Agent prompt 中的工具说明
- 业务流程中的调用方式

否则会出现新旧两套工具并存，增加决策噪音。

## 八、预期收益

### 1. 对 Agent 的收益

#### 能否减少推理压力

可以。

原因不是模型换了，而是让模型少管了很多接口细节：

迁移前，Agent 要同时决定：

- 调哪个 tool
- tool 名是什么
- 参数怎么拼
- 走 `bot` 还是 `user`
- 有无权限
- 是否需要 sandbox

迁移后，Agent 主要决定：

- 需要哪几个原子动作
- 这些动作的顺序是什么

这会明显减少不必要的推理。

#### 能否提升速度

通常可以，但要分开看：

- **首响应速度**：更容易提升，因为权限和身份问题能更早暴露
- **总耗时**：在多模块任务里有机会下降，但前提是同时缩减 tool 面和 sandbox 使用

保守判断：

- 平均迭代轮数会下降
- 失败重试会减少
- 多模块任务整体时延有机会下降 20% 到 50%

### 2. 对能力组合的收益

迁移后，Agent 更容易拥有跨模块原子能力。

例如：

- `contact +get-user -> im +messages-send -> docs +create -> drive +upload`
- `calendar +create -> docs +create -> task +create`
- `base +record-upsert -> sheets +write -> docs +update`

这类链路的好处是：

- 原子度高
- 错误边界清楚
- 更容易复用
- 更适合做演示和自动化

### 3. 对工程和组织的收益

- 飞书能力不再重复封装
- 私有化适配集中处理
- 多个机器人可共享同一能力层
- 人、脚本、Agent 都能用同一套接口

## 九、风险与约束

### 1. 最大现实风险：私有化飞书与标准飞书不一致

这个风险必须单独强调。

当前私有化飞书与标准飞书至少有两类明显差异：

#### A. MCP 支持缺失或不完整

- 部分 docs 类 shortcut 依赖 MCP。
- 如果私有化环境没有对应 MCP 服务，标准飞书里能展示的效果，私有化环境不一定能做出来。
- 这会直接影响 demo 的完整程度。

#### B. User 身份能力受限

- 私有化环境下，部分接口对 user token 支持不完整。
- 还有一些 scope 名称、授权方式和公有云不一致。
- 结果就是：命令层虽然设计完整，但 user 侧演示链路可能做不到标准飞书那样顺。

这一点已经在你本地适配和测试里反复出现，不能忽略。

### 2. 不能把环境限制误判成方案问题

如果后续做 demo，要先区分三件事：

1. `xfchat_cli` 是否具备该能力
2. 私有化环境是否开放该能力
3. 当前账号和身份是否具备权限

很多 demo 做不出来，并不是因为范式不对，而是因为环境能力没补齐。

### 3. 新旧调用面并存的过渡风险

如果只新增 `xfchat_cli`，但不逐步下线旧飞书封装，会出现双轨并存。  
这样会导致：

- 维护成本上升
- 排障路径变长
- Agent 决策噪音增加

### 4. 生产 token 存储必须重做

开发机模式不能直接搬到生产。  
必须把 token 和 config 托管化，否则后续多用户治理会失控。

## 十、建议的落地顺序

建议顺序如下：

1. 先盘点能力
2. 再做 adapter
3. 同时建设统一鉴权层
4. 然后按模块下线旧飞书封装
5. 最后收敛 Skill 和 Prompt

这样风险最小，也最容易分阶段验收。

## 十一、一句话结论

如果目标只是维持现有机器人继续运行，现有架构还能继续做。  
如果目标是让飞书能力变成多个机器人、多个 Agent、多个场景都能共用的底座，那么把飞书能力层迁到 `xfchat_cli` 范式，是更值得投入的方向。
