# 用户Drive模块测试

状态：通过

## 前置条件

- 身份： `user`
- 报告开始时状态： `lark-cli auth status` 显示当前没有已登录 user

## 测试记录

1. 上传本地文件 as user
命令： `lark-cli drive +upload --as user --file ./demo_reports/2026-03-29/artifacts/demo-upload.txt --name demo-upload-user.txt`
结果： `file_token=S2SCbJD6xo3eK2xLMGxcJhpon2d`

## 创建资源

- File token: `S2SCbJD6xo3eK2xLMGxcJhpon2d`

## 复现方式

1. 将源文件放在当前仓库工作目录下。
2. 使用相对路径上传。
3. Save the returned `file_token` for downstream demos.