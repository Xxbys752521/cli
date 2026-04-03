# 云文档表格（Sheets）测试结果

## 最终结论

- 云文档表格模块不走 MCP，当前 CLI 的 `sheets` 能力全部直接调用 OpenAPI。
- 私有化环境下可用的主路径为：
  - `POST /open-apis/sheets/v3/spreadsheets`
  - `GET /open-apis/sheets/v3/spreadsheets/{spreadsheet_token}`
  - `GET /open-apis/sheets/v3/spreadsheets/{spreadsheet_token}/sheets/query`
  - `GET /open-apis/sheets/v2/spreadsheets/{spreadsheet_token}/values/{range}`
  - `PUT /open-apis/sheets/v2/spreadsheets/{spreadsheet_token}/values`
  - `POST /open-apis/sheets/v2/spreadsheets/{spreadsheet_token}/values_append`
  - `POST /open-apis/sheets/v3/spreadsheets/{spreadsheet_token}/sheets/{sheet_id}/find`
- 已完成的代码适配是把 CLI 本地 scope 预检查从旧口径改到私有化 `scope.json` 里真实存在的口径：
  - 写能力：`sheets:spreadsheet`
  - 读能力：`sheets:spreadsheet:readonly`

## 代码适配

- 已更新 [shortcuts/sheets/sheet_create.go](/Users/wangqizhao/Developer/iflytek/cli/shortcuts/sheets/sheet_create.go)
- 已更新 [shortcuts/sheets/sheet_info.go](/Users/wangqizhao/Developer/iflytek/cli/shortcuts/sheets/sheet_info.go)
- 已更新 [shortcuts/sheets/sheet_read.go](/Users/wangqizhao/Developer/iflytek/cli/shortcuts/sheets/sheet_read.go)
- 已更新 [shortcuts/sheets/sheet_write.go](/Users/wangqizhao/Developer/iflytek/cli/shortcuts/sheets/sheet_write.go)
- 已更新 [shortcuts/sheets/sheet_append.go](/Users/wangqizhao/Developer/iflytek/cli/shortcuts/sheets/sheet_append.go)
- 已更新 [shortcuts/sheets/sheet_find.go](/Users/wangqizhao/Developer/iflytek/cli/shortcuts/sheets/sheet_find.go)

## Bot 侧测试结果

### 已打通

- `sheets +create`
- `sheets +info`
- `sheets +write`
- `sheets +append`
- `sheets +read`
- `sheets +find`

### 留痕对象

- Spreadsheet Token: `shtrz0wZluX5RNmrRPR4gJ3AW0b`
- URL: <https://yf2ljykclb.xfchat.iflytek.com/sheets/shtrz0wZluX5RNmrRPR4gJ3AW0b>
- Sheet ID: `3ae3aa`
- Title: `xfchat_cli Sheets Bot 测试`

### 关键命令与结果

```bash
./xfchat_cli sheets +create --as bot --title 'xfchat_cli Sheets Bot 测试'
```

返回成功，创建：

- `spreadsheet_token=shtrz0wZluX5RNmrRPR4gJ3AW0b`

```bash
./xfchat_cli sheets +info --as bot --spreadsheet-token shtrz0wZluX5RNmrRPR4gJ3AW0b
```

返回成功，确认：

- `title=xfchat_cli Sheets Bot 测试`
- `sheet_id=3ae3aa`
- `sheet_title=Sheet1`

```bash
./xfchat_cli sheets +write --as bot --spreadsheet-token shtrz0wZluX5RNmrRPR4gJ3AW0b --sheet-id 3ae3aa --range A1 --values '[["模块","状态"],["Sheets","通过"]]'
```

返回成功：

- `updatedRange=3ae3aa!A1:B2`

```bash
./xfchat_cli sheets +append --as bot --spreadsheet-token shtrz0wZluX5RNmrRPR4gJ3AW0b --sheet-id 3ae3aa --range A:B --values '[["OpenAPI","私有化已通"],["读写链路","正常"]]'
```

返回成功：

- `tableRange=3ae3aa!A3:B4`

```bash
./xfchat_cli sheets +read --as bot --spreadsheet-token shtrz0wZluX5RNmrRPR4gJ3AW0b --sheet-id 3ae3aa --range A1:B5
```

返回成功，读取到：

- `模块 | 状态`
- `Sheets | 通过`
- `OpenAPI | 私有化已通`
- `读写链路 | 正常`

```bash
./xfchat_cli sheets +find --as bot --spreadsheet-token shtrz0wZluX5RNmrRPR4gJ3AW0b --sheet-id 3ae3aa --range A1:B10 --find 'OpenAPI'
```

返回成功：

- `matched_cells=["A3"]`

## User 侧测试结果

### 已打通

- `sheets +create`
- `sheets +info`
- `sheets +write`
- `sheets +append`
- `sheets +read`
- `sheets +find`

