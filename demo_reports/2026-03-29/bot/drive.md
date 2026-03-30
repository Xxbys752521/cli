# 机器人Drive模块测试

状态：通过

## 前置条件

- 身份： `bot`

## 测试记录

1. 安全校验
命令： `lark-cli drive +upload --as bot --file /absolute/path/to/file`
结果： 被拒绝，原因是 `--file` 必须是当前目录下的相对路径

2. 上传本地文件
命令： `lark-cli drive +upload --as bot --file ./demo_reports/2026-03-29/artifacts/demo-upload.txt --name demo-upload.txt`
结果： `file_token=MtTnbNkFWoWZChxn8A3cb5vqnFb`

3. 下载已上传文件
命令： `lark-cli drive +download --as bot --file-token MtTnbNkFWoWZChxn8A3cb5vqnFb --output ./demo_reports/2026-03-29/artifacts/demo-upload.downloaded.txt --overwrite`
结果： 已在本地保存 `43` 字节

4. 一致性校验
命令： `cmp -s ./demo_reports/2026-03-29/artifacts/demo-upload.txt ./demo_reports/2026-03-29/artifacts/demo-upload.downloaded.txt && echo SAME`
结果： `SAME`

## 创建资源

- File token: `MtTnbNkFWoWZChxn8A3cb5vqnFb`

## 复现方式

1. 将测试文件放在仓库目录内。
2. 使用相对路径上传。
3. Download by `file_token`.
4. 比对源文件和下载文件。