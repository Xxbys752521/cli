# 机器人VC模块测试

状态：不支持

## 前置条件

- 身份： `bot`

## 测试记录

1. Search support check
命令： `lark-cli vc +search --as bot --start 2026-03-22 --end 2026-03-29`
结果： CLI 直接拒绝该调用，提示 `--as bot is not supported, this command only supports: user`

## 创建资源

无

## 复现方式

1. 以 `bot` 身份执行上述命令。
2. Confirm the shortcut is user-only.