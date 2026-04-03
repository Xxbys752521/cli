# 机器人IM模块测试

状态：通过

## 前置条件

- 身份： `bot`

## 测试记录

1. 给当前已登录用户发私信
命令： `xfchat_cli im +messages-send --as bot --user-id ou_4b6bc8c3ea123aa0d47b290828fc7891 --text 'Codex demo ping from bot at 2026-03-29 11:10'`
结果： `message_id=om_x100b53bd317214a0c2dbc64e498673b`, `chat_id=oc_42591f2bf31defb34a548c339962af5f`

2. 读取聊天记录
命令： `xfchat_cli im +chat-messages-list --as bot --chat-id oc_42591f2bf31defb34a548c339962af5f --start 2026-03-29T11:00:00+08:00 --end 2026-03-29T11:30:00+08:00`
结果： 返回 1 条文本消息，`message_id` 与发送结果一致

## 创建资源

- 会话： `oc_42591f2bf31defb34a548c339962af5f`
- 消息： `om_x100b53bd317214a0c2dbc64e498673b`

## 复现方式

1. 将 `--user-id` 替换为可达的用户 open_id。
2. 以 `bot` 身份发送一条文本消息。
3. 再通过 `+chat-messages-list` 读回消息。