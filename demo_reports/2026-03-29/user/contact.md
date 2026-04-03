# 用户Contact模块测试

状态：通过

## 前置条件

- 身份： `user`
- 报告开始时状态： `xfchat_cli auth status` 显示当前没有已登录 user

## 测试记录

1. 读取当前用户
命令： `xfchat_cli contact +get-user --as user`
结果： 调用成功，返回 `open_id=ou_4b6bc8c3ea123aa0d47b290828fc7891`

2. 按姓名搜索自己
命令： `xfchat_cli contact +search-user --as user --query 王琦钊`
结果： 调用成功，返回 1 条用户记录，open_id 与当前用户一致

## 创建资源

无

## 复现方式

1. 执行 `+get-user`。
2. 用当前显示名执行 `+search-user`。
3. 对比返回的 `open_id`。