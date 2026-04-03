# xfchat_cli 私有化能力总览

基于 `docs/result` 的最终结果，只保留场景结论。

当前总览只统计运行态已启用能力；`mail` 当前默认禁用，不纳入本矩阵。

## IM 即时通讯
- 能做：Bot 在群聊/私聊发文本、Markdown、图片、文件；做回复、线程、转发、下载、Pin、Reaction、拉人进群。
- 不能稳定做：User 读群历史、群消息搜索、按 `--user-id` 直接拉私聊历史。
- 佐证：`docs/result/im.md`

## Base 多维表格
- 能做：创建 Base、表、字段、记录、视图；搭业务样板；打通高级权限、角色、协作者、表单。
- 不能稳定做：`record-history-list`、`data-query`、地理位置写值；`Email` 字段会回落成文本。
- 佐证：`docs/result/base.md`，[排期样板](https://yf2ljykclb.xfchat.iflytek.com/base/basrz4CISgYO2yc3n4kfuESDLWb)

## Sheets 云文档表格
- 能做：创建表格、读写单元格、追加数据、查找内容、多工作表管理、筛选、协作者。
- 不能稳定做：图表、`filter_view` 还未测完；写后立刻读偶发短暂最终一致性。
- 佐证：`docs/result/sheets.md`，[User 表格](https://yf2ljykclb.xfchat.iflytek.com/sheets/shtrzlFRtdfQZ2o8qauJpuPO8Rf)

## Docs 云文档
- 能做：OpenAPI 创建/读取/追加/覆盖文档；批量写块；写 Heading、列表、引用、代码、Todo、Divider、Table、Sheet、图片、文件、Callout、Grid；写原生链接、`@用户`、`@文档`。
- 不能稳定做：`replace_range`、`replace_all`、`insert_before`、`insert_after`、`delete_range`；评论创建 404；复杂块读取不保真。
- 佐证：`docs/result/doc.md`，[全块演示](https://yf2ljykclb.xfchat.iflytek.com/docx/doxrz4AeBEYjGOoi06sHM7Dc8Kd)

## Wiki 知识库
- 能做：读取节点、解析 wiki token 对应的真实对象类型和 `obj_token`、列出当前 User 可见知识空间。
- 不能稳定做：创建知识空间、把云文档正式挂成新知识库节点、完整的 create/move/member 命令面。
- 佐证：`docs/result/wiki.md`，[知识库节点](https://yf2ljykclb.xfchat.iflytek.com/wiki/WG34wO7Nei7MIkkykdCrscUIz1t)

## Contact 通讯录
- 能做：User/Bot 获取指定用户信息；User 获取自己信息；User 搜索员工；原子 API 查用户、列部门成员、列子部门。
- 不能稳定做：Bot 不支持“获取自己信息”；`+search-user` 仅支持 User；Windows Git Bash 下原子 API 需禁用 MSYS 路径转换。
- 佐证：`docs/result/contact.md`

## Calendar 日历
- 能做：查日程、忙闲、主日历、日历列表；创建/更新/删除日程；创建共享日历；添加/列出/删除参会人。
- 不能稳定做：`+suggestion` 推荐空闲时段在私有化里是接口 404。
- 佐证：`docs/result/calendar.md`

## Task 任务
- 能做：User 查看我的任务；Bot 创建/更新/评论/提醒/完成/重开任务；创建子任务、清单、成员协作。
- 不能稳定做：`tasks.list` 和 `+get-my-tasks` 只支持 User，不支持 Bot。
- 佐证：`docs/result/task.md`

## 当前判断
- 已完成主链路并有演示留痕：IM、Base、Sheets、Docs、Calendar、Task。
- 已完成读取与解析但未完成建设型能力：Wiki。
- 已完成读取与搜索主链路：Contact。
