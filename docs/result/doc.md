# 云文档（Docs）测试结果

## 最终结论

- 云文档正文主链路已从 MCP 适配为 OpenAPI-first。
- 正文写入已改为“单次建块 + `blocks/batch_update` 批量更新文本元素”，不再逐块单请求写入。
- 当前已确认不依赖 MCP 的命令：
  - `docs +create`
  - `docs +fetch`
  - `docs +update`
- 主路径为：
  - `POST /open-apis/docx/v1/documents`
  - `GET /open-apis/docx/v1/documents/{document_id}`
  - `GET /open-apis/docx/v1/documents/{document_id}/blocks`
  - `GET /open-apis/docx/v1/documents/{document_id}/blocks/{document_id}/children`
  - `POST /open-apis/docx/v1/documents/{document_id}/blocks/{document_id}/children`
  - `PATCH /open-apis/docx/v1/documents/{document_id}/blocks/batch_update`
  - `DELETE /open-apis/docx/v1/documents/{document_id}/blocks/{document_id}/children/batch_delete`

## 代码适配

- 已更新 [shortcuts/doc/docs_create.go](/Users/wangqizhao/Developer/iflytek/cli/shortcuts/doc/docs_create.go)
- 已更新 [shortcuts/doc/docs_fetch.go](/Users/wangqizhao/Developer/iflytek/cli/shortcuts/doc/docs_fetch.go)
- 已更新 [shortcuts/doc/docs_update.go](/Users/wangqizhao/Developer/iflytek/cli/shortcuts/doc/docs_update.go)
- 已更新 [shortcuts/doc/helpers.go](/Users/wangqizhao/Developer/iflytek/cli/shortcuts/doc/helpers.go)

## 已打通能力

- `docs +create --as user`
- `docs +create --as bot`
- `docs +fetch --as user`
- `docs +fetch --as bot`
- `docs +update --mode append --as user`
- `docs +update --mode overwrite --as user`
- `docs +update --mode overwrite --as bot`
- 原生链接元素：`[text](url)`
- 原生提及元素：`@user[open_id]`
- 原生文档提及：`@doc[url或token]`

## 留痕对象

- User 文档：
  - `doxrzXQMuDvYEZocoNwy8vK9jTb`
  - 标题：`xfchat_cli Doc OpenAPI 测试`
- Bot 文档：
  - `doxrzy4IwmhwFekN58lWQTEjoag`
  - 标题：`xfchat_cli Doc OpenAPI Bot 测试`
- User 追加样例：
  - `doxrzKqENm2QjwcXoNKcZ5baFhe`
  - 标题：`xfchat_cli Doc OpenAPI Append 样例`
- User 覆盖样例：
  - `doxrzCv6C3ii2hnkTzYSi8n39Nf`
  - 标题：`xfchat_cli Doc OpenAPI 顺序验证`

## 关键命令与结果

```bash
./xfchat_cli docs +create --as user --title 'xfchat_cli Doc OpenAPI 测试' --markdown '# 标题\n\n第一段\n\n- 列表一\n- 列表二'
```

返回成功：

- `document_id=doxrzXQMuDvYEZocoNwy8vK9jTb`
- `backend=openapi`

```bash
./xfchat_cli docs +create --as bot --title 'xfchat_cli Doc OpenAPI Bot 测试' --markdown '# Bot标题\n\nBot 第一段'
```

返回成功：

- `document_id=doxrzy4IwmhwFekN58lWQTEjoag`
- `backend=openapi`

```bash
./xfchat_cli docs +create --as user --title 'xfchat_cli Doc OpenAPI Append 样例' --markdown $'原始一\n\n原始二'
./xfchat_cli docs +update --as user --doc doxrzKqENm2QjwcXoNKcZ5baFhe --mode append --markdown $'追加一\n\n追加二'
sleep 2
./xfchat_cli docs +fetch --as user --doc doxrzKqENm2QjwcXoNKcZ5baFhe
```

