# 状态机 Demo

## 文件

- 脚本：[run_lark_cli_state_machine_demo.sh](/Users/wangqizhao/Developer/iflytek/cli/demo_reports/2026-03-30/state_machine_demo/run_lark_cli_state_machine_demo.sh)
- 运行报告：[demo_run_report.md](/Users/wangqizhao/Developer/iflytek/cli/demo_reports/2026-03-30/state_machine_demo/demo_run_report.md)

## 目标

用当前 `lark-cli` 跑一遍“答疑状态机”最小闭环，包括：

- 建群
- 构造聊天记录
- 建答疑表
- 建知识库文档
- 跑自动答疑链路
- 跑转人工回写知识库链路

## 说明

为了保证可复现，聊天记录中的“提问人”和“人工”发言都由 bot 以“模拟角色”的方式发送。  
这不是权限缺陷，而是当前 `lark-cli` 的发消息能力本来就是 bot-only，用户侧消息发送不在 CLI 范围内。
