# 用户Calendar模块测试

状态：通过

## 前置条件

- 身份： `user`
- 报告开始时状态： `xfchat_cli auth status` 显示当前没有已登录 user

## 测试记录

1. 先读取当日日程
命令： `xfchat_cli calendar +agenda --as user --start 2026-03-29T00:00:00+08:00 --end 2026-03-29T23:59:59+08:00`
结果： `data=[]`

2. 创建日程
命令： `xfchat_cli calendar +create --as user --summary 'Codex Demo Event 2026-03-29' --description 'Created by xfchat_cli for reproducible demo testing' --start 2026-03-29T15:00:00+08:00 --end 2026-03-29T15:30:00+08:00`
结果： `event_id=f1557220-4fcc-4bc7-853b-48f23a21f9f2_0`

3. 读取日程 window
命令： `xfchat_cli calendar +agenda --as user --start 2026-03-29T14:30:00+08:00 --end 2026-03-29T16:00:00+08:00`
结果： 返回 1 条事件，主题、组织者和时间均与创建结果一致

## 创建资源

- Event: `f1557220-4fcc-4bc7-853b-48f23a21f9f2_0`

## 复现方式

1. 先确认 user 身份已完成授权。
2. 创建一个短时日程。
3. 读取对应时间窗口并确认日程已出现。