返回成功，最终内容为：

- `xfchat_cli Doc OpenAPI Append 样例`
- `原始一`
- `原始二`
- `追加一`
- `追加二`

```bash
./xfchat_cli docs +update --as user --doc doxrzCv6C3ii2hnkTzYSi8n39Nf --mode overwrite --markdown $'覆盖一\n\n覆盖二'
./xfchat_cli docs +fetch --as user --doc doxrzCv6C3ii2hnkTzYSi8n39Nf
```

返回成功，最终内容为：

- `xfchat_cli Doc OpenAPI 顺序验证`
- `覆盖一`
- `覆盖二`

```bash
./xfchat_cli docs +fetch --as bot --doc doxrzy4IwmhwFekN58lWQTEjoag
```

返回成功，确认 Bot 文档可读。

## 当前边界

- OpenAPI 适配当前只支持 `docs +update` 的：
  - `append`
  - `overwrite`
- 当前未适配：
  - `replace_range`
  - `replace_all`
  - `insert_before`
  - `insert_after`
  - `delete_range`
  - `--new-title`
  - 白板 markdown
- `docs +search` 这轮未改；它不属于 docx 正文编辑主链路。

## 富格式支持结论

### 当前已适配

- 一级标题：`block_type=3 / heading1`
- 二级标题：`block_type=4 / heading2`
- 三级标题：`block_type=5 / heading3`
- 无序列表：`block_type=12 / bullet`
- 有序列表：`block_type=13 / ordered`
- 引用：`block_type=15 / quote`
- 代码块：`block_type=14 / code`
- Todo：`block_type=17 / todo`
- Divider：`block_type=22 / divider`
- Callout：`block_type=19 / callout`
- Grid：`block_type=24 / grid`
- GridColumn：`block_type=25 / grid_column`
- Table：`block_type=31 / table`
- Table Cell：`block_type=32 / table_cell`
- Embedded Sheet：`block_type=30 / sheet`
- Image：`block_type=27 / image`
- File View：`block_type=33 / view`
- File：`block_type=23 / file`

### 实测样例

- 标题/列表/引用样例：
  - `doxrzKqENm2QjwcXoNKcZ5baFhe`
- 代码块样例：
  - `doxrzzY8UbCVzMkmVmJoo5kD3Wb`
- Todo/Divider 样例：
  - `doxrz2AvynPMxRfkSqgink9SXzc`
- 富块原子样例：
  - `doxrzxLI6i3Hnvdz5KSTOTbsvdd`

### 关键观测

```bash
./xfchat_cli api GET /open-apis/docx/v1/documents/doxrzKqENm2QjwcXoNKcZ5baFhe/blocks --as user
```

底层返回已确认存在：

- `block_type=3`，内容 `结构化标题`
- `block_type=12`，内容 `结构化列表一`
- `block_type=12`，内容 `结构化列表二`
- `block_type=15`，内容 `结构化引用`

```bash
./xfchat_cli api GET /open-apis/docx/v1/documents/doxrzzY8UbCVzMkmVmJoo5kD3Wb/blocks --as user
```

底层返回已确认：

- `block_type=14`
- `code.elements[0].text_run.content = fmt.Println("hello")`

```bash
./xfchat_cli docs +fetch --as user --doc doxrzKqENm2QjwcXoNKcZ5baFhe
```

当前已能回出 markdown 语义：

- `# 结构化标题`
- `- 结构化列表一`
- `- 结构化列表二`
- `> 结构化引用`

```bash
./xfchat_cli docs +fetch --as user --doc doxrzzY8UbCVzMkmVmJoo5kD3Wb
```

当前代码块可回出为：

