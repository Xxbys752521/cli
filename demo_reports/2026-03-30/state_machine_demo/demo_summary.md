# xfchat_cli 状态机 Demo 结果摘要

## 结果

已用当前 `xfchat_cli` 真实跑通一轮“答疑状态机”最小闭环，覆盖了两条链路：

- 自动答疑链路：
  - 提问
  - 补充信息
  - 检索知识库
  - 命中相似案例
  - 自动回答
  - 用户满意
  - 更新答疑结果

- 转人工链路：
  - 提问
  - 检索未命中
  - 转人工
  - 人工处理
  - 回写知识库
  - 更新答疑结果

## 已创建资源

- Demo 群：
  - 名称：`状态机Demo群-20260330_121533`
  - chat_id：`oc_0a491c7b930824ef3e0f7f4995ca1dd9`
  - 分享链接：[打开群聊](https://applink.feishu.cn/client/chat/chatter/add_by_link?link_token=a35oc694-49d7-4d8a-ace9-1728b9cf2dda)

- Demo 知识库文档：
  - 标题：`状态机Demo知识库-20260330_121533`
  - 文档链接：[打开文档](https://www.feishu.cn/docx/MUlrdNBH0o7ymPxYl7zcEcj5nqd)

## 最终状态

### CASE-A

- 提问人：`模拟提问人A`
- 问题：`私有化环境中点击知识库提示无权限`
- 是否命中案例：`是`
- 最终状态：`已结束`
- 处理结论：`命中私有化环境知识库权限案例，已给出授权与环境校验建议。`

### CASE-B

- 提问人：`模拟提问人B`
- 问题：`新员工入职后，外部知识库链接打开空白页，重新登录也不行`
- 是否命中案例：`否`
- 最终状态：`已结束`
- 处理结论：`人工确认是知识库可见范围缺失叠加旧租户缓存导致，补权限并清缓存后恢复。`

## 当前 CLI 的真实边界

- 群消息发送目前是 bot-only，所以 demo 中“提问人”和“人工”发言都是 bot 用模拟角色发出的
- `base` 的字段创建存在平台级限流，连续重复跑脚本时容易碰到 `OpenAPIAddField limited`
- 这说明当前 `xfchat_cli` 可以支撑状态机闭环，但如果要做成稳定可重复的产品能力，Base 建表和字段初始化需要加重试和退避

## 文件

- 脚本：[run_lark_cli_state_machine_demo.sh](/Users/wangqizhao/Developer/iflytek/cli/demo_reports/2026-03-30/state_machine_demo/run_lark_cli_state_machine_demo.sh)
- 摘要：[demo_summary.md](/Users/wangqizhao/Developer/iflytek/cli/demo_reports/2026-03-30/state_machine_demo/demo_summary.md)
- 原始报告：[demo_run_report.md](/Users/wangqizhao/Developer/iflytek/cli/demo_reports/2026-03-30/state_machine_demo/demo_run_report.md)
