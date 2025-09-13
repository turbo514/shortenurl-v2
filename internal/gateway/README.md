# api接口
基本api
http://localhost:8080/v1/



# 租户管理

|路径|方法|描述|
|---|---|---|
|/tenants|POST|创建租户|
|/users|POST|创建用户|
|/sessions|POST|登录|

## 创建租户
- 请求消息体
```json
{
    "name": "公司a"
}
```

- 响应消息体
```json
{
    "tenant_id":"01859604-0c2d-7d84-a1a7-0242ac130002",
    "api_key":"123456"
}
```

## 创建用户
- 请求消息体
```json
{
    "tenant_id":"01859604-0c2d-7d84-a1a7-0242ac130002",
    "api_key":"123456",
    "name":"张三",
    "password":"1234"
}
```

- 响应消息体
```json
{
    
}
```

## 用户登录
- 请求消息体
```json
{
    "name":"张三",
    "password":"1234"
}

- 响应消息体
```json
{
    "token":"张三abc123"
}
```

## 错误示例:
```json
{
    "errmsg":"wrong password or name",
    "errcode":100001
}
```

# 链接处理

|路径|方法|描述|
|---|---|---|
|/shorten|POST|生成短链接|
|/resolve/{short_code}|GET|还原链接|

## 生成短链接
- 请求消息体
```json
{
    "user_id":"3333",
    "tenant_id":"01859604-0c2d-7d84-a1a7-0242ac130002",
    "original_url":"https://baidu.com",
}
```

- 响应消息体
```json
{
    "original_url":"https://baidu.com",
    "short_code":"abc123"
}
```

## 还原短链接
- 请求消息体
```json
param
```

- 响应消息体
```json
{
    "original_url":"https://baiduc.com"
}
```

# 数据查询

//...
