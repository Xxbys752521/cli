# 用户Task模块测试

状态：通过

## 前置条件

- 身份： `user`
- 报告开始时状态： `lark-cli auth status` 显示当前没有已登录 user

## 测试记录

1. 创建任务清单
命令： `lark-cli task +tasklist-create --as user --name 'Codex Demo Tasklist 2026-03-29'`
结果： `guid=3ed5e0ff-e638-4374-88d4-b53d1db10fba`

2. 创建任务
命令： `lark-cli task +create --as user --tasklist-id 3ed5e0ff-e638-4374-88d4-b53d1db10fba --summary 'Codex demo follow-up' --description 'Created by lark-cli for reproducible demo testing' --due 2026-03-30`
结果： `guid=71b4b7a0-9536-4927-b83e-7ddfd9e06c2d`

3. 验证 shortcut 查询结果
命令： `lark-cli task +get-my-tasks --as user --query 'Codex demo follow-up'`
结果： 调用成功，但返回 `items=null`

4. 按 guid 直接读取
命令： `lark-cli task tasks get --params '{"task_guid":"71b4b7a0-9536-4927-b83e-7ddfd9e06c2d"}'`
结果： 返回完整任务详情，包含 `todo` 状态、截止日期和任务清单绑定信息

## 创建资源

- 任务清单： `3ed5e0ff-e638-4374-88d4-b53d1db10fba`
- 任务： `71b4b7a0-9536-4927-b83e-7ddfd9e06c2d`

## 复现方式

1. 创建任务清单。
2. 在该清单下创建任务。
3. 如果 shortcut 查询为空，再用原生 `tasks get` 校验任务是否存在。