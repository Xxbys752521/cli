# 机器人Minutes模块测试

状态：失败

## 前置条件

- 身份： `bot`

## 测试记录

1. Token 格式校验
命令： `xfchat_cli minutes minutes get --as bot --params '{"minute_token":"dummy"}'`
结果： failed with `code=99992402`, field validation reported `minute_token` min length `24`

## 创建资源

无

## 复现方式

1. Run the command above with an intentionally short token.
2. 确认 CLI 能把 API 字段校验细节透出。