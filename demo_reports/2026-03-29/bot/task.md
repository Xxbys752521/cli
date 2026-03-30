# 机器人Task模块测试

状态：不支持

## 前置条件

- 身份： `bot`

## 测试记录

1. 验证 shortcut 支持情况
命令： `lark-cli task +get-my-tasks --as bot`
结果： CLI 直接拒绝该调用，提示 `--as bot is not supported, this command only supports: user`

## 创建资源

无

## 复现方式

1. 以 `bot` 身份执行任意 user 侧 task shortcut。
2. 确认 CLI 在发起 API 调用前就拒绝该身份。