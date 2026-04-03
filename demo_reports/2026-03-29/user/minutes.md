# 用户Minutes模块测试

状态：失败

## 前置条件

- 身份： `user`
- 报告开始时状态： `xfchat_cli auth status` 显示当前没有已登录 user

## 测试记录

1. 检查 Minutes 权限
命令： `xfchat_cli minutes minutes get --params '{"minute_token":"dummy"}'`
结果： CLI stopped earlier with `insufficient permissions`, required scope `minutes:minutes:readonly`

## 创建资源

无

## 复现方式

1. 先为 user 身份补齐 `minutes:minutes:readonly` 授权。
2. Retry with a real minute token from a meeting artifact.