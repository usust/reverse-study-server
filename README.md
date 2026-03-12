# reverse-study-server

## 启动

```bash
go run ./cmd/server
```

默认监听：`http://localhost:10000`

## 主要接口（v1）

### 模型 API 配置管理（可保存多条）

- `GET /v1/model-api/configs`
- `GET /v1/model-api/configs/:id`
- `POST /v1/model-api/configs`
- `PUT /v1/model-api/configs/:id`
- `DELETE /v1/model-api/configs/:id`

### 代码生成与编译下载

- `POST /v1/model-api/create-c-code`
  - 入参支持 `configId`（优先）或 `config`（直接传配置）
- `POST /v1/model-api/create-c-code/download`
  - 先生成 C 代码，再后端编译，直接返回下载流
