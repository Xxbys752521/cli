# 组合能力案例

状态：通过

## 一、用途

这份文档只记录跨模块、跨身份的组合案例。  
单模块测试请看 `bot/*.md` 和 `user/*.md`。

## 二、候选案例

1. 通讯录查询 -> 发消息 -> 建文档 -> 上传附件
2. 建日程 -> 写会议纪要 -> 建跟进任务
3. 建 Base -> 写 Sheet -> 生成文档摘要
4. 查会议记录 -> 取妙记 -> 提取待办

## 三、已执行案例

### 案例 1：Bot 内容链路

流程：
`docs +create -> docs +update -> docs +fetch -> drive +upload -> drive +download -> cmp`

结果：
整条链路通过。已经覆盖文档创建、文档更新、文档回读、文件上传、文件下载、本地一致性校验。

### 案例 2：User 计划链路

流程：
`contact +get-user -> calendar +create -> calendar +agenda -> task +tasklist-create -> task +create -> task tasks get`

结果：
整条链路通过。已经覆盖用户身份识别、日程创建、日程回读、任务清单创建、任务创建和任务详情查询。

### 案例 3：跨身份协同链路

流程：
`sheets +create/write (bot) -> sheets +read (user) -> base +base-create/field-create/record-upsert (bot) -> base +base-get (user) -> docs +create --wiki-space (user) -> wiki spaces get_node (bot/user)`

结果：
整条链路通过。说明当前租户里部分资源可以在不同身份间协作读取，这一点适合做演示。

### 案例 4：边界与诊断链路

流程：
`calendar +agenda (bot fail) -> mail +triage (bot fail) -> mail profile (user fail) -> minutes get (user fail) -> task/vc/contact user-only reject on bot`

结果：
这组命令的价值不在成功，而在把几类真实边界说明白：

- 私有化环境下 scope 可能并不对齐
- user 身份并不一定具备完整业务能力
- 一部分 shortcut 本身就是单身份设计
- demo 是否能做完整，不只取决于命令是否存在，也取决于环境是否补齐

## 四、复现建议

1. 演示时先讲一个完整通过的案例，再讲一个失败边界案例。
2. 通过的案例建议优先用 Docs、Drive、Base、Sheets、Calendar、Task。
3. 边界案例建议用 Mail、Minutes、Bot Calendar 或 user-only shortcut。
