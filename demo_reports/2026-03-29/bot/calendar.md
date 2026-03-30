# 机器人Calendar模块测试

状态：失败

## 前置条件

- 身份： `bot`
- 权限状态：需通过实际命令验证

## 测试记录

1. 读取日程
命令： `lark-cli calendar +agenda --as bot --start 2026-03-29T00:00:00+08:00 --end 2026-03-29T23:59:59+08:00`
结果： failed with `code=99991672`, required scope `calendar:calendar.event:read`
结论： 当前环境下，这套 app 配置尚不能支撑 bot 身份读取日程。

## 创建资源

无

## 复现方式

1. 以 `bot` 身份执行上述命令。
2. 确认 CLI 返回权限错误以及对应的 `console_url`。