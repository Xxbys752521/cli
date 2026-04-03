# 机器人Sheets模块测试

状态：通过

## 前置条件

- 身份： `bot`

## 测试记录

1. 创建电子表格
命令： `xfchat_cli sheets +create --as bot --title 'Codex Demo Sheet 2026-03-29' --headers '["module","status","identity"]' --data '[["docs","pending","bot"],["sheets","pending","bot"]]'`
结果： `spreadsheet_token=SqSssOzYDh0l9Atatcwco6ffnyd`

2. 读取元信息
命令： `xfchat_cli sheets +info --as bot --spreadsheet-token SqSssOzYDh0l9Atatcwco6ffnyd`
结果： 默认工作表 ID 为 `ddcd3b`

3. 读取初始数据行
命令： `xfchat_cli sheets +read --as bot --spreadsheet-token SqSssOzYDh0l9Atatcwco6ffnyd --sheet-id ddcd3b --range A1:C3`
结果： 返回表头和两行初始数据

4. 验证 range 参数校验
命令： `xfchat_cli sheets +write --as bot --spreadsheet-token SqSssOzYDh0l9Atatcwco6ffnyd --range A4:C4 --values '[["drive","pending","bot"]]'`
结果： 被拒绝，原因是缺少 `--sheet-id` 或 `<sheetId>!` 前缀

5. 修正后写入
命令： `xfchat_cli sheets +write --as bot --spreadsheet-token SqSssOzYDh0l9Atatcwco6ffnyd --sheet-id ddcd3b --range A4:C4 --values '[["drive","pending","bot"]]'`
结果： 成功更新 `3` 个单元格，revision 为 `2`

## 创建资源

- 表格： `SqSssOzYDh0l9Atatcwco6ffnyd`
- 工作表： `ddcd3b`

## 复现方式

1. 创建带表头和初始数据的电子表格。
2. Call `+info` to get `sheet_id`.
3. 通过 `+read` 读取数据行。
4. Write new rows with an explicit `sheet_id`.