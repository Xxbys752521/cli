# lark-cli 状态机 Demo 运行报告

- 运行时间：2026-03-30 12:16:49 +0800
- 执行用户：王琦钊
- 目标：从建群、构造聊天记录、建答疑表和知识库，到按状态机跑完自动答疑和转人工两条链路

## 0. 准备上下文

- 当前用户 open_id：`ou_4b6bc8c3ea123aa0d47b290828fc7891`
- Demo 群名：`状态机Demo群-20260330_121648`
- Demo Base：`状态机Demo答疑表-20260330_121648`
- Demo 知识库文档：`状态机Demo知识库-20260330_121648`

## 1. 创建知识库文档

### doc_create

```bash
lark-cli docs +create --as user --title '状态机Demo知识库-20260330_121648' --markdown "# 常见答疑知识库

## 案例 1：私有化环境知识库打不开

- 现象：登录飞书后点击知识库页面打不开，或提示无权限
- 适用范围：私有化部署环境、新员工或新租户初次接入
- 处理建议：
  1. 确认当前账号已完成飞书登录
  2. 确认知识库对当前用户或部门已授权
  3. 如果是私有化环境，检查是否使用了正确的环境域名和对应租户
  4. 如果页面空白，优先检查是否存在权限未同步或网关配置异常
- 输出口径：先引导补充环境信息，再给出权限检查路径"
```

```json
{
  "ok": true,
  "identity": "user",
  "data": {
    "doc_id": "LE1Od5OW8opwZyxU9ztcAO6rn1f",
    "doc_url": "https://www.feishu.cn/docx/LE1Od5OW8opwZyxU9ztcAO6rn1f",
    "log_id": "202603301216496C3946D6C3A735EA05BD",
    "message": "文档创建成功"
  }
}
```


## 2. 创建答疑 Base 和表

### base_create

```bash
lark-cli base +base-create --as user --name '状态机Demo答疑表-20260330_121648' --time-zone Asia/Shanghai
```

```json
{
  "ok": true,
  "identity": "user",
  "data": {
    "base": {
      "base_token": "WI0TbNEbxain2EsAHlKcqHcrn25",
      "folder_token": "",
      "name": "状态机Demo答疑表-20260330_121648",
      "time_zone": "Asia/Shanghai",
      "url": "https://xcn596olvq0a.feishu.cn/base/WI0TbNEbxain2EsAHlKcqHcrn25"
    },
    "created": true
  }
}
```

### table_create

```bash
lark-cli base +table-create --as user --base-token WI0TbNEbxain2EsAHlKcqHcrn25 --name '答疑表' --fields '[
  {"name":"问题ID","type":"text"},
  {"name":"当前状态","type":"select","multiple":false,"options":[{"name":"待补充"},{"name":"待检索"},{"name":"已自动回答"},{"name":"已转人工"},{"name":"人工已处理"},{"name":"已回写知识库"},{"name":"已结束"}]},
  {"name":"提问人","type":"text"},
  {"name":"问题内容","type":"text"}
]'
```

```json

```

