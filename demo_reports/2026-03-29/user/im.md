# 用户IM模块测试

状态：部分通过

## 前置条件

- 身份： `user`
- 报告开始时状态： `lark-cli auth status` 显示当前没有已登录 user

## 测试记录

1. 搜索群聊
命令： `lark-cli im +chat-search --as user --query Codex`
结果： 调用成功，结果为 `total=0`

2. 读取与自己的 P2P 会话消息
命令： `lark-cli im +chat-messages-list --as user --user-id ou_4b6bc8c3ea123aa0d47b290828fc7891 --start 2026-03-29T11:00:00+08:00 --end 2026-03-29T11:30:00+08:00`
结果： 调用成功，结果为 `total=0`

## 创建资源

无

## 复现方式

1. 对群聊执行一次搜索。
2. 对 user 的 P2P 会话执行 `+chat-messages-list`。
3. 记录当前租户下结果是为空还是有数据。