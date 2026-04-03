# 机器人Wiki模块测试

状态：通过

## 前置条件

- 身份： `bot`

## 测试记录

1. 用实际 docx token 解析知识库节点
命令： `xfchat_cli wiki spaces get_node --as bot --params '{"token":"HHsbdDBROo3kboxdfrPcdZLinYd","obj_type":"docx"}'`
结果： success, returned `node_token=UhyQwyTvtiMteskyqH5cqVMJnVs`, `space_id=7593008347729628385`, `title='Codex Wiki Demo 2026-03-29'`

## 创建资源

- 解析出的节点： `UhyQwyTvtiMteskyqH5cqVMJnVs`

## 复现方式

1. 先创建一个挂在知识库下的文档。
2. Call `wiki spaces get_node` with the actual `docx` token and `obj_type=docx`.