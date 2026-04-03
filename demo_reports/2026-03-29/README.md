# xfchat_cli 演示测试总览

日期：2026-03-29  
仓库：[Xxbys752521/cli](https://github.com/Xxbys752521/cli)（上游参考：[larksuite/cli](https://github.com/larksuite/cli/tree/main)）  
工作目录：`/Users/wangqizhao/Developer/iflytek/cli`

> 本目录统一以当前私有化 fork `xfchat_cli` 为准。若材料里讨论“私有化 lark-cli”，指的就是当前这套发行版；二进制名与配置目录统一为 `xfchat_cli` / `~/.xfchat_cli`。

## 一、这份材料的用途

这份目录下的文档用于沉淀一套可复现的 `xfchat_cli` 演示测试记录，按模块和身份拆开留痕，便于后续做三件事：

1. 对外演示时快速查阅成功案例和失败边界。
2. 对内评估私有化飞书环境下哪些能力已经可用，哪些能力仍有缺口。
3. 为后续和 `spark` 做范式对比、做迁移方案提供事实依据。

## 二、目录结构

- `bot/*.md`：以 `bot` 身份执行的模块测试记录
- `user/*.md`：以 `user` 身份执行的模块测试记录
- `composite/cases.md`：跨模块组合案例
- `artifacts/`：测试过程中生成的本地文件
- `cli_vs_spark_comparison_report.md`：`spark` 与 `xfchat_cli` 对比报告
- `migration_technical_plan.md`：迁移技术方案

## 三、当前环境结论

- App 配置已完成。
- `bot` 身份可用。
- `user` 身份在测试过程中已完成重新授权。
- 本轮测试以“真实命令执行 + 实际返回结果”为准，没有使用伪造结果。

## 四、覆盖情况

| 模块 | Bot | User | 说明 |
| --- | --- | --- | --- |
| Calendar | 失败 | 通过 | `bot` 侧返回缺少 scope，`user` 侧创建和读取通过 |
| IM | 通过 | 部分通过 | `bot` 发消息和读消息通过，`user` 侧搜索和读取返回空结果 |
| Docs | 通过 | 通过 | `bot` 文档和 `user` 知识库文档都已创建并回读 |
| Drive | 通过 | 通过 | 两种身份都完成上传，`bot` 下载并完成内容校验 |
| Base | 通过 | 通过 | `bot` 建库、配字段、写记录通过，`user` 能读 |
| Sheets | 通过 | 通过 | `bot` 创建和写入通过，`user` 能读 |
| Task | 不支持 | 通过 | 相关 shortcut 以 `user` 为主，`bot` 被命令层拒绝 |
| Wiki | 通过 | 通过 | 两种身份都能解析节点信息 |
| Contact | 不支持 | 通过 | `user` 查询通过，`bot` 不支持对应 shortcut |
| Mail | 已禁用 | 已禁用 | 当前私有化运行时不注册 `mail` 服务命令，仅保留历史 skill 文档 |
| VC | 不支持 | 通过 | `user` 查询通过但为空，`bot` 不支持对应 shortcut |
| Minutes | 失败 | 失败 | `user` 缺少 minutes scope，`bot` 只验证到参数格式 |

## 五、本轮创建的关键资源

- Bot 文档：`AcqCdfl7dogRRIxwuB7cp3xWnHc`
- User 知识库文档：`HHsbdDBROo3kboxdfrPcdZLinYd`
- Bot 电子表格：`SqSssOzYDh0l9Atatcwco6ffnyd`
- Bot Base：`Gfi7bVQTQaR1sLsR18YcfsfpnjJ`
- User 任务清单：`3ed5e0ff-e638-4374-88d4-b53d1db10fba`
- User 任务：`71b4b7a0-9536-4927-b83e-7ddfd9e06c2d`
- User 日程：`f1557220-4fcc-4bc7-853b-48f23a21f9f2_0`
- Bot IM 会话：`oc_42591f2bf31defb34a548c339962af5f`

## 六、演示时值得强调的点

### 1. 身份差异比命令名更重要

同一个模块下，`bot` 和 `user` 的实际可用能力并不对称。  
这不是 `xfchat_cli` 独有问题，而是飞书权限模型本身决定的，但 `xfchat_cli` 把这个差异暴露得比较清楚。

### 2. 失败场景同样有价值

这套记录里保留了几类失败：

- 缺少 scope
- 私有化环境未开放接口
- 命令只支持某一身份
- 账号本身没有对应业务能力

这类失败在汇报时不要回避，因为它们更能说明平台边界和真实落地成本。

### 3. 私有化飞书不一定能达到标准飞书的演示效果

这是本轮材料必须单独说明的风险点。

当前私有化飞书与标准飞书存在两类明显差异：

1. **MCP 能力缺失或不完整**
   - 文档类 shortcut 里有一部分依赖 MCP。
   - 如果私有化环境没有部署对应 MCP 服务，演示效果会明显弱于标准飞书。
   - 这不是命令本身没有设计出来，而是运行环境缺件。

2. **User 身份能力受限**
   - 私有化环境下，部分接口对 user token 支持不完整，或者权限口径与公有云不一致。
   - 这会直接影响 demo 是否能做到“用户态全链路演示”。

这两个问题都已经在你本地适配和测试记录里出现过，不是理论风险。

## 七、如何复现

1. 先确认 `xfchat_cli --help` 可执行。
2. 执行 `xfchat_cli auth status` 确认当前身份状态。
3. 需要 `user` 侧测试时，先完成 `auth login`。
4. 逐个进入 `bot/*.md` 和 `user/*.md`，按文档中的命令重放。
5. 用本文件里的资源 ID 做跨模块串联。

## 八、相关文档

- 组合案例：[cases.md](/Users/wangqizhao/Developer/iflytek/cli/demo_reports/2026-03-29/composite/cases.md)
- 对比报告：[cli_vs_spark_comparison_report.md](/Users/wangqizhao/Developer/iflytek/cli/demo_reports/2026-03-29/cli_vs_spark_comparison_report.md)
- 迁移方案：[migration_technical_plan.md](/Users/wangqizhao/Developer/iflytek/cli/demo_reports/2026-03-29/migration_technical_plan.md)