### 留痕对象

- Spreadsheet Token: `shtrzlFRtdfQZ2o8qauJpuPO8Rf`
- URL: <https://yf2ljykclb.xfchat.iflytek.com/sheets/shtrzlFRtdfQZ2o8qauJpuPO8Rf>
- Sheet ID: `39a262`
- Title: `xfchat_cli Sheets User 测试`

### 关键命令与结果

```bash
./xfchat_cli sheets +create --as user --title 'xfchat_cli Sheets User 测试'
```

返回成功，创建：

- `spreadsheet_token=shtrzlFRtdfQZ2o8qauJpuPO8Rf`

```bash
./xfchat_cli sheets +info --as user --spreadsheet-token shtrzlFRtdfQZ2o8qauJpuPO8Rf
```

返回成功，确认：

- `title=xfchat_cli Sheets User 测试`
- `sheet_id=39a262`
- `sheet_title=Sheet1`

```bash
./xfchat_cli sheets +write --as user --spreadsheet-token shtrzlFRtdfQZ2o8qauJpuPO8Rf --sheet-id 39a262 --range A1 --values '[["模块","状态"],["Sheets","通过"]]'
```

返回成功：

- `updatedRange=39a262!A1:B2`

```bash
./xfchat_cli sheets +append --as user --spreadsheet-token shtrzlFRtdfQZ2o8qauJpuPO8Rf --sheet-id 39a262 --range A:B --values '[["OpenAPI","私有化已通"],["用户态","正常"]]'
```

返回成功：

- `tableRange=39a262!A3:B4`

```bash
./xfchat_cli sheets +read --as user --spreadsheet-token shtrzlFRtdfQZ2o8qauJpuPO8Rf --sheet-id 39a262 --range A1:B6
```

最终返回成功，读取到：

- `模块 | 状态`
- `Sheets | 通过`
- `OpenAPI | 私有化已通`
- `用户态 | 正常`

```bash
./xfchat_cli sheets +find --as user --spreadsheet-token shtrzlFRtdfQZ2o8qauJpuPO8Rf --sheet-id 39a262 --range A1:B10 --find 'OpenAPI'
```

返回成功：

- `matched_cells=["A3"]`

### 一次真实的私有化现象

- `write` 或 `append` 后立即 `read`，有一次读到了旧 revision 的空行。
- 延迟约 2 秒后再次读取，`revision` 与内容均正常。
- 原始 OpenAPI `GET /open-apis/sheets/v2/spreadsheets/{token}/values/{range}` 也返回相同结果。
- 当前判断是这套私有化环境存在短暂最终一致性窗口，不是 CLI 路径或参数错误。

### 原始接口复测结果

```bash
./xfchat_cli api POST /open-apis/sheets/v3/spreadsheets --as user --data '{"title":"raw-sheets-user-20260402-b"}'
```

返回成功：

- `spreadsheet_token=shtrzC9a4ZBISq0mPtT9cVa3yme`

这说明在 user 侧补齐 `sheets:spreadsheet` / `sheets:spreadsheet:readonly` 后，原始建表接口已不再要求额外旧口径。

## 权限结论

### 已确认存在于 `docs/scope.json` 的权限

- `sheets:spreadsheet`
- `sheets:spreadsheet:readonly`

### 私有化旧口径结论

- 早期未授权时，服务端报错中曾出现 `sheets:spreadsheet:create` 与 `drive:drive`
- 重新授权后，原始建表接口已直接成功
- 当前不再需要额外申请这两个旧口径

## 下一步

- 继续扩展 `sheets` 模块时，可优先补：
  - 导出 `+export`
  - 更复杂范围写入
  - 多工作表场景
- 如果后续出现“写后立刻读空”的现象，应按私有化最终一致性处理，优先延迟重读，而不是改接口路径。

## 外围能力补测

### 已打通

- `spreadsheets patch`：可修改表格标题
- `spreadsheet.sheet.filters create`
- `spreadsheet.sheet.filters get`
- `drive permission.members auth`
- `drive permission.members create`

### 关键结果

```bash
./xfchat_cli sheets spreadsheets patch --as user --params '{"spreadsheet_token":"shtrzlFRtdfQZ2o8qauJpuPO8Rf"}' --data '{"title":"xfchat_cli Sheets User 测试-已改名"}'
```

返回成功。约 2 秒后回读：

- `title=xfchat_cli Sheets User 测试-已改名`

```bash
./xfchat_cli sheets spreadsheet.sheet.filters create --as user --params '{"spreadsheet_token":"shtrzlFRtdfQZ2o8qauJpuPO8Rf","sheet_id":"39a262"}' --data '{"range":"39a262!A1:B4","col":"A","condition":{"filter_type":"multiValue","expected":["OpenAPI"]}}'
```

