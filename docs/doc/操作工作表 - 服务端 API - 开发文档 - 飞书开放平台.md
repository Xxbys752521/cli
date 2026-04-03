## 操作工作表

最后更新于 2024-11-28

本文内容

根据电子表格的 token 对工作表进行操作，包括增加工作表、复制工作表、删除工作表。

该接口和 [更新工作表属性](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/ugjMzUjL4IzM14COyMTN) 的请求地址相同，但参数不同，调用前请仔细阅读文档。

## 请求

| 基本 |  |
| --- | --- |
| HTTP URL | <https://open.xfchat.iflytek.com/open-apis/sheets/v2/spreadsheets/:spreadsheet\_token/sheets\_batch\_update> |
| HTTP Method | POST |
| 接口频率限制 | [100 次/秒](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/uUzN04SN3QjL1cDN) |
| 支持的应用类型 | 自建应用  商店应用 |
| 权限要求  开启任一权限即可 | 查看、评论、编辑和管理云空间中所有文件  查看、评论、编辑和管理电子表格 |

### 请求头

| 名称 | 类型 | 必填 | 描述 |
| --- | --- | --- | --- |
| Authorization | string | 是 | 应用调用 API 时，需要通过访问凭证（access\_token）进行身份鉴权，参考 [选择并获取访问凭证](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/uMTNz4yM1MjLzUzM#5aa2e490) 。  **值格式** ："Bearer `access_token` "  可选值如下：  - tenant\_access\_token  ：  以应用身份调用 API，可读写的数据范围由应用自身的 [数据权限范围](https://open.xfchat.iflytek.com/document/home/introduction-to-scope-and-authorization/configure-app-data-permissions) 决定。参考 [自建应用获取 tenant\_access\_token](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/tenant_access_token_internal) 或 [商店应用获取 tenant\_access\_token](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/ukDNz4SO0MjL5QzM/auth-v3/auth/tenant_access_token) 。 **示例值** ："Bearer t-g1044qeGEDXTB6NDJOGV4JQCYDGHRBARFTGT1234" - user\_access\_token  ：  以用户身份调用 API，可读写的数据范围由用户的的数据权限范围决定。参考 [获取 user\_access\_token](https://open.xfchat.iflytek.com/document/uAjLw4CM/ukTMukTMukTM/authentication-management/access-token/get-user-access-token) 。。 **示例值** ："Bearer u-cjz1eKCEx289x1TXEiQJqAh5171B4gDHPq00l0GE1234" |
| Content-Type | string | 是 | **固定值** ："application/json; charset=utf-8" |

### 路径参数

| 参数 | 类型 | 描述 |
| --- | --- | --- |
| spreadsheet\_token | string | 电子表格的 token。可通过以下两种方式获取。了解更多，参考 [电子表格概述](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/uATMzUjLwEzM14CMxMTN/overview) 。  - 电子表格的 URL：<https://sample.xfchat.iflytek.com/sheets/> ==Ios7sNNEphp3WbtnbCscPqabcef== - 调用 [获取文件夹中的文件清单](https://open.xfchat.iflytek.com/document/uAjLw4CM/ukTMukTMukTM/reference/drive-v1/file/list) |

### 请求体

### 请求体示例

```
{
  "requests": [
    {
      "addSheet": {
        "properties": {
          "title": "new_sheet",
          "index": 1
        }
      }
    },
    {
      "copySheet": {
        "source": {
          "sheetId": "2jm6f7"
        },
        "destination": {
          "title": "sheet_copy"
        }
      }
    },
    {
      "deleteSheet": {
        "sheetId": "l8Gdub"
      }
    }
  ]
}
```

### cURL 请求示例

```
curl --location --request POST 'https://open.xfchat.iflytek.com/open-apis/sheets/v2/spreadsheets/Ios7sNNEphp3WbtnbCscPqabcef/sheets_batch_update' \
--header 'Authorization: Bearer t-e346617a4acfc3a11d4ed24dca0d0c0fc8e0067e' \
--header 'Content-Type: application/json' \
--data-raw '{
  "requests": [
    {
      "addSheet": {
        "properties": {
          "title": "new_sheet",
          "index": 1
        }
      }
    },
    {
      "copySheet": {
        "source": {
          "sheetId": "2jm6f7"
        },
        "destination": {
          "title": "sheet_copy"
        }
      }
    },
    {
      "deleteSheet": {
        "sheetId": "l8Gdub"
      }
    }
  ]
}
'
```

## 响应

### 响应体

### 响应体示例

```
{
  "code": 0,
  "msg": "Success",
  "data": {
    "replies": [
      {
        "addSheet": {
          "properties": {
            "sheetId": "l8Gddg",
            "title": "new_sheet",
            "index": 1
          }
        },
        "copySheet": {
          "properties": {
            "sheetId": "dso4jn",
            "title": "sheet_copy",
            "index": 0
          }
        },
        "deleteSheet": {
          "result": true,
          "sheetId": "l8Gdub"
        }
      }
    ]
  }
}
```

### 错误码

具体可参考： [服务端错误码说明](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/ugjM14COyUjL4ITN)

本文内容