- `````text``
- `fmt.Println("hello")`
- ``````

### 当前结论

- 文本内容：通
- 标题语义块：通
- 列表语义块：通
- 引用语义块：通
- 代码块语义块：通
- 更复杂的飞书专用块：未适配

当前 `docs` OpenAPI 版已经不再只是“纯文本写入”，而是至少打通了常用富文本块。

## 更多块类型补测

### Todo / Divider

```bash
./xfchat_cli docs +create --as user --title 'xfchat_cli Doc Todo Divider 结构化' --markdown $'- [ ] 未完成待办\n- [x] 已完成待办\n\n---\n\n普通正文'
```

底层返回已确认：

- `block_type=17 / todo`，`style.done=false`
- `block_type=17 / todo`，`style.done=true`
- `block_type=22 / divider`

`fetch` 当前可回出：

- `- [ ] 未完成待办`
- `- [x] 已完成待办`
- `---`

### Table

通过原始 OpenAPI 创建了 `2x2` 表格块，并向 4 个单元格分别写入：

- `表头A`
- `表头B`
- `数据1`
- `数据2`

底层返回已确认：

- `block_type=31 / table`
- `block_type=32 / table_cell`

### Callout

通过原始 OpenAPI 创建了高亮块：

- `block_type=19 / callout`
- `background_color=4`
- `border_color=4`
- `text_color=2`
- `emoji_id=gift`

服务端自动生成了子块，随后继续向其子树写入：

- `Callout 内容区已写入`

### Grid / GridColumn

通过原始 OpenAPI 创建了 `2` 列分栏：

- `block_type=24 / grid`
- `block_type=25 / grid_column`

服务端自动生成两列，随后分别写入：

- 左列：`左列内容`
- 右列：`右列内容`

### Embedded Sheet

通过原始 OpenAPI 创建了内嵌电子表格块：

- `sheet.token = shtrz8mCgyGDqNnrF9wbZrF2cpe_iK98dV`

随后继续调用 Sheets OpenAPI 写入：

- `列A | 列B`
- `值1 | 值2`
- `值3 | 值4`

延迟约 2 秒后回读正常。

### Image / File

通过 `docs +media-insert` 已插入：

- 图片块：
  - `block_type=27 / image`
  - `image.token = boxrzHsGZh65TPz8R2BPMKpx4Vd`
- 文件块：
  - `block_type=33 / view`
  - `block_type=23 / file`
  - `file.token = boxrzcSShc4J3BGjfEMg7kteLq4`
  - `file.name = README.md`

### 当前读取边界

- `docs +fetch` 现在已经能回出：
  - heading
  - bullet
  - ordered
  - quote
  - code
  - todo
  - divider
- 但对这些块仍不会完整保真为结构化 markdown：
  - callout
  - grid/grid_column
  - table
  - embedded sheet
  - image
  - file/view

例如在 `doxrzxLI6i3Hnvdz5KSTOTbsvdd` 上，`fetch` 当前会把：

- table 单元格文本
- callout 内部文本
- grid 两列文本

全部平铺到 markdown 中，但不会保留它们各自的块级结构；同样也不会把内嵌 sheet、图片、文件视图完整序列化回 markdown。

## 私有化现象

- 文档写入后立即读取，偶发会看到旧内容。
- 延迟约 2 秒后再次读取，内容正常。
- 这与 Sheets 模块表现一致，属于短暂最终一致性窗口。
- 在连续高频写入时，个别请求会返回空响应并被 CLI 识别为 `response parse error: EOF (body: )`，更像私有化写接口稳定性问题，不是某一种格式专属报错。

## 全块演示文档

- 文档：`doxrz4AeBEYjGOoi06sHM7Dc8Kd`
- 标题：`xfchat_cli 云文档全块炫技演示`
- 链接：
  - `https://yf2ljykclb.xfchat.iflytek.com/docx/doxrz4AeBEYjGOoi06sHM7Dc8Kd`
