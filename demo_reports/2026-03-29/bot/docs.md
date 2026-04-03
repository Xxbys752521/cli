# 机器人Docs模块测试

状态：通过

## 前置条件

- 身份： `bot`

## 测试记录

1. 创建文档
命令： `xfchat_cli docs +create --as bot --title 'Codex Demo Doc 2026-03-29' --markdown '# Demo Report\n\n- created by bot\n- date: 2026-03-29\n- purpose: cli test'`
结果： `doc_id=AcqCdfl7dogRRIxwuB7cp3xWnHc`

2. 追加内容
命令： `xfchat_cli docs +update --as bot --doc AcqCdfl7dogRRIxwuB7cp3xWnHc --mode append --markdown '\n## Update\n- append from bot\n- step: docs update'`
结果： 更新成功

3. 读取文档
命令： `xfchat_cli docs +fetch --as bot --doc AcqCdfl7dogRRIxwuB7cp3xWnHc`
结果： 成功读取 Markdown，长度为 `127`，内容包含追加段落

## 创建资源

- 文档： `AcqCdfl7dogRRIxwuB7cp3xWnHc`
- URL: `https://www.feishu.cn/docx/AcqCdfl7dogRRIxwuB7cp3xWnHc`

## 复现方式

1. 用 Markdown 创建文档。
2. 通过 `--mode append` 追加一段内容。
3. 读取文档并确认前后两部分内容都在。