# 机器人Mail模块测试

状态：失败

## 前置条件

- 身份： `bot`

## 测试记录

1. 邮件摘要查询
命令： `xfchat_cli mail +triage --as bot --max 3 --format json`
结果： failed with `code=99991672`, required scope `mail:user_mailbox.message:readonly`

## 创建资源

无

## 复现方式

1. 以 `bot` 身份执行上述命令。
2. 确认 CLI 返回权限错误以及对应的 `console_url`。