- 当前已集中展示的块：
  - `block_type=2 / text`
  - `block_type=3 / heading1`
  - `block_type=4 / heading2`
  - `block_type=5 / heading3`
  - `block_type=12 / bullet`
  - `block_type=13 / ordered`
  - `block_type=14 / code`
  - `block_type=15 / quote`
  - `block_type=17 / todo`
  - `block_type=19 / callout`
  - `block_type=22 / divider`
  - `block_type=24 / grid`
  - `block_type=25 / grid_column`
  - `block_type=27 / image`
  - `block_type=30 / sheet`
  - `block_type=31 / table`
  - `block_type=32 / table_cell`
  - `block_type=33 / view`
  - `block_type=23 / file`

## User / Bot 对等性

- 不完全对等。
- 已确认基本对等的主链路：
  - `docs +create`
  - `docs +fetch`
  - `docs +update`
  - `docs +media-insert`
  - `docs +media-download`
- 明确不对等的能力：
  - `docs +search`：仅支持 `user`，`bot` 直接被命令限制拒绝
  - `drive.file.comments.create_v2`：
    - `user`、`bot` 两侧对真实 `docx` 都直接返回 `HTTP 404`
- 资源归属相关差异：
  - `drive.permission.members.auth` 在 `user` 文档上对 `user` 返回 `view/edit/manage_public = true`
  - 同一用户文档对 `bot` 返回 `edit = false`
  - `bot` 自己创建的文档上，`bot` 的 `view = true`

## 外围能力

- 搜索：
  - `docs +search --as user` 已通，但当前对标题 `xfchat_cli 云文档全块炫技演示` 返回空结果
  - `docs +search --as bot` 不支持
- 权限鉴权：
  - `drive.permission.members.auth` 已通
  - `user` 文档 `doxrz4AeBEYjGOoi06sHM7Dc8Kd`：
    - `user/view = true`
    - `user/manage_public = true`
    - `bot/edit = false`
  - `bot` 文档 `doxrzy4IwmhwFekN58lWQTEjoag`：
    - `bot/view = true`
- 协作者：
  - `drive.permission.members.create` 已通
  - `user`、`bot` 都已成功把测试群 `oc_010df6b42b975ce056cc7c2e717abde8` 加为 `view`
  - 当前 CLI 暴露的 `permission.members` 能力只有：
    - `auth`
    - `create`
    - `transfer_owner`
  - 未暴露：
    - `retrieve`
    - `delete`
    - `update`
  - `transfer_owner` 结果：
    - `user` 自己的临时文档 `doxrz6U0hXq6VuCrsuLRu0kUuCe` 对自己执行转移，返回 `1063002 Permission denied`
    - `bot` 临时文档 `doxrzvkz7i6esfVWFGfh40viSDd` 转移到用户 `ou_6b90260ffbda660ec3e47f27c0871ec9` 返回 `Success`
    - 转移后 `user`、`bot` 两侧仍可 `fetch` 该文档
    - 转移后权限探针确认：
      - `user/edit = true`
      - `bot/view = true`
- 评论：
  - `drive.file.comments.list` 已通
  - `user` 文档和 `bot` 文档当前都可正常列出评论，现阶段均为空
  - `drive.file.comments.create_v2` 尚未打通：
    - `user` token 已刷新到 `docs:document.comment:create`
    - 但 `user`、`bot` 两侧调用都返回 `HTTP 404`
    - 直接探测原始路径 `POST /open-apis/drive/v1/files/{file_token}/new_comments` 结果一致
    - 可定性为私有化服务未提供评论创建接口，不是 CLI 或权限问题

## 群内留痕

- `om_x100b53e9c1d53ca0385658a9972dfb6`
- `om_x100b53e9a1747ca438600d32e111c9f`
- `om_x100b53e9b289b0a8386947c689ef847`
- `om_x100b53ea5dd124a03862caea27bf7bd`
- `om_x100b53ea59df88a0386a2a2d1822eba`
- `om_x100b53ea5102a4a03859d8e1818e599`
