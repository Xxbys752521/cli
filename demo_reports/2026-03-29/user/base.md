# 用户Base模块测试

状态：通过

## 前置条件

- 身份： `user`
- 报告开始时状态： `xfchat_cli auth status` 显示当前没有已登录 user

## 测试记录

1. 读取 bot 创建的 Base
命令： `xfchat_cli base +base-get --as user --base-token Gfi7bVQTQaR1sLsR18YcfsfpnjJ`
结果： 调用成功，返回了 `Codex Demo Base 2026-03-29` 的元信息

## 创建资源

本模块未以 user 身份新建资源

## 复现方式

1. 先以 `bot` 身份创建 Base。
2. 再切换到 `user` 身份。
3. 调用 `+base-get`，确认该 Base 可见。