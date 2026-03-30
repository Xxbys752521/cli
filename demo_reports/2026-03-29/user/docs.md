# 用户Docs模块测试

状态：通过

## 前置条件

- 身份： `user`
- 报告开始时状态： `lark-cli auth status` 显示当前没有已登录 user

## 测试记录

1. 创建挂在知识库下的文档
命令： `lark-cli docs +create --as user --wiki-space my_library --title 'Codex Wiki Demo 2026-03-29' --markdown '# Wiki Demo\n\nCreated from lark-cli as user.'`
结果： `doc_id=HHsbdDBROo3kboxdfrPcdZLinYd`, returned wiki URL

2. 读取内容
命令： `lark-cli docs +fetch --as user --doc HHsbdDBROo3kboxdfrPcdZLinYd`
结果： 成功读取 Markdown，标题和正文与预期一致

## 创建资源

- 文档： `HHsbdDBROo3kboxdfrPcdZLinYd`
- Wiki URL token: `UhyQwyTvtiMteskyqH5cqVMJnVs`

## 复现方式

1. 在 `my_library` 下创建文档。
2. Fetch by returned `doc_id`.
3. 确认回读的 Markdown 与写入内容一致。