返回成功。

```bash
./xfchat_cli sheets spreadsheet.sheet.filters get --as user --params '{"spreadsheet_token":"shtrzlFRtdfQZ2o8qauJpuPO8Rf","sheet_id":"39a262"}'
```

返回成功，确认：

- `range=39a262!A1:B4`
- `filter_infos[0].col=A`
- `expected=["OpenAPI"]`
- `filtered_out_rows=[2,4]`

```bash
./xfchat_cli drive permission.members auth --as user --params '{"token":"shtrzlFRtdfQZ2o8qauJpuPO8Rf","type":"sheet","action":"view"}'
./xfchat_cli drive permission.members auth --as user --params '{"token":"shtrzlFRtdfQZ2o8qauJpuPO8Rf","type":"sheet","action":"edit"}'
./xfchat_cli drive permission.members auth --as user --params '{"token":"shtrzlFRtdfQZ2o8qauJpuPO8Rf","type":"sheet","action":"manage_public"}'
```

以上均返回：

- `auth_result=true`

```bash
./xfchat_cli drive permission.members create --as user --params '{"token":"shtrzlFRtdfQZ2o8qauJpuPO8Rf","type":"sheet","need_notification":false}' --data '{"member_id":"oc_010df6b42b975ce056cc7c2e717abde8","member_type":"openchat","perm":"view","type":"chat"}'
```

返回成功，已把测试群加入 User 表格协作者。

```bash
./xfchat_cli drive permission.members create --as bot --params '{"token":"shtrz0wZluX5RNmrRPR4gJ3AW0b","type":"sheet","need_notification":false}' --data '{"member_id":"oc_010df6b42b975ce056cc7c2e717abde8","member_type":"openchat","perm":"view","type":"chat"}'
```

返回成功，已把测试群加入 Bot 表格协作者。

### 当前未继续的项

- 图表：本地 registry 与 `scope.json` 中未检到明确的 `chart` 资源或 schema，未盲打。
- 更复杂的工作表属性操作：当前仅依据“操作工作表”文档验证了新增、复制、删除。
- 筛选视图（`filter_view` / `filter_view.condition`）：本地能看到文档线索，但当前未拿到私有化接口文档，也未在 CLI 中完成适配与实测。

### 群内留痕

- `om_x100b53e9132814a0386da5793903a03`

## 多工作表补测

### 已打通

- `addSheet`
- `copySheet`
- `deleteSheet`

### 使用文档

- [操作工作表 - 服务端 API - 开发文档 - 飞书开放平台](/Users/wangqizhao/Developer/iflytek/cli/docs/doc/操作工作表 - 服务端 API - 开发文档 - 飞书开放平台.md)

### 关键命令与结果

```bash
./xfchat_cli api POST /open-apis/sheets/v2/spreadsheets/shtrzlFRtdfQZ2o8qauJpuPO8Rf/sheets_batch_update --as user --data '{"requests":[{"addSheet":{"properties":{"title":"多Sheet演示-新表","index":1}}},{"copySheet":{"source":{"sheetId":"39a262"},"destination":{"title":"多Sheet演示-副本"}}}]}'
```

返回成功：

- `addSheet.sheetId=1jDLTq`
- `addSheet.title=多Sheet演示-新表`
- `copySheet.sheetId=1jDLTr`
- `copySheet.title=多Sheet演示-副本`

```bash
./xfchat_cli sheets +info --as user --spreadsheet-token shtrzlFRtdfQZ2o8qauJpuPO8Rf
```

延迟约 2 秒后回读，确认当前保留的工作表有：

- `39a262 / Sheet1`
- `1jDLTr / 多Sheet演示-副本`

```bash
./xfchat_cli sheets +read --as user --spreadsheet-token shtrzlFRtdfQZ2o8qauJpuPO8Rf --sheet-id 1jDLTr --range A1:B6
```

返回成功，确认复制后的工作表 `1jDLTr` 保留了原数据：

- `模块 | 状态`
- `Sheets | 通过`
- `OpenAPI | 私有化已通`
- `用户态 | 正常`

```bash
./xfchat_cli api POST /open-apis/sheets/v2/spreadsheets/shtrzlFRtdfQZ2o8qauJpuPO8Rf/sheets_batch_update --as user --data '{"requests":[{"deleteSheet":{"sheetId":"1jDLTq"}}]}'
```

返回成功：

- `deleteSheet.result=true`
- `deleteSheet.sheetId=1jDLTq`

### 私有化现象

- `addSheet` 与 `copySheet` 返回成功后，`+info` 立即回读未必能同时看到两个新表。
- 延迟约 2 秒后再次回读，状态正常。
- 这与前面的写后读延迟一致，仍属于短暂最终一致性现象。

### 群内留痕

- `om_x100b53e9388e74a03864ed71d4f19fe`
