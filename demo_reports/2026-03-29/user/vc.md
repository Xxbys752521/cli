# 用户VC模块测试

状态：通过

## 前置条件

- 身份： `user`
- 报告开始时状态： `lark-cli auth status` 显示当前没有已登录 user

## 测试记录

1. 查询会议记录
命令： `lark-cli vc +search --as user --start 2026-03-22 --end 2026-03-29`
结果： 调用成功，结果为 `total=0`

## 创建资源

无

## 复现方式

1. 提供一个有效时间范围。
2. 执行 `vc +search`。
3. Record whether the result is empty or populated.