# 机器人Contact模块测试

状态：不支持

## 前置条件

- 身份： `bot`

## 测试记录

1. 查询当前用户
命令： `lark-cli contact +get-user --as bot`
结果： 校验失败，bot 身份如果不显式传 `--user-id`，无法读取当前用户信息

2. 验证 shortcut 是否支持
命令： `lark-cli contact +search-user --as bot --query 王琦钊`
结果： CLI 直接拒绝该身份，提示 `--as bot is not supported, this command only supports: user`

## 创建资源

无

## 复现方式

1. 以 `bot` 身份执行上述两条命令。
2. 确认 CLI 能区分“参数校验失败”和“该身份不支持”两类情况。