# 用户Wiki模块测试

状态：通过

## 前置条件

- 身份： `user`
- 报告开始时状态： `lark-cli auth status` 显示当前没有已登录 user

## 测试记录

1. 解析已创建的知识库文档
命令： `lark-cli wiki spaces get_node --as user --params '{"token":"HHsbdDBROo3kboxdfrPcdZLinYd","obj_type":"docx"}'`
结果： success, returned `node_token=UhyQwyTvtiMteskyqH5cqVMJnVs`, `space_id=7593008347729628385`, `obj_type=docx`

## 创建资源

- 节点： `UhyQwyTvtiMteskyqH5cqVMJnVs`

## 复现方式

1. Create a wiki-backed doc.
2. 用返回的 `doc_id` 作为 `token`，并传 `obj_type=docx`。
3. 确认返回中的 `space_id`、`node_token` 和 `obj_token`。