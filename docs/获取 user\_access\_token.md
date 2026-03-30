##

最后更新于 2025-04-08

本文内容

OAuth 令牌接口，可用于获取 `user_access_token` 以及 `refresh_token` 。 `user_access_token` 为用户访问凭证，使用该凭证可以以用户身份调用 OpenAPI。 `refresh_token` 为刷新凭证，可以用来获取新的 `user_access_token` 。

- 获取 `user_access_token` 前需要先获取授权码，详见 [获取授权码](https://open.xfchat.iflytek.com/document/common-capabilities/sso/api/obtain-oauth-code) 。请注意授权码的有效期为 5 分钟，且只能被使用一次。
- 获取到的 `user_access_token` 存在有效期，如何刷新 `user_access_token` 详见 [刷新 user\_access\_token](https://open.xfchat.iflytek.com/document/uAjLw4CM/ukTMukTMukTM/authentication-management/access-token/refresh-user-access-token) 。
- 如果你需要获取用户信息，详见 [获取用户信息](https://open.xfchat.iflytek.com/document/uAjLw4CM/ukTMukTMukTM/reference/authen-v1/user_info/get) 。

## 请求

| 基本 |  |
| --- | --- |
| HTTP URL | <https://open.xfchat.iflytek.com/open-apis/authen/v2/oauth/token> |
| HTTP Method | POST |
| 接口频率限制 | [1000 次/分钟、50 次/秒](https://open.xfchat.iflytek.com/document/ukTMukTMukTM/uUzN04SN3QjL1cDN) |
| 支持的应用类型 | 自建应用  商店应用 |
| 权限要求 | 无 |
| 字段权限要求 | `refresh_token` 以及 `refresh_token_expires_in` 字段仅在具备以下权限时返回：  离线访问已授权数据 |

### 请求头

| 名称 | 类型 | 必填 | 描述 |
| --- | --- | --- | --- |
| Content-Type | string | 是 | 请求体类型。  **固定值：** `application/json; charset=utf-8` |

### 请求体

| 名称 | 类型 | 必填 | 描述 |
| --- | --- | --- | --- |
| grant\_type | string | 是 | 授权类型。  **固定值：** `authorization_code` |
| client\_id | string | 是 | 应用的 App ID。应用凭证 App ID 和 App Secret 获取方式：  1. 登录 [飞书开发者后台](https://open.xfchat.iflytek.com/app) 。 2. 进入应用详情页，在左侧导航栏，单击 **凭证与基础信息** 。 3. 在 **应用凭证** 区域，获取并保存 **App ID** 和 **App Secret** 。  **示例值：** `cli_a5ca35a685b0x26e` |
| client\_secret | string | 是 | 应用的 App Secret。应用凭证 App ID 和 App Secret 获取方式：  1. 登录 [飞书开发者后台](https://open.xfchat.iflytek.com/app) 。 2. 进入应用详情页，在左侧导航栏，单击 **凭证与基础信息** 。 3. 在 **应用凭证** 区域，获取并保存 **App ID** 和 **App Secret** 。  **示例值：** `baBqE5um9LbFGDy3X7LcfxQX1sqpXlwy` |
| code | string | 是 | 授权码，详见 [获取授权码](https://open.xfchat.iflytek.com/document/common-capabilities/sso/api/obtain-oauth-code) 。  **示例值：** `a61hb967bd094dge949h79bbexd16dfe` |
| redirect\_uri | string | 否 | 在构造授权页页面链接时所拼接的应用回调地址。  **示例值：** `https://example.com/api/oauth/callback` |
| code\_verifier | string | 否 | 在发起授权前，本地生成的随机字符串，用于 PKCE（Proof Key for Code Exchange）流程。使用 PKCE 时，该值为必填项。  有关 PKCE 的详细介绍，请参阅 [RFC 7636 - Proof Key for Code Exchange by OAuth Public Clients](https://datatracker.ietf.org/doc/html/rfc7636) 。  **长度限制：** 最短 43 字符，最长 128 字符  **可用字符集：** \[A-Z\] / \[a-z\] / \[0-9\] / "-" / "." / "\_" / "~"  **示例值：** `TxYmzM4PHLBlqm5NtnCmwxMH8mFlRWl_ipie3O0aVzo` |
| scope | string | 否 | 用于缩减 `user_access_token` 的权限范围。如果未指定，生成的 `user_access_token` 将包含用户授权时的所有权限。请务必注意该列表中不得包含重复项（错误码 20067）以及未授权项（错误码 20068）。  **说明** ：  - 多次调用当前接口缩减不同的权限不会叠加。例如，第一调用该接口缩减了权限 A，则 `user_access_token` 不包含权限 A；第二次调用该接口缩减了权限 B，则 `user_access_token` 不包含权限 B。 - 生效的权限列表可通过本接口返回值 scope 查看。  **格式要求：** 以空格分隔的 `scope` 列表  **示例值：** `auth:user.id:read task:task:read` |

### 请求体示例

```
{
    "grant_type": "authorization_code",
    "client_id": "cli_a5ca35a685b0x26e",
    "client_secret": "baBqE5um9LbFGDy3X7LcfxQX1sqpXlwy",
    "code": "a61hb967bd094dge949h79bbexd16dfe",
    "redirect_uri": "https://example.com/api/oauth/callback",
    "code_verifier": "TxYmzM4PHLBlqm5NtnCmwxMH8mFlRWl_ipie3O0aVzo"
}
```

## 响应

响应体类型为 `application/json; charset=utf-8` 。

### 响应体

| 名称 | 类型 | 描述 |
| --- | --- | --- |
| code | string | 错误码，为 0 时表明请求成功，非 0 表示失败，请参照下文 [错误码](https://open.xfchat.iflytek.com/document/uAjLw4CM/ukTMukTMukTM/authentication-management/access-token/get-user-access-token#%E9%94%99%E8%AF%AF%E7%A0%81) 一节进行相应处理 |
| access\_token | string | 即 `user_access_token` ，仅在请求成功时返回 |
| expires\_in | int | 即 `user_access_token` 的有效期，单位为秒，仅在请求成功时返回 |
| refresh\_token | string | 用于刷新 `user_access_token` ，详见 [刷新 user\_access\_token](https://open.xfchat.iflytek.com/document/uAjLw4CM/ukTMukTMukTM/authentication-management/access-token/refresh-user-access-token) 。该字段仅在请求成功且用户授予 `offline_access` 权限时返回。 |
| refresh\_token\_expires\_in | int | 即 `refresh_token` 的有效期，单位为秒，仅在返回 `refresh_token` 时返回。 |
| token\_type | string | 值固定为 `Bearer` ，仅在请求成功时返回 |
| scope | string | 本次请求所获得的 `access_token` 所具备的权限列表，以空格分隔，仅在请求成功时返回 |
| error | string | 错误类型，仅在请求失败时返回 |
| error\_description | string | 具体的错误信息，仅在请求失败时返回 |

### 响应体示例

成功响应示例：

```
{
    "code": 0,
    "access_token": "eyJhbGciOiJFUzI1NiIs**********X6wrZHYKDxJkWwhdkrYg",
    "expires_in": 7200, // 非固定值，请务必根据响应体中返回的实际值来确定 access_token 的有效期
    "refresh_token": "eyJhbGciOiJFUzI1NiIs**********XXOYOZz1mfgIYHwM8ZJA",
    "refresh_token_expires_in": 604800, // 非固定值，请务必根据响应体中返回的实际值来确定 refresh_token 的有效期
    "scope": "auth:user.id:read offline_access task:task:read user_profile",
    "token_type": "Bearer"
}
```

失败响应示例：

```
{
    "code": 20050,
    "error": "server_error",
    "error_description": "An unexpected server error occurred. Please retry your request."
}
```

### 错误码

| HTTP 状态码 | 错误码 | 描述 | 排查建议 |
| --- | --- | --- | --- |
| 400 | 20001 | The request is missing a required parameter. | 必要参数缺失，请检查请求时传入的参数是否有误 |
| 400 | 20002 | The client secret is invalid. | 应用认证失败，请检查提供的 `client_id` 与 `client_secret` 是否正确 |
| 400 | 20003 | The authorization code is not found. Please note that an authorization code can only be used once. | 无效的授权码，请检查授权码是否有效，注意授权码仅能使用一次 |
| 400 | 20004 | The authorization code has expired. | 授权码已经过期，请在授权码生成后的 5 分钟内使用 |
| 400 | 20008 | The user does not exist. | 用户不存在，请检查发起授权的用户的当前状态 |
| 400 | 20009 | The specified app is not installed. | 租户未安装应用，请检查应用状态 |
| 400 | 20010 | The user does not have permission to use this app. | 用户无应用使用权限，请检查发起授权的用户是否仍具有应用使用权限 |
| 400 | 20024 | The provided authorization code or refresh token does not match the provided client ID. | 提供的授权码与 `client_id` 不匹配，请勿混用不同应用的凭证 |
| 400 | 20036 | The specified grant\_type is not supported. | 无效的 `grant_type` ，请检查请求体中 `grant_type` 字段的取值 |
| 400 | 20048 | The specified app does not exist. | 应用不存在，请检查应用状态 |
| 400 | 20049 | PKCE code challenge failed. | PKCE 校验失败，请检查请求体中 `code_verifier` 字段是否存在且有效 |
| 500 | 20050 | An unexpected server error occurred. Please retry your request. | 内部服务错误，请稍后重试，如果持续报错请联系 [技术支持](https://applink.xfchat.iflytek.com/TLJpeNdW) |
| 400 | 20063 | The request is malformed. Please check your request. | 请求体中缺少必要字段，请根据具体的错误信息补齐字段 |
| 400 | 20065 | The authorization code has been used. Please note that an authorization code can only be used once. | 授权码已被使用，授权码仅能使用一次，请检查是否有被重复使用 |
| 400 | 20066 | The user status is invalid. | 用户状态非法，请检查发起授权的用户的当前状态 |
| 400 | 20067 | The provided scope list contains duplicate scopes. Please ensure all scopes are unique. | 无效的 `scope` 列表，其中存在重复项，请确保传入的 `scope` 列表中没有重复项 |
| 400 | 20068 | The provided scope list contains scopes that are not permitted. Please ensure all scopes are allowed. | 无效的 `scope` 列表，其中存在应用未申请的权限项，请确保应用具备 `scope` 列表内列出的所有权限 |
| 400 | 20069 | The specified app is not enabled. | 应用未启用，请检查应用状态 |
| 400 | 20070 | Multiple authentication methods were provided. Please only use one to proceed. | 请求使用了 `Basic Authentication` ，同时也在请求体中传递了 `client_secret` ，请仅使用一种认证方式 |
| 400 | 20071 | The provided redirect URI does not match the one used during authorization. | 无效的 `redirect_uri` ，请确保 `redirect_uri` 与 [获取授权码](https://open.xfchat.iflytek.com/document/common-capabilities/sso/api/obtain-oauth-code) 时传入的 `redirect_uri` 保持一致 |
| 503 | 20072 | The server is temporarily unavailable. Please retry your request. | 服务暂不可用，请稍后重试 |

## 代码示例

### Golang

运行下面示例程序的步骤：

1. 点击下方代码块右上角复制按钮，将代码复制到本地文件中，保存为 `main.go` ；
2. 参照注释部分，完成配置；
3. 在 `main.go` 所在目录下新建 `.env` 文件，内容如下：

 ```
 APP_ID=cli_xxxxxx # 仅为示例值，请使用你的应用的 App ID，获取方式：开发者后台 -> 基础信息 -> 凭证与基础信息 -> 应用凭证 -> App ID
 APP_SECRET=xxxxxx # 仅为示例值，请使用你的应用的 App Secret，获取方式：开发者后台 -> 基础信息 -> 凭证与基础信息 -> 应用凭证 -> App Secret
 ```

1. 在 `main.go` 所在目录执行以下命令：

 ```
 go mod init oauth-test
 go get github.com/gin-gonic/gin
 go get github.com/gin-contrib/sessions
 go get github.com/gin-contrib/sessions/cookie
 go get github.com/joho/godotenv
 go get golang.org/x/oauth2
 go run main.go
 ```

1. 浏览器打开 [http://localhost:8080](http://localhost:8080/) ，按照页面提示完成授权流程；

```
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "os"

    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
    _ "github.com/joho/godotenv/autoload"
    "golang.org/x/oauth2"
)

var oauthEndpoint = oauth2.Endpoint{
    AuthURL:  "https://accounts.xfchat.iflytek.com/open-apis/authen/v1/authorize",
    TokenURL: "https://open.xfchat.iflytek.com/open-apis/authen/v2/oauth/token",
}

var oauthConfig = &oauth2.Config{
    ClientID:     os.Getenv("APP_ID"),
    ClientSecret: os.Getenv("APP_SECRET"),
    RedirectURL:  "http://localhost:8080/callback", // 请先添加该重定向 URL，配置路径：开发者后台 -> 开发配置 -> 安全设置 -> 重定向 URL -> 添加
    Endpoint:     oauthEndpoint,
    Scopes:       []string{"offline_access"}, // 如果你不需要 refresh_token，请注释掉该行，否则你需要先申请 offline_access 权限方可使用，配置路径：开发者后台 -> 开发配置 -> 权限管理      
}

func main() {
    r := gin.Default()

    // 使用 Cookie 存储 session
    store := cookie.NewStore([]byte("secret")) // 此处仅为示例，务必不要硬编码密钥
    r.Use(sessions.Sessions("mysession", store))

    r.GET("/", indexController)
    r.GET("/login", loginController)
    r.GET("/callback", oauthCallbackController)

    fmt.Println("Started running on http://localhost:8080")
    log.Fatal(r.Run(":8080"))
}

func indexController(c *gin.Context) {
    c.Header("Content-Type", "text/html; charset=utf-8")
    var username string
    session := sessions.Default(c)
    if session.Get("user") != nil {
       username = session.Get("user").(string)
    }
    html := fmt.Sprintf(\`<html><head><style>body{font-family:Arial,sans-serif;background:#f4f4f4;margin:0;display:flex;justify-content:center;align-items:center;height:100vh}.container{text-align:center;background:#fff;padding:30px;border-radius:10px;box-shadow:0 0 10px rgba(0,0,0,0.1)}a{padding:10px 20px;font-size:16px;color:#fff;background:#007bff;border-radius:5px;text-decoration:none;transition:0.3s}a:hover{background:#0056b3}}</style></head><body><div class="container"><h2>欢迎%s！</h2><a href="/login">使用飞书登录</a></div></body></html>\`, username)
    c.String(http.StatusOK, html)
}

func loginController(c *gin.Context) {
    session := sessions.Default(c)

    // 生成随机 state 字符串，你也可以用其他有意义的信息来构建 state
    state := fmt.Sprintf("%d", rand.Int())
    // 将 state 存入 session 中
    session.Set("state", state)
    // 生成 PKCE 需要的 code verifier
    verifier := oauth2.GenerateVerifier()
    // 将 code verifier 存入 session 中
    session.Set("code_verifier", verifier)
    session.Save()

    url := oauthConfig.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))
    // 用户点击登录时，重定向到授权页面
    c.Redirect(http.StatusTemporaryRedirect, url)
}

func oauthCallbackController(c *gin.Context) {
    session := sessions.Default(c)
    ctx := context.Background()

    // 从 session 中获取 state
    expectedState := session.Get("state")
    state := c.Query("state")

    // 如果 state 不匹配，说明是 CSRF 攻击，拒绝处理
    if state != expectedState {
       log.Printf("invalid oauth state, expected '%s', got '%s'\n", expectedState, state)
       c.Redirect(http.StatusTemporaryRedirect, "/")
       return
    }

    code := c.Query("code")
    // 如果 code 为空，说明用户拒绝了授权
    if code == "" {
       log.Printf("error: %s", c.Query("error"))
       c.Redirect(http.StatusTemporaryRedirect, "/")
       return
    }
    codeVerifier, _ := session.Get("code_verifier").(string)
    // 使用获取到的 code 获取 token
    token, err := oauthConfig.Exchange(ctx, code, oauth2.VerifierOption(codeVerifier))
    if err != nil {
       log.Printf("oauthConfig.Exchange() failed with '%s'\n", err)
       c.Redirect(http.StatusTemporaryRedirect, "/")
       return
    }

    client := oauthConfig.Client(ctx, token)

    req, err := http.NewRequest("GET", "https://open.xfchat.iflytek.com/open-apis/authen/v1/user_info", nil)
    if err != nil {
       log.Fatal(err)
    }
    req.Header.Set("Authorization", "Bearer "+token.AccessToken)

    // 使用 token 发起请求，获取用户信息
    resp, err := client.Do(req)
    if err != nil {
       log.Printf("client.Get() failed with '%s'\n", err)
       c.Redirect(http.StatusTemporaryRedirect, "/")
       return
    }
    defer resp.Body.Close()

    var user struct {
       Data struct {
          Name string \`json:"name"\`
       } \`json:"data"\`
    }
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
       log.Printf("json.NewDecoder() failed with '%s'\n", err)
       c.Redirect(http.StatusTemporaryRedirect, "/")
       return
    }
    // 后续可以用获取到的用户信息构建登录态，此处仅为示例，请勿直接使用
    session.Set("user", user.Data.Name)
    session.Save()

    c.Header("Content-Type", "text/html; charset=utf-8")
    html := fmt.Sprintf(\`<html><head><style>body{font-family:Arial,sans-serif;background:#f4f4f4;margin:0;display:flex;justify-content:center;align-items:center;height:100vh}.container{text-align:center;background:#fff;padding:30px;border-radius:10px;box-shadow:0 0 10px rgba(0,0,0,0.1)}a{padding:10px 20px;font-size:16px;color:#fff;background:#007bff;border-radius:5px;text-decoration:none;transition:0.3s}a:hover{background:#0056b3}}</style></head><body><div class="container"><h2>你好，%s！</h2><p>你已成功完成授权登录流程。</p><a href="/">返回主页</a></div></body></html>\`, user.Data.Name)
    c.String(http.StatusOK, html)
}
```

本文内容
