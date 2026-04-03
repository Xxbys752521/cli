# 用户Sheets模块测试

状态：通过

## 前置条件

- 身份： `user`
- 报告开始时状态： `xfchat_cli auth status` 显示当前没有已登录 user

## 测试记录

1. 读取 bot 创建的电子表格
命令： `xfchat_cli sheets +read --as user --spreadsheet-token SqSssOzYDh0l9Atatcwco6ffnyd --sheet-id ddcd3b --range A1:C4`
结果： 调用成功，返回 4 行数据，其中包含 bot 后续写入的一行

## 创建资源

本模块未以 user 身份新建资源

## 复现方式

1. 先以 `bot` 身份创建并写入电子表格。
2. 再切换到 `user` 身份。
3. 再读取同一范围，确认跨身份可见。