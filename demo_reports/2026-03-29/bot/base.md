# 机器人Base模块测试

状态：通过

## 前置条件

- 身份： `bot`

## 测试记录

1. 创建 Base
命令： `xfchat_cli base +base-create --as bot --name 'Codex Demo Base 2026-03-29' --time-zone Asia/Shanghai`
结果： `base_token=Gfi7bVQTQaR1sLsR18YcfsfpnjJ`

2. 第一次建表尝试
命令： `xfchat_cli base +table-create ... --fields '[{\"field_name\":\"module\",\"type\":1},...]'`
结果： 失败，原因是 shortcut 期望现代字段 JSON，不接受旧版数字枚举字段格式

3. 创建数据表
命令： `xfchat_cli base +table-create --as bot --base-token Gfi7bVQTQaR1sLsR18YcfsfpnjJ --name DemoTable --fields '[{\"type\":\"text\",\"name\":\"module\"}]'`
结果： 第一次重试提示 `DemoTable` 已存在，对应表 ID 为 `tblOsV9eJFNn27Nd`

4. 按 `lark-base` 规范补字段
命令：
- `xfchat_cli base +field-create --as bot --base-token Gfi7bVQTQaR1sLsR18YcfsfpnjJ --table-id tblOsV9eJFNn27Nd --json '{"type":"text","name":"module","style":{"type":"plain"}}'`
- `... --json '{"type":"text","name":"status","style":{"type":"plain"}}'`
- `... --json '{"type":"text","name":"identity","style":{"type":"plain"}}'`
结果： 字段均已创建；其中一次因加字段限流，在短暂等待后重试成功

5. 写入并读取记录
命令：
- `xfchat_cli base +record-upsert --as bot --base-token Gfi7bVQTQaR1sLsR18YcfsfpnjJ --table-id tblOsV9eJFNn27Nd --json '{"module":"base","status":"ok","identity":"bot"}'`
- `xfchat_cli base +record-list --as bot --base-token Gfi7bVQTQaR1sLsR18YcfsfpnjJ --table-id tblOsV9eJFNn27Nd`
结果： record created as `recvfdQuXNOYwa`, list output showed `["bot","NO.001","base","ok"]`

## 创建资源

- Base： `Gfi7bVQTQaR1sLsR18YcfsfpnjJ`
- 数据表： `tblOsV9eJFNn27Nd`
- 记录： `recvfdQuXNOYwa`

## 复现方式

1. 创建一个 Base。
2. 写记录前先查看表和字段结构。
3. 使用 `+field-create` 时应传 `type/name/style` 结构，不要再用旧版数字字段类型。
4. 写入记录后再用 `+record-list` 回读。