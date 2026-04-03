## 批量更新块的内容

最后更新于 2025-03-05

本文内容

批量更新块的富文本内容。

## 前提条件

调用此接口前，请确保当前调用身份（tenant\_access\_token 或 user\_access\_token）已有云文档的阅读、编辑等文档权限，否则接口将返回 HTTP 403 或 400 状态码。了解更多，参考 [如何为应用或用户开通文档权限](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/uczNzUjL3czM14yN3MTN#16c6475a) 。

## 请求

| 基本 |  |
| --- | --- |
| HTTP URL | <https://open.xfchat.iflytek.com/open-apis/docx/v1/documents/:document\_id/blocks/batch\_update> |
| HTTP Method | PATCH |
| 支持的应用类型 | 自建应用  商店应用 |
| 权限要求 | 创建及编辑新版文档 |
| 字段权限要求 | 获取用户 user ID  仅自建应用 |

### 请求头

| 名称 | 类型 | 必填 | 描述 |
| --- | --- | --- | --- |
| Authorization | string | 是 | tenant\_access\_token  或  user\_access\_token  **值格式** ："Bearer `access_token` "  **示例值** ："Bearer u-7f1bcd13fc57d46bac21793a18e560"  [了解更多：如何选择与获取 access token](https://open.xfchat.iflytek.com/document/uAjLw4CM/ugTN1YjL4UTN24CO1UjN/trouble-shooting/how-to-choose-which-type-of-token-to-use) |
| Content-Type | string | 是 | **固定值** ："application/json; charset=utf-8" |

### 路径参数

### 查询参数

### 请求体

### 请求示例

```
curl --location --request PATCH 'https://open.xfchat.iflytek.com/open-apis/docx/v1/documents/doxcnAJ9VRRJqVMYZ1MyKnavXWe/blocks/batch_update' \
--header 'Authorization: Bearer u-xxx'
--header 'Content-Type: application/json' \
--data-raw '{
    "requests": [
        {
            "block_id": "doxcnk0i44OMOaouw8AdXuXrp6b",
            "merge_table_cells": {
                "column_end_index": 2,
                "column_start_index": 0,
                "row_end_index": 1,
                "row_start_index": 0
            }
        },
        {
            "block_id": "doxcn0K8iGSMW4Mqgs9qlyTP50d",
            "update_text_elements": {
                "elements": [
                    {
                        "text_run": {
                            "content": "Hello",
                            "text_element_style": {
                                "background_color": 2,
                                "bold": true,
                                "italic": true,
                                "strikethrough": true,
                                "text_color": 2,
                                "underline": true                            }
                        }
                    },
                    {
                        "text_run": {
                            "content": "World ",
                            "text_element_style": {
                                "italic": true
                            }
                        }
                    },
                    {
                        "mention_doc": {
                            "obj_type": 22,
                            "token": "doxcnAJ9VRRJqVMYZ1MyKnayXWe",
                            "url": "https://test.xfchat.iflytek.com/docx/doxcnAJ9VRRJqVMYZ1MyKnayXWe"
                        }
                    }
                ]
            }
        }
    ]
}'
# 调用前请替换 'Authorization: Bearer u-xxx' 中的 'u-xxx' 为真实的访问令牌
```

### 请求体示例

## 响应

### 响应体

### 响应体示例

### 错误码

| HTTP状态码 | 错误码 | 描述 | 排查建议 |
| --- | --- | --- | --- |
| 400 | 1770001 | invalid param | 确认传入的参数是否合法 |
| 404 | 1770002 | not found | **文档场景中** ：   文档的 `document_id` 不存在。请确认文档是否已被删除或 `document_id` 是否填写正确。参考 [文档概述](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/uUDN04SN0QjL1QDN/document-docx/docx-overview) 了解如何获取文档的 `document_id` 。   **群公告场景中** ：   群 ID `chat_id` 不存在。请确认群是否被解散或 `chat_id` 是否填写正确。 |
| 400 | 1770003 | resource deleted | 确认资源是否已被删除 |
| 400 | 1770004 | too many blocks in document | 确认文档 Block 数量是否超上限 |
| 400 | 1770005 | too deep level in document | 确认文档 Block 层级是否超上限 |
| 400 | 1770006 | schema mismatch | 确认文档结构是否合法 |
| 400 | 1770007 | too many children in block | 确认指定 Block 的 Children 数量是否超上限 |
| 400 | 1770008 | too big file size | 确认上传的文件尺寸是否超上限 |
| 400 | 1770010 | too many table column | 确认表格列数是否超上限 |
| 400 | 1770011 | too many table cell | 确认表格单元格数量是否超上限 |
| 400 | 1770012 | too many grid column | 确认 Grid 列数量是否超上限 |
| 400 | 1770013 | relation mismatch | 图片、文件等资源的关联关系不正确。请确保在创建图片、文件块时，同时上传了相关图片或文件素材至对应的文档块中。详情参考文档 [常见问题 3 和 4](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/uUDN04SN0QjL1QDN/document-docx/docx-v1/faq#1908ddf0) |
| 400 | 1770014 | parent children relation mismatch | 确认 Block 父子关系是否正确 |
| 400 | 1770015 | single edit with multi document | 确认 Block 所属文档与指定的 Document 是否相同 |
| 400 | 1770019 | repeated blockID in document | 确认 Document 中的 BlockID 是否有重复 |
| 400 | 1770020 | operation denied on copying document | 确认 Document 是否正在创建副本中 |
| 400 | 1770021 | too old document | 确认指定的 Document 版本（Revision\_id）是否过旧。指定的版本号与文档最新版本号差值不能超过 1000 |
| 400 | 1770022 | invalid page token | 确认查询参数中的 page\_token 是否合法 |
| 400 | 1770024 | invalid operation | 确认操作是否合法 |
| 400 | 1770025 | operation and block not match | 确认指定 Block 应用对应操作是否合法 |
| 400 | 1770026 | row operation over range | 确认行操作下标是否越界 |
| 400 | 1770027 | column operation over range | 确认列操作下标是否越界 |
| 400 | 1770028 | block not support create children | 确认指定 Block 添加 Children 是否合法 |
| 400 | 1770029 | block not support to create | 确认指定 Block 是否支持创建 |
| 400 | 1770030 | invalid parent children relation | 确认指定操作其父子关系是否合法 |
| 400 | 1770031 | block not support to delete children | 确认指定 Block 是否支持删除 Children |
| 400 | 1770033 | raw content size exceed limited | 纯文本内容大小超过限制 |
| 400 | 1770034 | operation count exceed limited | 当前请求中涉及单元格个数过多，请拆分成多次请求 |
| 400 | 1770035 | resource count exceed limit | 当前请求中资源的数目超限，请拆分成多次请求。各类资源上限为：ChatCard 200 张，File 200 个，MentionDoc 200 个，MentionUser 200 个，Image 20 张，ISV 20 个，Sheet 5 篇，Bitable 5 篇。 |
| 400 | 1770038 | resource not found | 未查询到插入的资源或资源无权限插入，请检查资源标识是否正确。 |
| 403 | 1770032 | forbidden | **文档场景中** ：  确认当前调用身份是否有文档阅读（获取相关接口）或编辑（更新、删除、创建相关接口）权限。请参考以下方式解决：  - 如果你使用的是 `tenant_access_token` ，意味着当前应用没有文档权限。你需通过云文档网页页面右上方 **「...」** -> **「...更多」** -> **「添加文档应用」** 入口为应用添加文档权限。  **说明** ：在 **添加文档应用** 前，你需确保目标应用至少开通了一个云文档或多维表格的 [API 权限](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/uYTM5UjL2ETO14iNxkTN/scope-list) 。否则你将无法在文档应用窗口搜索到目标应用。  ![](https://sf1-lark-tos.xfchat.iflytek.com/obj/open-platform-opendoc/open-platform-opendoc_22c027f63c540592d3ca8f41d48bb107_CSas7OYJBR.png?height=1994&maxWidth=550&width=3278) - 如果你使用的是 `user_access_token` ，意味着当前用户没有文档权限。你需通过云文档网页页面右上方 **分享** 入口为当前用户添加文档权限。  ![image.png](https://sf1-lark-tos.xfchat.iflytek.com/obj/open-platform-opendoc/open-platform-opendoc_3e052d3bac56f9441296ae22e2969d63_a2DEYrJup8.png?height=278&maxWidth=550&width=1383)  了解具体操作步骤或其它添加权限方式，参考 [云文档常见问题 3](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/uczNzUjL3czM14yN3MTN#16c6475a) 。  对于创建和更新相关接口，你还需要确认：  - 当前调用身份是否有 MentionDoc 即 @文档 中文档的阅读权限 - MentionUser 即 @用户 中的用户是否在职且与当前调用身份互为联系人 - 当前调用身份是否具有群卡片的查看和分享权限 - 当前调用身份是否具有访问指定 Wiki 即知识库子目录的权限 - 当前调用身份是否具有 OKR、ISV、Add-Ons 等文档块的查看权限  **群公告场景中** ：  当前的操作者没有群公告的编辑权限。解决方案：  - 方案一：调用 [指定群管理员](https://open.xfchat.iflytek.com/document/uAjLw4CM/ukTMukTMukTM/reference/im-v1/chat-managers/add_managers) 接口，将当前操作者置为群管理员，然后重试。 - 方案二：在 **飞书客户端 > 群组 > 设置 > 群管理** 中，将 **谁可以编辑群信息** 设置为 **所有群成员** ，然后重试。  对于创建和更新相关接口，你还需要确认：  - 当前调用身份是否有 MentionDoc 即 @文档 中文档的阅读权限 - MentionUser 即 @用户 中的用户是否在职且与当前调用身份互为联系人 - 当前调用身份是否具有群卡片的查看和分享权限 - 当前调用身份是否具有访问指定 Wiki 即知识库子目录的权限 - 当前调用身份是否具有 OKR、ISV、Add-Ons 等文档块的查看权限 |
| 500 | 1771001 | server internal error | 服务器内部错误。请重试，若仍无法解决请咨询 [技术支持](https://applink.xfchat.iflytek.com/TLJpeNdW) 。 |
| 500 | 1771006 | mount folder failed | 挂载文档到云空间文件夹失败 |
| 500 | 1771002 | gateway server internal error | 网关服务内部错误。请重试，若仍无法解决请咨询 [技术支持](https://applink.xfchat.iflytek.com/TLJpeNdW) 。 |
| 500 | 1771003 | gateway marshal error | 网关服务解析错误。请重试，若仍无法解决请咨询 [技术支持](https://applink.xfchat.iflytek.com/TLJpeNdW) 。 |
| 500 | 1771004 | gateway unmarshal error | 网关服务反解析错误。请重试，若仍无法解决请咨询 [技术支持](https://applink.xfchat.iflytek.com/TLJpeNdW) 。 |
| 503 | 1771005 | system under maintenance | 系统服务正在维护中 |

本文内容
