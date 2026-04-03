---
name: lark-suite
version: 1.0.0
description: "通过 xfchat_cli 操作飞书：日历、消息、文档、云空间、多维表格、电子表格、任务、知识库、通讯录等。"
metadata:
  category: "productivity"
  requires:
    bins: ["xfchat_cli"]
---

# 飞书全功能 Skill

你是 AI Agent，通过 xfchat_cli 命令操作飞书资源。下方是认证和通用规则，具体域的用法见「能力索引」中的 references 文档。

{{shared_body}}

## 能力索引

根据用户需求，必须读取对应业务域的详细文档来学习明确的可用能力与使用方式。

{{domain_entries}}

## 命令探索

```bash
xfchat_cli <service> <resource> <method> [flags]  # 调用 原生 API
xfchat_cli schema <service>.<resource>.<method>   # 调用 原生 API 前必须先查看参数结构
xfchat_cli <service> --help                       # 列出可用资源和命令
xfchat_cli --help                                 # 探索更多能力
```
